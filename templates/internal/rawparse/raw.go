// Package rawparse provides the lossless dummy format parser used to exercise
// the shared template parser without HTML or SQL knowledge.
package rawparse

import (
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/tinybind-go/templates/internal/syntax"
)

type Parser struct{}

type TextNode struct {
	Kind string          `json:"kind"`
	Pos  syntax.Position `json:"pos"`
	Text string          `json:"text"`
}

func (n *TextNode) NodeType() string { return n.Kind }

// Root creates a configurable root declaration backed by the dummy parser.
func Root(keyword, nodeType, outputPrefix string) syntax.RootDeclaration {
	return syntax.RootDeclaration{
		Keyword:      keyword,
		NodeType:     nodeType,
		OutputPrefix: outputPrefix,
		Context:      "raw:text",
		Parser:       Parser{},
	}
}

// Parse uses a generic template declaration so shared parser fixtures do not
// depend on component/HTML or statement/SQL behavior.
func Parse(filename string, source []byte) (*syntax.Module, error) {
	return syntax.ParseModule(filename, string(source), []syntax.RootDeclaration{
		Root("template", "raw:template", "raw"),
	})
}

func (Parser) ParseBody(context *syntax.BodyContext, insertionContext string) ([]syntax.Node, *syntax.Terminator, error) {
	source := context.Source()
	pos := context.Offset()
	textStart := pos
	var nodes []syntax.Node
	flush := func(end int) {
		if textStart == end {
			return
		}
		nodes = appendRawText(nodes, source[textStart:end], context.Position(textStart))
	}
	for pos < len(source) {
		switch {
		case strings.HasPrefix(source[pos:], "{{"):
			end := strings.Index(source[pos+2:], "}}")
			if end < 0 {
				return nil, nil, context.ErrorAt(pos, "unterminated escaped template text")
			}
			pos += end + 4
		case source[pos] == '{':
			flush(pos)
			fragment, end, err := readEmbedded(context, pos)
			if err != nil {
				return nil, nil, err
			}
			context.SetOffset(end)
			node, terminator, err := context.ParseEmbedded(fragment, insertionContext)
			if err != nil {
				return nil, nil, err
			}
			if terminator != nil {
				return nodes, terminator, nil
			}
			nodes = append(nodes, node)
			pos = context.Offset()
			textStart = pos
		case source[pos] == '}':
			flush(pos)
			terminator := &syntax.Terminator{Kind: syntax.TerminatorRoot, Pos: context.Position(pos)}
			pos++
			context.SetOffset(pos)
			return nodes, terminator, nil
		default:
			_, size := utf8.DecodeRuneInString(source[pos:])
			pos += size
		}
	}
	return nil, nil, context.ErrorAt(pos, "unterminated declaration body")
}

func readEmbedded(context *syntax.BodyContext, start int) (syntax.Embedded, int, error) {
	source := context.Source()
	pos := start + 1
	contentStart := pos
	depth := 0
	quote := byte(0)
	for pos < len(source) {
		c := source[pos]
		if quote != 0 {
			if c == '\\' {
				pos += 2
				continue
			}
			pos++
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			quote = c
			pos++
			continue
		}
		if c == '(' || c == '[' {
			depth++
		} else if c == ')' || c == ']' {
			depth--
		} else if c == '}' && depth == 0 {
			return syntax.Embedded{Text: source[contentStart:pos], StartOffset: start, ContentOffset: contentStart}, pos + 1, nil
		}
		pos++
	}
	return syntax.Embedded{}, 0, context.ErrorAt(start, "unterminated template expression")
}

func appendRawText(nodes []syntax.Node, text string, pos syntax.Position) []syntax.Node {
	if text == "" {
		return nodes
	}
	if len(nodes) > 0 {
		if previous, ok := nodes[len(nodes)-1].(*TextNode); ok {
			previous.Text += text
			return nodes
		}
	}
	return append(nodes, &TextNode{Kind: "raw:text", Pos: pos, Text: text})
}

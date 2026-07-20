package syntax

import "unicode/utf8"

// Position is a one-based source position. Offset is intentionally not part of
// the serialized AST; parsers use offsets internally while later compiler
// stages report the stable line and column pair.
type Position struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

func positionAt(source string, offset int) Position {
	if offset < 0 {
		offset = 0
	}
	if offset > len(source) {
		offset = len(source)
	}
	pos := Position{Line: 1, Col: 1}
	lineStart := 0
	for i := 0; i < offset; {
		r, size := utf8.DecodeRuneInString(source[i:])
		i += size
		if r == '\n' {
			pos.Line++
			lineStart = i
		}
	}
	pos.Col = utf8.RuneCountInString(source[lineStart:offset]) + 1
	return pos
}

func positionInFragment(source string, offset int, base Position) Position {
	local := positionAt(source, offset)
	if local.Line == 1 {
		return Position{Line: base.Line, Col: base.Col + local.Col - 1}
	}
	return Position{Line: base.Line + local.Line - 1, Col: local.Col}
}

// Package htmlbind parses typed HTML template sources into an AST.
package htmlbind

import "github.com/shibukawa/tinybind-go/templates/internal/syntax"

type Module = syntax.Module
type PackageDecl = syntax.PackageDecl
type ImportDecl = syntax.ImportDecl
type Declaration = syntax.Declaration
type TypeDecl = syntax.TypeDecl
type EnumDecl = syntax.EnumDecl
type EnumMember = syntax.EnumMember
type ExternalDecl = syntax.ExternalDecl
type TemplateDecl = syntax.TemplateDecl
type Field = syntax.Field
type Parameter = syntax.Parameter
type TypeRef = syntax.TypeRef
type Position = syntax.Position
type Expr = syntax.Expr
type IdentifierExpr = syntax.IdentifierExpr
type LiteralExpr = syntax.LiteralExpr
type MemberExpr = syntax.MemberExpr
type IndexExpr = syntax.IndexExpr
type CallExpr = syntax.CallExpr
type UnaryExpr = syntax.UnaryExpr
type BinaryExpr = syntax.BinaryExpr
type ConditionalExpr = syntax.ConditionalExpr
type ParseError = syntax.ParseError
type ExpressionNode = syntax.ExpressionNode
type IfNode = syntax.IfNode
type ForNode = syntax.ForNode

type Node = syntax.Node

// Body is the body stored in TemplateDecl.Body.
type Body = []Node

type TextNode struct {
	Kind string   `json:"kind"`
	Pos  Position `json:"pos"`
	Text string   `json:"text"`
}

func (n *TextNode) NodeType() string { return n.Kind }

type CommentNode struct {
	Kind string   `json:"kind"`
	Pos  Position `json:"pos"`
	Text string   `json:"text"`
}

func (n *CommentNode) NodeType() string { return n.Kind }

type DoctypeNode struct {
	Kind string   `json:"kind"`
	Pos  Position `json:"pos"`
	Text string   `json:"text"`
}

func (n *DoctypeNode) NodeType() string { return n.Kind }

type ElementNode struct {
	Kind        string      `json:"kind"`
	Pos         Position    `json:"pos"`
	Name        string      `json:"name"`
	Attributes  []Attribute `json:"attributes,omitempty"`
	Children    []Node      `json:"children,omitempty"`
	SelfClosing bool        `json:"selfClosing,omitempty"`
}

func (n *ElementNode) NodeType() string { return n.Kind }

type ComponentNode struct {
	Kind        string      `json:"kind"`
	Pos         Position    `json:"pos"`
	Name        string      `json:"name"`
	Arguments   []Attribute `json:"arguments,omitempty"`
	Children    []Node      `json:"children,omitempty"`
	SelfClosing bool        `json:"selfClosing,omitempty"`
}

func (n *ComponentNode) NodeType() string { return n.Kind }

type Attribute struct {
	Kind    string          `json:"kind"`
	Pos     Position        `json:"pos"`
	Name    string          `json:"name"`
	Boolean bool            `json:"boolean,omitempty"`
	Value   []AttributePart `json:"value,omitempty"`
}

type AttributePart struct {
	Kind       string   `json:"kind"`
	Pos        Position `json:"pos"`
	Context    string   `json:"context,omitempty"`
	Text       string   `json:"text,omitempty"`
	Expression Expr     `json:"expression,omitempty"`
}

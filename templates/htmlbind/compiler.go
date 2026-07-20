package htmlbind

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/tinybind-go/templates/internal/syntax"
)

type valueKind string

const (
	kindInvalid     valueKind = "invalid"
	kindString      valueKind = "string"
	kindBool        valueKind = "bool"
	kindInt         valueKind = "int"
	kindFloat       valueKind = "float"
	kindDecimal     valueKind = "decimal"
	kindDateTime    valueKind = "datetime"
	kindDate        valueKind = "date"
	kindTime        valueKind = "time"
	kindURL         valueKind = "url"
	kindBytes       valueKind = "bytes"
	kindRecord      valueKind = "record"
	kindEnum        valueKind = "enum"
	kindArray       valueKind = "array"
	kindHTML        valueKind = "html"
	kindTrustedHTML valueKind = "trusted_html"
	kindTrustedCSS  valueKind = "trusted_css"
	kindTrustedJS   valueKind = "trusted_javascript"
	kindScriptJSON  valueKind = "script_json"
)

type valueType struct {
	kind     valueKind
	name     string
	elem     *valueType
	optional bool
}

func (t valueType) String() string {
	base := t.name
	if base == "" {
		base = string(t.kind)
	}
	if t.kind == kindArray && t.elem != nil {
		base = "[]" + t.elem.String()
	}
	if t.optional {
		base += "?"
	}
	return base
}

func (t valueType) required() valueType { t.optional = false; return t }

type functionSig struct {
	params []valueType
	result valueType
}

type componentInfo struct {
	decl   *TemplateDecl
	params map[string]valueType
	order  []Parameter
}

type compiler struct {
	filename    string
	source      string
	module      *Module
	records     map[string]*TypeDecl
	enums       map[string]*EnumDecl
	enumMembers map[string]valueType
	externals   map[string]functionSig
	components  map[string]*componentInfo
	exprTypes   map[Expr]valueType
}

type CompileError struct {
	Filename string
	Pos      Position
	Message  string
}

func (e *CompileError) Error() string {
	name := e.Filename
	if name == "" {
		name = "<template>"
	}
	return fmt.Sprintf("%s:%d:%d: %s", name, e.Pos.Line, e.Pos.Col, e.Message)
}

func newCompiler(filename, source string, module *Module) *compiler {
	return &compiler{
		filename: filename, source: source, module: module,
		records: map[string]*TypeDecl{}, enums: map[string]*EnumDecl{},
		enumMembers: map[string]valueType{}, externals: map[string]functionSig{},
		components: map[string]*componentInfo{}, exprTypes: map[Expr]valueType{},
	}
}

func (c *compiler) analyze() error {
	for _, declaration := range c.module.Declarations {
		switch declaration := declaration.(type) {
		case *TypeDecl:
			if c.nameExists(declaration.Name) {
				return c.error(declaration.Pos, "duplicate declaration "+declaration.Name)
			}
			c.records[declaration.Name] = declaration
		case *EnumDecl:
			if c.nameExists(declaration.Name) {
				return c.error(declaration.Pos, "duplicate declaration "+declaration.Name)
			}
			c.enums[declaration.Name] = declaration
			for _, member := range declaration.Members {
				if _, exists := c.enumMembers[member.Name]; exists {
					return c.error(member.Pos, "duplicate enum member "+member.Name)
				}
				c.enumMembers[member.Name] = valueType{kind: kindEnum, name: declaration.Name}
			}
		}
	}
	for _, declaration := range c.module.Declarations {
		switch declaration := declaration.(type) {
		case *TypeDecl:
			seen := map[string]bool{}
			for _, field := range declaration.Fields {
				if seen[field.Name] {
					return c.error(field.Pos, "duplicate field "+field.Name)
				}
				seen[field.Name] = true
				if _, err := c.resolveType(field.Type); err != nil {
					return err
				}
			}
		case *ExternalDecl:
			if c.nameExists(declaration.Name) {
				return c.error(declaration.Pos, "duplicate declaration "+declaration.Name)
			}
			var sig functionSig
			for _, parameter := range declaration.Parameters {
				t, err := c.resolveType(parameter.Type)
				if err != nil {
					return err
				}
				sig.params = append(sig.params, t)
			}
			result, err := c.resolveType(declaration.Result)
			if err != nil {
				return err
			}
			sig.result = result
			c.externals[declaration.Name] = sig
		case *TemplateDecl:
			if declaration.Kind != "html:component" || declaration.Output.Name != "html" {
				return c.error(declaration.Pos, "HTML generator only accepts html:component declarations")
			}
			if c.nameExists(declaration.Name) {
				return c.error(declaration.Pos, "duplicate declaration "+declaration.Name)
			}
			info := &componentInfo{decl: declaration, params: map[string]valueType{}, order: declaration.Parameters}
			for _, parameter := range declaration.Parameters {
				if _, exists := info.params[parameter.Name]; exists {
					return c.error(parameter.Pos, "duplicate parameter "+parameter.Name)
				}
				t, err := c.resolveType(parameter.Type)
				if err != nil {
					return err
				}
				info.params[parameter.Name] = t
			}
			c.components[declaration.Name] = info
		}
	}
	for _, declaration := range c.module.Declarations {
		component, ok := declaration.(*TemplateDecl)
		if !ok {
			continue
		}
		info := c.components[component.Name]
		scope := copyScope(info.params)
		body, ok := component.Body.([]syntax.Node)
		if !ok {
			return c.error(component.Pos, "invalid HTML component body")
		}
		if err := c.analyzeNodes(body, scope); err != nil {
			return err
		}
	}
	return nil
}

func (c *compiler) nameExists(name string) bool {
	_, record := c.records[name]
	_, enum := c.enums[name]
	_, external := c.externals[name]
	_, component := c.components[name]
	return record || enum || external || component
}

func (c *compiler) resolveType(ref TypeRef) (valueType, error) {
	var result valueType
	switch ref.Name {
	case "string":
		result.kind = kindString
	case "bool":
		result.kind = kindBool
	case "int":
		result.kind = kindInt
	case "float":
		result.kind = kindFloat
	case "decimal":
		result.kind = kindDecimal
	case "datetime":
		result.kind = kindDateTime
	case "date":
		result.kind = kindDate
	case "time":
		result.kind = kindTime
	case "url":
		result.kind = kindURL
	case "bytes":
		result.kind = kindBytes
	case "html":
		result.kind = kindHTML
	case "trusted_html":
		result.kind = kindTrustedHTML
	case "trusted_css":
		result.kind = kindTrustedCSS
	case "trusted_javascript":
		result.kind = kindTrustedJS
	case "script_json":
		result.kind = kindScriptJSON
	default:
		if _, ok := c.records[ref.Name]; ok {
			result = valueType{kind: kindRecord, name: ref.Name}
		} else if _, ok := c.enums[ref.Name]; ok {
			result = valueType{kind: kindEnum, name: ref.Name}
		} else {
			return valueType{}, c.error(ref.Pos, "unknown type "+ref.Name)
		}
	}
	if ref.Array {
		elem := result
		result = valueType{kind: kindArray, elem: &elem}
	}
	result.optional = ref.Optional
	return result, nil
}

func (c *compiler) analyzeNodes(nodes []syntax.Node, scope map[string]valueType) error {
	for _, node := range nodes {
		switch node := node.(type) {
		case *TextNode, *CommentNode, *DoctypeNode:
		case *syntax.ExpressionNode:
			t, err := c.infer(node.Expression, scope)
			if err != nil {
				return err
			}
			if err := c.validateInsertion(node.Context, t, exprPos(node.Expression)); err != nil {
				return err
			}
		case *syntax.IfNode:
			t, err := c.infer(node.Condition, scope)
			if err != nil {
				return err
			}
			if t.kind != kindBool || t.optional {
				return c.error(exprPos(node.Condition), "if condition must be bool")
			}
			if err := c.analyzeNodes(node.Then, copyScope(scope)); err != nil {
				return err
			}
			if err := c.analyzeNodes(node.Else, copyScope(scope)); err != nil {
				return err
			}
		case *syntax.ForNode:
			t, err := c.infer(node.Iterable, scope)
			if err != nil {
				return err
			}
			if t.kind != kindArray || t.optional || t.elem == nil {
				return c.error(exprPos(node.Iterable), "for expression must be an array")
			}
			inner := copyScope(scope)
			inner[node.Variable] = *t.elem
			if node.Index != "" {
				inner[node.Index] = valueType{kind: kindInt}
			}
			if err := c.analyzeNodes(node.Body, inner); err != nil {
				return err
			}
		case *ElementNode:
			for _, attribute := range node.Attributes {
				if err := c.analyzeAttribute(node.Name, attribute, scope); err != nil {
					return err
				}
			}
			if err := c.analyzeNodes(node.Children, scope); err != nil {
				return err
			}
		case *ComponentNode:
			if err := c.analyzeComponentCall(node, scope); err != nil {
				return err
			}
		default:
			return c.error(Position{Line: 1, Col: 1}, fmt.Sprintf("unsupported HTML node %T", node))
		}
	}
	return nil
}

func (c *compiler) analyzeAttribute(element string, attribute Attribute, scope map[string]valueType) error {
	mixed := len(attribute.Value) > 1 || (len(attribute.Value) == 1 && attribute.Value[0].Expression == nil)
	for _, part := range attribute.Value {
		if part.Expression == nil {
			continue
		}
		t, err := c.infer(part.Expression, scope)
		if err != nil {
			return err
		}
		if isURLAttribute(attribute.Name) && t.required().kind != kindURL {
			return c.error(part.Pos, "attribute "+attribute.Name+" requires url, got "+t.String())
		}
		if mixed && t.optional {
			return c.error(part.Pos, "optional expression must be the entire attribute value")
		}
		if err := c.validateInsertion("html:attribute", t, part.Pos); err != nil {
			return err
		}
	}
	return nil
}

func (c *compiler) analyzeComponentCall(node *ComponentNode, scope map[string]valueType) error {
	component, ok := c.components[node.Name]
	if !ok {
		return c.error(node.Pos, "unknown component "+node.Name)
	}
	provided := map[string]bool{}
	for _, argument := range node.Arguments {
		want, ok := component.params[argument.Name]
		if !ok {
			return c.error(argument.Pos, "unknown argument "+argument.Name+" for "+node.Name)
		}
		if provided[argument.Name] {
			return c.error(argument.Pos, "duplicate argument "+argument.Name)
		}
		provided[argument.Name] = true
		got, err := c.attributeValueType(argument, scope)
		if err != nil {
			return err
		}
		if !assignable(want, got) {
			return c.error(argument.Pos, "argument "+argument.Name+" expects "+want.String()+", got "+got.String())
		}
	}
	for name := range component.params {
		if name == "children" {
			continue
		}
		if !provided[name] {
			return c.error(node.Pos, "missing argument "+name+" for "+node.Name)
		}
	}
	childrenType, acceptsChildren := component.params["children"]
	if len(node.Children) > 0 && (!acceptsChildren || childrenType.kind != kindHTML) {
		return c.error(node.Pos, "component "+node.Name+" does not accept children")
	}
	if len(node.Children) > 0 && provided["children"] {
		return c.error(node.Pos, "component "+node.Name+" received children twice")
	}
	if len(node.Children) > 0 && childrenType.optional {
		return c.error(node.Pos, "component "+node.Name+" children parameter must be html, not html?")
	}
	if acceptsChildren && !provided["children"] && len(node.Children) == 0 {
		return c.error(node.Pos, "component "+node.Name+" requires children")
	}
	return c.analyzeNodes(node.Children, scope)
}

func (c *compiler) attributeValueType(attribute Attribute, scope map[string]valueType) (valueType, error) {
	if attribute.Boolean {
		return valueType{kind: kindBool}, nil
	}
	if len(attribute.Value) == 1 && attribute.Value[0].Expression != nil {
		return c.infer(attribute.Value[0].Expression, scope)
	}
	for _, part := range attribute.Value {
		if part.Expression == nil {
			continue
		}
		t, err := c.infer(part.Expression, scope)
		if err != nil {
			return valueType{}, err
		}
		if t.required().kind != kindString {
			return valueType{}, c.error(part.Pos, "mixed attribute value expressions must be string")
		}
		if t.optional {
			return valueType{}, c.error(part.Pos, "optional expression must be the entire attribute value")
		}
	}
	return valueType{kind: kindString}, nil
}

func (c *compiler) validateInsertion(context string, t valueType, pos Position) error {
	base := t.required().kind
	switch context {
	case "html:child":
		if base == kindTrustedCSS || base == kindTrustedJS || base == kindScriptJSON || base == kindRecord || base == kindArray {
			return c.error(pos, "cannot insert "+t.String()+" into html:child")
		}
	case "html:attribute":
		if base == kindTrustedHTML || base == kindTrustedCSS || base == kindTrustedJS || base == kindScriptJSON || base == kindRecord || base == kindArray || base == kindHTML {
			return c.error(pos, "cannot insert "+t.String()+" into html:attribute")
		}
	case "html:script":
		if base != kindTrustedJS && base != kindScriptJSON {
			return c.error(pos, "html:script requires trusted_javascript or script_json")
		}
	case "html:style":
		if base != kindTrustedCSS {
			return c.error(pos, "html:style requires trusted_css")
		}
	default:
		return c.error(pos, "unknown HTML insertion context "+context)
	}
	return nil
}

func (c *compiler) infer(expr Expr, scope map[string]valueType) (valueType, error) {
	if known, ok := c.exprTypes[expr]; ok {
		return known, nil
	}
	var result valueType
	var err error
	switch expr := expr.(type) {
	case *IdentifierExpr:
		if t, ok := scope[expr.Name]; ok {
			result = t
		} else if t, ok := c.enumMembers[expr.Name]; ok {
			result = t
		} else {
			err = c.error(expr.Pos, "unknown identifier "+expr.Name)
		}
	case *LiteralExpr:
		switch expr.ValueKind {
		case "string":
			result.kind = kindString
		case "bool":
			result.kind = kindBool
		case "number":
			if strings.Contains(expr.Value.(string), ".") {
				result.kind = kindFloat
			} else {
				result.kind = kindInt
			}
		case "null":
			result = valueType{kind: kindInvalid, optional: true}
		default:
			err = c.error(expr.Pos, "unknown literal type")
		}
	case *MemberExpr:
		var object valueType
		object, err = c.infer(expr.Object, scope)
		if err == nil {
			if object.optional {
				err = c.error(expr.Pos, "member access on optional "+object.String())
			} else if object.kind != kindRecord {
				err = c.error(expr.Pos, "member access requires a record")
			} else {
				field, ok := findField(c.records[object.name], expr.Member)
				if !ok {
					err = c.error(expr.Pos, "unknown field "+expr.Member+" on "+object.name)
				} else {
					result, err = c.resolveType(field.Type)
				}
			}
		}
	case *IndexExpr:
		var object, index valueType
		object, err = c.infer(expr.Object, scope)
		if err == nil {
			index, err = c.infer(expr.Index, scope)
		}
		if err == nil && (object.kind != kindArray || object.optional) {
			err = c.error(expr.Pos, "indexing requires an array")
		}
		if err == nil && index.kind != kindInt {
			err = c.error(expr.Pos, "array index must be int")
		}
		if err == nil {
			result = *object.elem
		}
	case *CallExpr:
		result, err = c.inferCall(expr, scope)
	case *UnaryExpr:
		var operand valueType
		operand, err = c.infer(expr.Operand, scope)
		if err == nil {
			switch expr.Operator {
			case "!", "not":
				if operand.kind != kindBool || operand.optional {
					err = c.error(expr.Pos, "not requires bool")
				} else {
					result = operand
				}
			case "+", "-":
				if !numeric(operand) {
					err = c.error(expr.Pos, "numeric unary operator requires number")
				} else {
					result = operand
				}
			default:
				err = c.error(expr.Pos, "unsupported unary operator "+expr.Operator)
			}
		}
	case *BinaryExpr:
		var left, right valueType
		left, err = c.infer(expr.Left, scope)
		if err == nil {
			right, err = c.infer(expr.Right, scope)
		}
		if err == nil {
			result, err = c.binaryType(expr, left, right)
		}
	case *ConditionalExpr:
		var condition, thenType, elseType valueType
		condition, err = c.infer(expr.Condition, scope)
		if err == nil && (condition.kind != kindBool || condition.optional) {
			err = c.error(expr.Pos, "conditional condition must be bool")
		}
		if err == nil {
			thenType, err = c.infer(expr.Then, scope)
		}
		if err == nil {
			elseType, err = c.infer(expr.Else, scope)
		}
		if err == nil {
			if !assignable(thenType, elseType) || !assignable(elseType, thenType) {
				err = c.error(expr.Pos, "conditional branches must have the same type")
			} else {
				result = thenType
			}
		}
	default:
		err = c.error(Position{Line: 1, Col: 1}, fmt.Sprintf("unsupported expression %T", expr))
	}
	if err != nil {
		return valueType{}, err
	}
	c.exprTypes[expr] = result
	return result, nil
}

func (c *compiler) inferCall(call *CallExpr, scope map[string]valueType) (valueType, error) {
	identifier, ok := call.Callee.(*IdentifierExpr)
	if !ok {
		return valueType{}, c.error(call.Pos, "only named functions can be called")
	}
	if intrinsic, ok := intrinsicResult(identifier.Name); ok {
		if len(call.Arguments) != 1 {
			return valueType{}, c.error(call.Pos, identifier.Name+" expects one argument")
		}
		argument, err := c.infer(call.Arguments[0], scope)
		if err != nil {
			return valueType{}, err
		}
		if identifier.Name != "JsonForScript" && (argument.kind != kindString || argument.optional) {
			return valueType{}, c.error(call.Pos, identifier.Name+" expects string")
		}
		if identifier.Name == "JsonForScript" && !c.jsonSerializable(argument, map[string]bool{}) {
			return valueType{}, c.error(call.Pos, "JsonForScript argument is not statically serializable")
		}
		return intrinsic, nil
	}
	sig, ok := c.externals[identifier.Name]
	if !ok {
		return valueType{}, c.error(call.Pos, "unknown function "+identifier.Name)
	}
	if len(call.Arguments) != len(sig.params) {
		return valueType{}, c.error(call.Pos, fmt.Sprintf("%s expects %d arguments", identifier.Name, len(sig.params)))
	}
	for i, argument := range call.Arguments {
		got, err := c.infer(argument, scope)
		if err != nil {
			return valueType{}, err
		}
		if !assignable(sig.params[i], got) {
			return valueType{}, c.error(exprPos(argument), fmt.Sprintf("argument %d expects %s, got %s", i+1, sig.params[i], got))
		}
	}
	return sig.result, nil
}

func (c *compiler) binaryType(expr *BinaryExpr, left, right valueType) (valueType, error) {
	switch expr.Operator {
	case "and", "&&", "or", "||":
		if left.kind != kindBool || right.kind != kindBool || left.optional || right.optional {
			return valueType{}, c.error(expr.Pos, "boolean operator requires bool")
		}
		return valueType{kind: kindBool}, nil
	case "==", "!=":
		if left.kind == kindInvalid && left.optional {
			if !right.optional {
				return valueType{}, c.error(expr.Pos, "null can only compare with optional")
			}
			return valueType{kind: kindBool}, nil
		}
		if right.kind == kindInvalid && right.optional {
			if !left.optional {
				return valueType{}, c.error(expr.Pos, "null can only compare with optional")
			}
			return valueType{kind: kindBool}, nil
		}
		if !assignable(left, right) && !assignable(right, left) {
			return valueType{}, c.error(expr.Pos, "incompatible comparison")
		}
		if !c.comparable(left, map[string]bool{}) || !c.comparable(right, map[string]bool{}) {
			return valueType{}, c.error(expr.Pos, "values are not comparable")
		}
		return valueType{kind: kindBool}, nil
	case "<", "<=", ">", ">=":
		if !numeric(left) || !numeric(right) {
			return valueType{}, c.error(expr.Pos, "ordered comparison requires numbers")
		}
		return valueType{kind: kindBool}, nil
	case "+":
		if left.kind == kindString && right.kind == kindString && !left.optional && !right.optional {
			return valueType{kind: kindString}, nil
		}
		fallthrough
	case "-", "*", "/", "%":
		if !numeric(left) || !numeric(right) || left.kind != right.kind {
			return valueType{}, c.error(expr.Pos, "arithmetic operands must have the same numeric type")
		}
		return left, nil
	default:
		return valueType{}, c.error(expr.Pos, "unsupported binary operator "+expr.Operator)
	}
}

func intrinsicResult(name string) (valueType, bool) {
	switch name {
	case "RawHTML":
		return valueType{kind: kindTrustedHTML}, true
	case "RawCSS":
		return valueType{kind: kindTrustedCSS}, true
	case "RawJavaScript":
		return valueType{kind: kindTrustedJS}, true
	case "JsonForScript":
		return valueType{kind: kindScriptJSON}, true
	}
	return valueType{}, false
}

func assignable(want, got valueType) bool {
	if got.kind == kindInvalid && got.optional {
		return want.optional
	}
	if want.kind != got.kind || want.name != got.name || want.optional != got.optional {
		return false
	}
	if want.kind == kindArray {
		return want.elem != nil && got.elem != nil && assignable(*want.elem, *got.elem)
	}
	return true
}

func numeric(t valueType) bool {
	return !t.optional && (t.kind == kindInt || t.kind == kindFloat)
}
func (c *compiler) jsonSerializable(t valueType, visiting map[string]bool) bool {
	if t.optional {
		t.optional = false
	}
	switch t.kind {
	case kindString, kindBool, kindInt, kindFloat, kindDecimal, kindEnum:
		return true
	case kindArray:
		return t.elem != nil && c.jsonSerializable(*t.elem, visiting)
	case kindRecord:
		record, ok := c.records[t.name]
		if !ok {
			return false
		}
		if visiting[t.name] {
			return true
		}
		visiting[t.name] = true
		defer delete(visiting, t.name)
		for _, field := range record.Fields {
			fieldType, err := c.resolveType(field.Type)
			if err != nil || !c.jsonSerializable(fieldType, visiting) {
				return false
			}
		}
		return true
	}
	return false
}

func (c *compiler) comparable(t valueType, visiting map[string]bool) bool {
	if t.optional {
		return true
	}
	switch t.kind {
	case kindString, kindBool, kindInt, kindFloat, kindDecimal, kindDateTime, kindDate, kindTime, kindURL,
		kindEnum, kindTrustedHTML, kindTrustedCSS, kindTrustedJS, kindScriptJSON:
		return true
	case kindRecord:
		record := c.records[t.name]
		if record == nil {
			return false
		}
		if visiting[t.name] {
			return false
		}
		visiting[t.name] = true
		defer delete(visiting, t.name)
		for _, field := range record.Fields {
			fieldType, err := c.resolveType(field.Type)
			if err != nil || !c.comparable(fieldType, visiting) {
				return false
			}
		}
		return true
	}
	return false
}
func isURLAttribute(name string) bool {
	switch name {
	case "href", "src", "action", "formaction", "poster":
		return true
	}
	return false
}
func findField(record *TypeDecl, name string) (Field, bool) {
	if record == nil {
		return Field{}, false
	}
	for _, field := range record.Fields {
		if field.Name == name {
			return field, true
		}
	}
	return Field{}, false
}
func copyScope(in map[string]valueType) map[string]valueType {
	out := make(map[string]valueType, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
func (c *compiler) error(pos Position, message string) error {
	return &CompileError{Filename: c.filename, Pos: pos, Message: message}
}
func goPublicName(name string) string {
	if name == "" {
		return name
	}
	r, size := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[size:]
}
func exprPos(expr Expr) Position {
	switch expr := expr.(type) {
	case *IdentifierExpr:
		return expr.Pos
	case *LiteralExpr:
		return expr.Pos
	case *MemberExpr:
		return expr.Pos
	case *IndexExpr:
		return expr.Pos
	case *CallExpr:
		return expr.Pos
	case *UnaryExpr:
		return expr.Pos
	case *BinaryExpr:
		return expr.Pos
	case *ConditionalExpr:
		return expr.Pos
	}
	return Position{Line: 1, Col: 1}
}

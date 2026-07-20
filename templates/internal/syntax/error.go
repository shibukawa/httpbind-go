package syntax

import "fmt"

// ParseError reports a stable source location for syntax diagnostics.
type ParseError struct {
	Filename string
	Offset   int
	Line     int
	Column   int
	Message  string
}

func (e *ParseError) Error() string {
	name := e.Filename
	if name == "" {
		name = "<template>"
	}
	return fmt.Sprintf("%s:%d:%d: %s", name, e.Line, e.Column, e.Message)
}

func errorAt(filename, source string, offset int, message string) error {
	if offset < 0 {
		offset = 0
	}
	if offset > len(source) {
		offset = len(source)
	}
	line, column := 1, 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return &ParseError{Filename: filename, Offset: offset, Line: line, Column: column, Message: message}
}

// ErrorAt constructs a parser diagnostic for a format-specific body parser.
func ErrorAt(filename, source string, localOffset, baseOffset int, message string) error {
	return ErrorAtPosition(filename, source, localOffset, baseOffset, Position{Line: 1, Col: 1}, message)
}

// ErrorAtPosition constructs a diagnostic for a source fragment while keeping
// its file-global byte offset and line/column position.
func ErrorAtPosition(filename, source string, localOffset, baseOffset int, basePos Position, message string) error {
	err := errorAt(filename, source, localOffset, message).(*ParseError)
	err.Offset += baseOffset
	pos := positionInFragment(source, localOffset, basePos)
	err.Line = pos.Line
	err.Column = pos.Col
	return err
}

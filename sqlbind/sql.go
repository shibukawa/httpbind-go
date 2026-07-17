//go:build !tinygo

// Package sqlbind provides generated, reflection-free database/sql row mapping.
package sqlbind

import "database/sql"

// ScanRows maps joined SQL rows into a grouped object tree using generated code.
func ScanRows[T any](rows *sql.Rows) ([]T, error) {
	var zero []T
	fn, ok := lookupScanner(typeKey[T]())
	if !ok {
		return zero, missingScannerError(typeKey[T]())
	}
	out, err := fn(rows)
	if err != nil {
		return zero, err
	}
	return out.([]T), nil
}

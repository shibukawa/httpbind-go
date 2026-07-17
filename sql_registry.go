//go:build !tinygo

package httpbinder

import (
	"database/sql"
	"fmt"
	"sync"
)

type scanRowsFunc func(*sql.Rows) (any, error)

var scanners sync.Map

// RegisterScanRows registers a generated SQL tree scanner for T.
func RegisterScanRows[T any](fn func(*sql.Rows) ([]T, error)) {
	scanners.Store(typeKey[T](), scanRowsFunc(func(rows *sql.Rows) (any, error) { return fn(rows) }))
}

func lookupScanner(t any) (scanRowsFunc, bool) {
	v, ok := scanners.Load(t)
	if !ok {
		return nil, false
	}
	return v.(scanRowsFunc), true
}

func missingScannerError(t interface{ String() string }) error {
	return Internal(fmt.Errorf("httpbinder: no SQL scanner registered for %s", t.String()))
}

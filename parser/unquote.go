package parser

import "strconv"

func unquote(s string) (string, error) {
	return strconv.Unquote(s)
}

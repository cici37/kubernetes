//go:build NO_CEL_IDENT_ESCAPE

package model

func Escape(ident string) (string, bool) {
	return ident, true
}

func Unescape(escaped string) (string, bool) {
	return escaped, true
}

// Package sanitize provides whitelist-based sanitization for values written to logs.
// This addresses CWE-117 (Improper Output Neutralization for Logs) by ensuring
// only known-safe characters reach log output, preventing log forging via CR/LF injection.
package sanitize

import "strings"

// allowedLogRune defines the whitelist of characters permitted in log values.
// Any rune outside this set is replaced. The set covers all legitimate CloudEvent
// source and type values (URIs, dot-delimited identifiers) while excluding
// control characters, ANSI escapes, and unicode trickery by construction.
func allowedLogRune(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '.' || r == '_' || r == ':' ||
		r == '/' || r == '@' || r == '-' ||
		r == ' '
}

// LogValue sanitizes a string before it is written to a log entry.
// It applies a whitelist: only characters matching [a-zA-Z0-9._:/@\- ] are kept,
// all others are replaced with '_'.
func LogValue(s string) string {
	return strings.Map(func(r rune) rune {
		if allowedLogRune(r) {
			return r
		}
		return '_'
	}, s)
}

package xform

import "strings"

// FieldPath specifies the address of a field.
type FieldPath []string

// ParseFieldPath parses a string into a FieldPath.
// We may in future support jsonpath or similar, for now we simply support dot-separated field paths.
func ParseFieldPath(fieldPath string) (FieldPath, error) {
	return FieldPath(strings.Split(fieldPath, ".")), nil
}

// ParseFieldPaths parses string slice into a FieldPath slice.
// This simplifies error handling.
func ParseFieldPaths(fieldPaths []string) ([]FieldPath, error) {
	var out []FieldPath
	for _, fieldPath := range fieldPaths {
		p, err := ParseFieldPath(fieldPath)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

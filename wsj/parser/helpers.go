// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package parser

// bdup returns a copy of a slice
func bdup(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// toAnySlice helps us navigate Pigeon's nodes
func toAnySlice(v any) []any {
	if v == nil {
		return nil
	}
	return v.([]any)
}

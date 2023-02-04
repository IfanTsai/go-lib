//go:build go1.20

package byteutils

import (
	"unsafe"
)

// B2S converts byte slice to a string without memory allocation.
func B2S(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

// S2B converts string to a byte slice without memory allocation.
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

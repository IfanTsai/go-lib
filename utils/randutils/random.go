package randutils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt generates a random integer between min and max.
func RandomInt(min, max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))

	return min + n.Int64()
}

// RandomString generates a random string of length n.
func RandomString(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[RandomInt(0, int64(len(alphabet)-1))])
	}

	return sb.String()
}

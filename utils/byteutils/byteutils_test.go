package byteutils_test

import (
	"testing"

	"github.com/IfanTsai/go-lib/utils/byteutils"
	"github.com/stretchr/testify/require"
)

const s = "hello world"

func TestB2S(t *testing.T) {
	b := []byte(s)
	require.Equal(t, s, byteutils.B2S(b))
}

func TestS2B(t *testing.T) {
	b := []byte(s)
	require.Equal(t, b, byteutils.S2B(s))
}

func BenchmarkB2S(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		byteutils.B2S([]byte(s))
	}
}

func BenchmarkS2B(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		byteutils.S2B(s)
	}
}

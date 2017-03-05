package utils

import (
	"crypto/rand"
	"io"
	"testing"
)

const tokenLength = 32

// shortReader provides a broken implementation of io.Reader for testing.
type shortReader struct{}

func (sr shortReader) Read(p []byte) (int, error) {
	return len(p) % 2, io.ErrUnexpectedEOF
}

// TestRandomBytes tests the (extremely rare) case that crypto/rand does
// not return the expected number of bytes.
func TestRandomBytes(t *testing.T) {
	// Pioneered by https://github.com/justinas/nosurf
	original := rand.Reader
	rand.Reader = shortReader{}
	defer func() {
		rand.Reader = original
	}()

	var b = make([]byte, tokenLength)
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("RandomBytes did not report a short read: only read %d bytes", len(b))
		}
	}()

	b = RandomBytes(tokenLength)
}

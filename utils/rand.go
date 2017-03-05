package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomBytes returns securely generated random bytes. It will panic
// if the system's secure random number generator fails to function correctly.
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}

	return b
}

// RandomString returns a URL-safe, base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
func RandomString(n int) string {
	d := base64.URLEncoding.DecodedLen(n)
	buf := make([]byte, base64.URLEncoding.EncodedLen(d), n)
	base64.URLEncoding.Encode(buf, RandomBytes(d))
	for i := n - len(buf); i > 0; i-- {
		buf = append(buf, byte((base64.StdPadding)))
	}
	return string(buf)
}

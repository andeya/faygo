package utils

import (
	"crypto/rand"
	r "math/rand"
	"time"
)

var alphaNum = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`)

// RandomBytes generate random []byte by specify chars.
func RandomBytes(length int, alphabets ...byte) []byte {
	if len(alphabets) == 0 {
		alphabets = alphaNum
	}
	var bytes = make([]byte, length)
	var randBy bool
	if num, err := rand.Read(bytes); num != length || err != nil {
		r.Seed(time.Now().UnixNano())
		randBy = true
	}
	for i, b := range bytes {
		if randBy {
			bytes[i] = alphabets[r.Intn(len(alphabets))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return bytes
}

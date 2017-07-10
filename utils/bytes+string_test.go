package utils

import (
	"testing"
)

func TestBytesToString(t *testing.T) {
	bb := []byte("testing: BytesToString")
	ss := BytesToString(bb)
	t.Logf("type: %T, value: %v", ss, ss)
}

func TestStringToBytes(t *testing.T) {
	s := "testing: StringToBytes"
	b := StringToBytes(s)
	t.Logf("type: %T, value: %v, val-string: %s\n", b, b, b)
	b = append(b, '!')
	t.Logf("after append:\ntype: %T, value: %v, val-string: %s\n", b, b, b)
}

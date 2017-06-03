package utils

import (
	"testing"
)

func TestBytes2String(t *testing.T) {
	bb := []byte("testing: Bytes2String")
	ss := Bytes2String(bb)
	t.Logf("type: %T, value: %v", ss, ss)
}

func TestString2Bytes(t *testing.T) {
	s := "testing: String2Bytes"
	b := String2Bytes(&s)
	t.Logf("type: %T, value: %v, val-string: %s\n", b, b, b)
	b = append(b, '!')
	t.Logf("after append:\ntype: %T, value: %v, val-string: %s\n", b, b, b)
}

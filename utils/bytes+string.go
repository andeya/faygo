package utils

import (
	"unsafe"
)

// Byte2String convert []byte type to string type.
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes convert *string type to []byte type.
// NOTE: panic if modify the member value of the []byte.
func String2Bytes(s *string) []byte {
	sp := *(*[2]uintptr)(unsafe.Pointer(s))
	bp := [3]uintptr{sp[0], sp[1], sp[1]}
	return *(*[]byte)(unsafe.Pointer(&bp))
}

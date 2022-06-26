package ameda

import (
	"math"
	"strconv"
)

// Uint8ToInterface converts uint8 to interface.
func Uint8ToInterface(v uint8) interface{} {
	return v
}

// Uint8ToInterfacePtr converts uint8 to *interface.
func Uint8ToInterfacePtr(v uint8) *interface{} {
	r := Uint8ToInterface(v)
	return &r
}

// Uint8ToString converts uint8 to string.
func Uint8ToString(v uint8) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint8ToStringPtr converts uint8 to *string.
func Uint8ToStringPtr(v uint8) *string {
	r := Uint8ToString(v)
	return &r
}

// Uint8ToBool converts uint8 to bool.
func Uint8ToBool(v uint8) bool {
	return v != 0
}

// Uint8ToBoolPtr converts uint8 to *bool.
func Uint8ToBoolPtr(v uint8) *bool {
	r := Uint8ToBool(v)
	return &r
}

// Uint8ToFloat32 converts uint8 to float32.
func Uint8ToFloat32(v uint8) float32 {
	return float32(v)
}

// Uint8ToFloat32Ptr converts uint8 to *float32.
func Uint8ToFloat32Ptr(v uint8) *float32 {
	r := Uint8ToFloat32(v)
	return &r
}

// Uint8ToFloat64 converts uint8 to float64.
func Uint8ToFloat64(v uint8) float64 {
	return float64(v)
}

// Uint8ToFloat64Ptr converts uint8 to *float64.
func Uint8ToFloat64Ptr(v uint8) *float64 {
	r := Uint8ToFloat64(v)
	return &r
}

// Uint8ToInt converts uint8 to int.
func Uint8ToInt(v uint8) int {
	return int(v)
}

// Uint8ToIntPtr converts uint8 to *int.
func Uint8ToIntPtr(v uint8) *int {
	r := Uint8ToInt(v)
	return &r
}

// Uint8ToInt8 converts uint8 to int8.
func Uint8ToInt8(v uint8) (int8, error) {
	if v > math.MaxInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Uint8ToInt8Ptr converts uint8 to *int8.
func Uint8ToInt8Ptr(v uint8) (*int8, error) {
	r, err := Uint8ToInt8(v)
	return &r, err
}

// Uint8ToInt16 converts uint8 to int16.
func Uint8ToInt16(v uint8) int16 {
	return int16(v)
}

// Uint8ToInt16Ptr converts uint8 to *int16.
func Uint8ToInt16Ptr(v uint8) *int16 {
	r := Uint8ToInt16(v)
	return &r
}

// Uint8ToInt32 converts uint8 to int32.
func Uint8ToInt32(v uint8) int32 {
	return int32(v)
}

// Uint8ToInt32Ptr converts uint8 to *int32.
func Uint8ToInt32Ptr(v uint8) *int32 {
	r := Uint8ToInt32(v)
	return &r
}

// Uint8ToInt64 converts uint8 to int64.
func Uint8ToInt64(v uint8) int64 {
	return int64(v)
}

// Uint8ToInt64Ptr converts uint8 to *int64.
func Uint8ToInt64Ptr(v uint8) *int64 {
	r := Uint8ToInt64(v)
	return &r
}

// Uint8ToUint converts uint8 to uint.
func Uint8ToUint(v uint8) uint {
	return uint(v)
}

// Uint8ToUintPtr converts uint8 to *uint.
func Uint8ToUintPtr(v uint8) *uint {
	r := Uint8ToUint(v)
	return &r
}

// Uint8ToUint8Ptr converts uint8 to *uint8.
func Uint8ToUint8Ptr(v uint8) *uint8 {
	return &v
}

// Uint8ToUint16 converts uint8 to uint16.
func Uint8ToUint16(v uint8) uint16 {
	return uint16(v)
}

// Uint8ToUint16Ptr converts uint8 to *uint16.
func Uint8ToUint16Ptr(v uint8) *uint16 {
	r := Uint8ToUint16(v)
	return &r
}

// Uint8ToUint32 converts uint8 to uint32.
func Uint8ToUint32(v uint8) uint32 {
	return uint32(v)
}

// Uint8ToUint32Ptr converts uint8 to *uint32.
func Uint8ToUint32Ptr(v uint8) *uint32 {
	r := Uint8ToUint32(v)
	return &r
}

// Uint8ToUint64 converts uint8 to uint64.
func Uint8ToUint64(v uint8) uint64 {
	return uint64(v)
}

// Uint8ToUint64Ptr converts uint8 to *uint64.
func Uint8ToUint64Ptr(v uint8) *uint64 {
	r := Uint8ToUint64(v)
	return &r
}

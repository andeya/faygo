package ameda

import (
	"math"
	"strconv"
)

// Uint32ToInterface converts uint32 to interface.
func Uint32ToInterface(v uint32) interface{} {
	return v
}

// Uint32ToInterfacePtr converts uint32 to *interface.
func Uint32ToInterfacePtr(v uint32) *interface{} {
	r := Uint32ToInterface(v)
	return &r
}

// Uint32ToString converts uint32 to string.
func Uint32ToString(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint32ToStringPtr converts uint32 to *string.
func Uint32ToStringPtr(v uint32) *string {
	r := Uint32ToString(v)
	return &r
}

// Uint32ToBool converts uint32 to bool.
func Uint32ToBool(v uint32) bool {
	return v != 0
}

// Uint32ToBoolPtr converts uint32 to *bool.
func Uint32ToBoolPtr(v uint32) *bool {
	r := Uint32ToBool(v)
	return &r
}

// Uint32ToFloat32 converts uint32 to float32.
func Uint32ToFloat32(v uint32) float32 {
	return float32(v)
}

// Uint32ToFloat32Ptr converts uint32 to *float32.
func Uint32ToFloat32Ptr(v uint32) *float32 {
	r := Uint32ToFloat32(v)
	return &r
}

// Uint32ToFloat64 converts uint32 to float64.
func Uint32ToFloat64(v uint32) float64 {
	return float64(v)
}

// Uint32ToFloat64Ptr converts uint32 to *float64.
func Uint32ToFloat64Ptr(v uint32) *float64 {
	r := Uint32ToFloat64(v)
	return &r
}

// Uint32ToInt converts uint32 to int.
func Uint32ToInt(v uint32) int {
	return int(v)
}

// Uint32ToIntPtr converts uint32 to *int.
func Uint32ToIntPtr(v uint32) *int {
	r := Uint32ToInt(v)
	return &r
}

// Uint32ToInt8 converts uint32 to int8.
func Uint32ToInt8(v uint32) (int8, error) {
	if v > math.MaxInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Uint32ToInt8Ptr converts uint32 to *int8.
func Uint32ToInt8Ptr(v uint32) (*int8, error) {
	r, err := Uint32ToInt8(v)
	return &r, err
}

// Uint32ToInt16 converts uint32 to int16.
func Uint32ToInt16(v uint32) (int16, error) {
	if v > math.MaxInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// Uint32ToInt16Ptr converts uint32 to *int16.
func Uint32ToInt16Ptr(v uint32) (*int16, error) {
	r, err := Uint32ToInt16(v)
	return &r, err
}

// Uint32ToInt32 converts uint32 to int32.
func Uint32ToInt32(v uint32) (int32, error) {
	if v > math.MaxInt32 {
		return 0, errOverflowValue
	}
	return int32(v), nil
}

// Uint32ToInt32Ptr converts uint32 to *int32.
func Uint32ToInt32Ptr(v uint32) (*int32, error) {
	r, err := Uint32ToInt32(v)
	return &r, err
}

// Uint32ToInt64 converts uint32 to int64.
func Uint32ToInt64(v uint32) int64 {
	return int64(v)
}

// Uint32ToInt64Ptr converts uint32 to *int64.
func Uint32ToInt64Ptr(v uint32) *int64 {
	r := Uint32ToInt64(v)
	return &r
}

// Uint32ToUint converts uint32 to uint.
func Uint32ToUint(v uint32) uint {
	return uint(v)
}

// Uint32ToUintPtr converts uint32 to *uint.
func Uint32ToUintPtr(v uint32) *uint {
	r := Uint32ToUint(v)
	return &r
}

// Uint32ToUint8 converts uint32 to uint8.
func Uint32ToUint8(v uint32) (uint8, error) {
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Uint32ToUint8Ptr converts uint32 to *uint8.
func Uint32ToUint8Ptr(v uint32) (*uint8, error) {
	r, err := Uint32ToUint8(v)
	return &r, err
}

// Uint32ToUint16 converts uint32 to uint16.
func Uint32ToUint16(v uint32) (uint16, error) {
	if v > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// Uint32ToUint16Ptr converts uint32 to *uint16.
func Uint32ToUint16Ptr(v uint32) (*uint16, error) {
	r, err := Uint32ToUint16(v)
	return &r, err
}

// Uint32ToUint32Ptr converts uint32 to *uint32.
func Uint32ToUint32Ptr(v uint32) *uint32 {
	return &v
}

// Uint32ToUint64 converts uint32 to uint64.
func Uint32ToUint64(v uint32) uint64 {
	return uint64(v)
}

// Uint32ToUint64Ptr converts uint32 to *uint64.
func Uint32ToUint64Ptr(v uint32) *uint64 {
	r := Uint32ToUint64(v)
	return &r
}

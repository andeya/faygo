package ameda

import (
	"math"
	"strconv"
)

// Uint16ToInterface converts uint16 to interface.
func Uint16ToInterface(v uint16) interface{} {
	return v
}

// Uint16ToInterfacePtr converts uint16 to *interface.
func Uint16ToInterfacePtr(v uint16) *interface{} {
	r := Uint16ToInterface(v)
	return &r
}

// Uint16ToString converts uint16 to string.
func Uint16ToString(v uint16) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint16ToStringPtr converts uint16 to *string.
func Uint16ToStringPtr(v uint16) *string {
	r := Uint16ToString(v)
	return &r
}

// Uint16ToBool converts uint16 to bool.
func Uint16ToBool(v uint16) bool {
	return v != 0
}

// Uint16ToBoolPtr converts uint16 to *bool.
func Uint16ToBoolPtr(v uint16) *bool {
	r := Uint16ToBool(v)
	return &r
}

// Uint16ToFloat32 converts uint16 to float32.
func Uint16ToFloat32(v uint16) float32 {
	return float32(v)
}

// Uint16ToFloat32Ptr converts uint16 to *float32.
func Uint16ToFloat32Ptr(v uint16) *float32 {
	r := Uint16ToFloat32(v)
	return &r
}

// Uint16ToFloat64 converts uint16 to float64.
func Uint16ToFloat64(v uint16) float64 {
	return float64(v)
}

// Uint16ToFloat64Ptr converts uint16 to *float64.
func Uint16ToFloat64Ptr(v uint16) *float64 {
	r := Uint16ToFloat64(v)
	return &r
}

// Uint16ToInt converts uint16 to int.
func Uint16ToInt(v uint16) int {
	return int(v)
}

// Uint16ToIntPtr converts uint16 to *int.
func Uint16ToIntPtr(v uint16) *int {
	r := Uint16ToInt(v)
	return &r
}

// Uint16ToInt8 converts uint16 to int8.
func Uint16ToInt8(v uint16) (int8, error) {
	if v > math.MaxInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Uint16ToInt8Ptr converts uint16 to *int8.
func Uint16ToInt8Ptr(v uint16) (*int8, error) {
	r, err := Uint16ToInt8(v)
	return &r, err
}

// Uint16ToInt16 converts uint16 to int16.
func Uint16ToInt16(v uint16) (int16, error) {
	if v > math.MaxInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// Uint16ToInt16Ptr converts uint16 to *int16.
func Uint16ToInt16Ptr(v uint16) (*int16, error) {
	r, err := Uint16ToInt16(v)
	return &r, err
}

// Uint16ToInt32 converts uint16 to int32.
func Uint16ToInt32(v uint16) int32 {
	return int32(v)
}

// Uint16ToInt32Ptr converts uint16 to *int32.
func Uint16ToInt32Ptr(v uint16) *int32 {
	r := Uint16ToInt32(v)
	return &r
}

// Uint16ToInt64 converts uint16 to int64.
func Uint16ToInt64(v uint16) int64 {
	return int64(v)
}

// Uint16ToInt64Ptr converts uint16 to *int64.
func Uint16ToInt64Ptr(v uint16) *int64 {
	r := Uint16ToInt64(v)
	return &r
}

// Uint16ToUint converts uint16 to uint.
func Uint16ToUint(v uint16) uint {
	return uint(v)
}

// Uint16ToUintPtr converts uint16 to *uint.
func Uint16ToUintPtr(v uint16) *uint {
	r := Uint16ToUint(v)
	return &r
}

// Uint16ToUint8 converts uint16 to uint8.
func Uint16ToUint8(v uint16) (uint8, error) {
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Uint16ToUint8Ptr converts uint16 to *uint8.
func Uint16ToUint8Ptr(v uint16) (*uint8, error) {
	r, err := Uint16ToUint8(v)
	return &r, err
}

// Uint16ToUint16Ptr converts uint16 to *uint16.
func Uint16ToUint16Ptr(v uint16) *uint16 {
	return &v
}

// Uint16ToUint32 converts uint16 to uint32.
func Uint16ToUint32(v uint16) uint32 {
	return uint32(v)
}

// Uint16ToUint32Ptr converts uint16 to *uint32.
func Uint16ToUint32Ptr(v uint16) *uint32 {
	r := Uint16ToUint32(v)
	return &r
}

// Uint16ToUint64 converts uint16 to uint64.
func Uint16ToUint64(v uint16) uint64 {
	return uint64(v)
}

// Uint16ToUint64Ptr converts uint16 to *uint64.
func Uint16ToUint64Ptr(v uint16) *uint64 {
	r := Uint16ToUint64(v)
	return &r
}

package ameda

import (
	"math"
	"strconv"
)

// Uint64ToInterface converts uint64 to interface.
func Uint64ToInterface(v uint64) interface{} {
	return v
}

// Uint64ToInterfacePtr converts uint64 to *interface.
func Uint64ToInterfacePtr(v uint64) *interface{} {
	r := Uint64ToInterface(v)
	return &r
}

// Uint64ToString converts uint64 to string.
func Uint64ToString(v uint64) string {
	return strconv.FormatUint(uint64(v), 10)
}

// Uint64ToStringPtr converts uint64 to *string.
func Uint64ToStringPtr(v uint64) *string {
	r := Uint64ToString(v)
	return &r
}

// Uint64ToBool converts uint64 to bool.
func Uint64ToBool(v uint64) bool {
	return v != 0
}

// Uint64ToBoolPtr converts uint64 to *bool.
func Uint64ToBoolPtr(v uint64) *bool {
	r := Uint64ToBool(v)
	return &r
}

// Uint64ToFloat32 converts uint64 to float32.
func Uint64ToFloat32(v uint64) float32 {
	return float32(v)
}

// Uint64ToFloat32Ptr converts uint64 to *float32.
func Uint64ToFloat32Ptr(v uint64) *float32 {
	r := Uint64ToFloat32(v)
	return &r
}

// Uint64ToFloat64 converts uint64 to float64.
func Uint64ToFloat64(v uint64) float64 {
	return float64(v)
}

// Uint64ToFloat64Ptr converts uint64 to *float64.
func Uint64ToFloat64Ptr(v uint64) *float64 {
	r := Uint64ToFloat64(v)
	return &r
}

// Uint64ToInt converts uint64 to int.
func Uint64ToInt(v uint64) int {
	return int(v)
}

// Uint64ToIntPtr converts uint64 to *int.
func Uint64ToIntPtr(v uint64) *int {
	r := Uint64ToInt(v)
	return &r
}

// Uint64ToInt8 converts uint64 to int8.
func Uint64ToInt8(v uint64) (int8, error) {
	if v > math.MaxInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Uint64ToInt8Ptr converts uint64 to *int8.
func Uint64ToInt8Ptr(v uint64) (*int8, error) {
	r, err := Uint64ToInt8(v)
	return &r, err
}

// Uint64ToInt16 converts uint64 to int16.
func Uint64ToInt16(v uint64) (int16, error) {
	if v > math.MaxInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// Uint64ToInt16Ptr converts uint64 to *int16.
func Uint64ToInt16Ptr(v uint64) (*int16, error) {
	r, err := Uint64ToInt16(v)
	return &r, err
}

// Uint64ToInt32 converts uint64 to int32.
func Uint64ToInt32(v uint64) (int32, error) {
	if v > math.MaxInt32 {
		return 0, errOverflowValue
	}
	return int32(v), nil
}

// Uint64ToInt32Ptr converts uint64 to *int32.
func Uint64ToInt32Ptr(v uint64) (*int32, error) {
	r, err := Uint64ToInt32(v)
	return &r, err
}

// Uint64ToInt64 converts uint64 to int64.
func Uint64ToInt64(v uint64) (int64, error) {
	if v > math.MaxInt64 {
		return 0, errOverflowValue
	}
	return int64(v), nil
}

// Uint64ToInt64Ptr converts uint64 to *int64.
func Uint64ToInt64Ptr(v uint64) (*int64, error) {
	r, err := Uint64ToInt64(v)
	return &r, err
}

// Uint64ToUint converts uint64 to uint.
func Uint64ToUint(v uint64) (uint, error) {
	if !Host64bit && v > math.MaxUint32 {
		return 0, errOverflowValue
	}
	return uint(v), nil
}

// Uint64ToUintPtr converts uint64 to *uint.
func Uint64ToUintPtr(v uint64) (*uint, error) {
	r, err := Uint64ToUint(v)
	return &r, err
}

// Uint64ToUint8 converts uint64 to uint8.
func Uint64ToUint8(v uint64) (uint8, error) {
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Uint64ToUint8Ptr converts uint64 to *uint8.
func Uint64ToUint8Ptr(v uint64) (*uint8, error) {
	r, err := Uint64ToUint8(v)
	return &r, err
}

// Uint64ToUint16 converts uint64 to uint16.
func Uint64ToUint16(v uint64) (uint16, error) {
	if v > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// Uint64ToUint16Ptr converts uint64 to *uint16.
func Uint64ToUint16Ptr(v uint64) (*uint16, error) {
	r, err := Uint64ToUint16(v)
	return &r, err
}

// Uint64ToUint32 converts uint64 to uint32.
func Uint64ToUint32(v uint64) (uint32, error) {
	if v > math.MaxUint32 {
		return 0, errOverflowValue
	}
	return uint32(v), nil
}

// Uint64ToUint32Ptr converts uint64 to *uint32.
func Uint64ToUint32Ptr(v uint64) (*uint32, error) {
	r, err := Uint64ToUint32(v)
	return &r, err
}

// Uint64ToUint64Ptr converts uint64 to *uint64.
func Uint64ToUint64Ptr(v uint64) *uint64 {
	return &v
}

package ameda

import (
	"math"
	"strconv"
)

// Int32ToInterface converts int32 to interface.
func Int32ToInterface(v int32) interface{} {
	return v
}

// Int32ToInterfacePtr converts int32 to *interface.
func Int32ToInterfacePtr(v int32) *interface{} {
	r := Int32ToInterface(v)
	return &r
}

// Int32ToString converts int32 to string.
func Int32ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int32ToStringPtr converts int32 to *string.
func Int32ToStringPtr(v int32) *string {
	r := Int32ToString(v)
	return &r
}

// Int32ToBool converts int32 to bool.
func Int32ToBool(v int32) bool {
	return v != 0
}

// Int32ToBoolPtr converts int32 to *bool.
func Int32ToBoolPtr(v int32) *bool {
	r := Int32ToBool(v)
	return &r
}

// Int32ToFloat32 converts int32 to float32.
func Int32ToFloat32(v int32) float32 {
	return float32(v)
}

// Int32ToFloat32Ptr converts int32 to *float32.
func Int32ToFloat32Ptr(v int32) *float32 {
	r := Int32ToFloat32(v)
	return &r
}

// Int32ToFloat64 converts int32 to float64.
func Int32ToFloat64(v int32) float64 {
	return float64(v)
}

// Int32ToFloat64Ptr converts int32 to *float64.
func Int32ToFloat64Ptr(v int32) *float64 {
	r := Int32ToFloat64(v)
	return &r
}

// Int32ToInt converts int32 to int.
func Int32ToInt(v int32) int {
	return int(v)
}

// Int32ToIntPtr converts int32 to *int.
func Int32ToIntPtr(v int32) *int {
	r := Int32ToInt(v)
	return &r
}

// Int32ToInt8 converts int32 to int8.
func Int32ToInt8(v int32) (int8, error) {
	if v > math.MaxInt8 || v < math.MinInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Int32ToInt8Ptr converts int32 to *int8.
func Int32ToInt8Ptr(v int32) (*int8, error) {
	r, err := Int32ToInt8(v)
	return &r, err
}

// Int32ToInt16 converts int32 to int16.
func Int32ToInt16(v int32) (int16, error) {
	if v > math.MaxInt16 || v < math.MinInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// Int32ToInt16Ptr converts int32 to *int16.
func Int32ToInt16Ptr(v int32) (*int16, error) {
	r, err := Int32ToInt16(v)
	return &r, err
}

// Int32ToInt32Ptr converts int32 to *int32.
func Int32ToInt32Ptr(v int32) *int32 {
	return &v
}

// Int32ToInt64 converts int32 to int64.
func Int32ToInt64(v int32) int64 {
	return int64(v)
}

// Int32ToInt64Ptr converts int32 to *int64.
func Int32ToInt64Ptr(v int32) *int64 {
	r := Int32ToInt64(v)
	return &r
}

// Int32ToUint converts int32 to uint.
func Int32ToUint(v int32) (uint, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint(v), nil
}

// Int32ToUintPtr converts int32 to *uint.
func Int32ToUintPtr(v int32) (*uint, error) {
	r, err := Int32ToUint(v)
	return &r, err
}

// Int32ToUint8 converts int32 to uint8.
func Int32ToUint8(v int32) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Int32ToUint8Ptr converts int32 to *uint8.
func Int32ToUint8Ptr(v int32) (*uint8, error) {
	r, err := Int32ToUint8(v)
	return &r, err
}

// Int32ToUint16 converts int32 to uint16.
func Int32ToUint16(v int32) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// Int32ToUint16Ptr converts int32 to *uint16.
func Int32ToUint16Ptr(v int32) (*uint16, error) {
	r, err := Int32ToUint16(v)
	return &r, err
}

// Int32ToUint32 converts int32 to uint32.
func Int32ToUint32(v int32) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint32(v), nil
}

// Int32ToUint32Ptr converts int32 to *uint32.
func Int32ToUint32Ptr(v int32) (*uint32, error) {
	r, err := Int32ToUint32(v)
	return &r, err
}

// Int32ToUint64 converts int32 to uint64.
func Int32ToUint64(v int32) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint64(v), nil
}

// Int32ToUint64Ptr converts int32 to *uint64.
func Int32ToUint64Ptr(v int32) (*uint64, error) {
	r, err := Int32ToUint64(v)
	return &r, err
}

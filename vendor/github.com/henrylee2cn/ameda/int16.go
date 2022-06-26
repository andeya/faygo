package ameda

import (
	"math"
	"strconv"
)

// Int16ToInterface converts int16 to interface.
func Int16ToInterface(v int16) interface{} {
	return v
}

// Int16ToInterfacePtr converts int16 to *interface.
func Int16ToInterfacePtr(v int16) *interface{} {
	r := Int16ToInterface(v)
	return &r
}

// Int16ToString converts int16 to string.
func Int16ToString(v int16) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int16ToStringPtr converts int16 to *string.
func Int16ToStringPtr(v int16) *string {
	r := Int16ToString(v)
	return &r
}

// Int16ToBool converts int16 to bool.
func Int16ToBool(v int16) bool {
	return v != 0
}

// Int16ToBoolPtr converts int16 to *bool.
func Int16ToBoolPtr(v int16) *bool {
	r := Int16ToBool(v)
	return &r
}

// Int16ToFloat32 converts int16 to float32.
func Int16ToFloat32(v int16) float32 {
	return float32(v)
}

// Int16ToFloat32Ptr converts int16 to *float32.
func Int16ToFloat32Ptr(v int16) *float32 {
	r := Int16ToFloat32(v)
	return &r
}

// Int16ToFloat64 converts int16 to float64.
func Int16ToFloat64(v int16) float64 {
	return float64(v)
}

// Int16ToFloat64Ptr converts int16 to *float64.
func Int16ToFloat64Ptr(v int16) *float64 {
	r := Int16ToFloat64(v)
	return &r
}

// Int16ToInt converts int16 to int.
func Int16ToInt(v int16) int {
	return int(v)
}

// Int16ToIntPtr converts int16 to *int.
func Int16ToIntPtr(v int16) *int {
	r := Int16ToInt(v)
	return &r
}

// Int16ToInt8 converts int16 to int8.
func Int16ToInt8(v int16) (int8, error) {
	if v > math.MaxInt8 || v < math.MinInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Int16ToInt8Ptr converts int16 to *int8.
func Int16ToInt8Ptr(v int16) (*int8, error) {
	r, err := Int16ToInt8(v)
	return &r, err
}

// Int16ToInt16Ptr converts int16 to *int16.
func Int16ToInt16Ptr(v int16) *int16 {
	return &v
}

// Int16ToInt32 converts int16 to int32.
func Int16ToInt32(v int16) int32 {
	return int32(v)
}

// Int16ToInt32Ptr converts int16 to *int32.
func Int16ToInt32Ptr(v int16) *int32 {
	r := Int16ToInt32(v)
	return &r
}

// Int16ToInt64 converts int16 to int64.
func Int16ToInt64(v int16) int64 {
	return int64(v)
}

// Int16ToInt64Ptr converts int16 to *int64.
func Int16ToInt64Ptr(v int16) *int64 {
	r := Int16ToInt64(v)
	return &r
}

// Int16ToUint converts int16 to uint.
func Int16ToUint(v int16) (uint, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint(v), nil
}

// Int16ToUintPtr converts int16 to *uint.
func Int16ToUintPtr(v int16) (*uint, error) {
	r, err := Int16ToUint(v)
	return &r, err
}

// Int16ToUint8 converts int16 to uint8.
func Int16ToUint8(v int16) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Int16ToUint8Ptr converts int16 to *uint8.
func Int16ToUint8Ptr(v int16) (*uint8, error) {
	r, err := Int16ToUint8(v)
	return &r, err
}

// Int16ToUint16 converts int16 to uint16.
func Int16ToUint16(v int16) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint16(v), nil
}

// Int16ToUint16Ptr converts int16 to *uint16.
func Int16ToUint16Ptr(v int16) (*uint16, error) {
	r, err := Int16ToUint16(v)
	return &r, err
}

// Int16ToUint32 converts int16 to uint32.
func Int16ToUint32(v int16) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint32(v), nil
}

// Int16ToUint32Ptr converts int16 to *uint32.
func Int16ToUint32Ptr(v int16) (*uint32, error) {
	r, err := Int16ToUint32(v)
	return &r, err
}

// Int16ToUint64 converts int16 to uint64.
func Int16ToUint64(v int16) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint64(v), nil
}

// Int16ToUint64Ptr converts int16 to *uint64.
func Int16ToUint64Ptr(v int16) (*uint64, error) {
	r, err := Int16ToUint64(v)
	return &r, err
}

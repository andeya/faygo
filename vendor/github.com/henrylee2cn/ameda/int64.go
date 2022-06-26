package ameda

import (
	"math"
	"strconv"
)

// Int64ToInterface converts int64 to interface.
func Int64ToInterface(v int64) interface{} {
	return v
}

// Int64ToInterfacePtr converts int64 to *interface.
func Int64ToInterfacePtr(v int64) *interface{} {
	r := Int64ToInterface(v)
	return &r
}

// Int64ToString converts int64 to string.
func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

// Int64ToStringPtr converts int64 to *string.
func Int64ToStringPtr(v int64) *string {
	r := Int64ToString(v)
	return &r
}

// Int64ToBool converts int64 to bool.
func Int64ToBool(v int64) bool {
	return v != 0
}

// Int64ToBoolPtr converts int64 to *bool.
func Int64ToBoolPtr(v int64) *bool {
	r := Int64ToBool(v)
	return &r
}

// Int64ToFloat32 converts int64 to float32.
func Int64ToFloat32(v int64) float32 {
	return float32(v)
}

// Int64ToFloat32Ptr converts int64 to *float32.
func Int64ToFloat32Ptr(v int64) *float32 {
	r := Int64ToFloat32(v)
	return &r
}

// Int64ToFloat64 converts int64 to float64.
func Int64ToFloat64(v int64) float64 {
	return float64(v)
}

// Int64ToFloat64Ptr converts int64 to *float64.
func Int64ToFloat64Ptr(v int64) *float64 {
	r := Int64ToFloat64(v)
	return &r
}

// Int64ToInt converts int64 to int.
func Int64ToInt(v int64) (int, error) {
	if !Host64bit && v > math.MaxInt32 {
		return 0, errOverflowValue
	}
	return int(v), nil
}

// Int64ToIntPtr converts int64 to *int.
func Int64ToIntPtr(v int64) (*int, error) {
	r, err := Int64ToInt(v)
	return &r, err
}

// Int64ToInt8 converts int64 to int8.
func Int64ToInt8(v int64) (int8, error) {
	if v > math.MaxInt8 || v < math.MinInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Int64ToInt8Ptr converts int64 to *int8.
func Int64ToInt8Ptr(v int64) (*int8, error) {
	r, err := Int64ToInt8(v)
	return &r, err
}

// Int64ToInt16 converts int64 to int16.
func Int64ToInt16(v int64) (int16, error) {
	if v > math.MaxInt16 || v < math.MinInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// Int64ToInt16Ptr converts int64 to *int16.
func Int64ToInt16Ptr(v int64) (*int16, error) {
	r, err := Int64ToInt16(v)
	return &r, err
}

// Int64ToInt32 converts int64 to int32.
func Int64ToInt32(v int64) (int32, error) {
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, errOverflowValue
	}
	return int32(v), nil
}

// Int64ToInt32Ptr converts int64 to *int32.
func Int64ToInt32Ptr(v int64) (*int32, error) {
	r, err := Int64ToInt32(v)
	return &r, err
}

// Int64ToInt64Ptr converts int64 to *int64.
func Int64ToInt64Ptr(v int64) *int64 {
	return &v
}

// Int64ToUint converts int64 to uint.
func Int64ToUint(v int64) (uint, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if !Host64bit && v > math.MaxUint32 {
		return 0, errOverflowValue
	}
	return uint(v), nil
}

// Int64ToUintPtr converts int64 to *uint.
func Int64ToUintPtr(v int64) (*uint, error) {
	r, err := Int64ToUint(v)
	return &r, err
}

// Int64ToUint8 converts int64 to uint8.
func Int64ToUint8(v int64) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Int64ToUint8Ptr converts int64 to *uint8.
func Int64ToUint8Ptr(v int64) (*uint8, error) {
	r, err := Int64ToUint8(v)
	return &r, err
}

// Int64ToUint16 converts int64 to uint16.
func Int64ToUint16(v int64) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// Int64ToUint16Ptr converts int64 to *uint16.
func Int64ToUint16Ptr(v int64) (*uint16, error) {
	r, err := Int64ToUint16(v)
	return &r, err
}

// Int64ToUint32 converts int64 to uint32.
func Int64ToUint32(v int64) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint32(v), nil
}

// Int64ToUint32Ptr converts int64 to *uint32.
func Int64ToUint32Ptr(v int64) (*uint32, error) {
	r, err := Int64ToUint32(v)
	return &r, err
}

// Int64ToUint64 converts int64 to uint64.
func Int64ToUint64(v int64) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint64(v), nil
}

// Int64ToUint64Ptr converts int64 to *uint64.
func Int64ToUint64Ptr(v int64) (*uint64, error) {
	r, err := Int64ToUint64(v)
	return &r, err
}

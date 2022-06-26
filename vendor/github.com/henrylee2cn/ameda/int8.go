package ameda

import (
	"strconv"
)

// Int8ToInterface converts int8 to interface.
func Int8ToInterface(v int8) interface{} {
	return v
}

// Int8ToInterfacePtr converts int8 to *interface.
func Int8ToInterfacePtr(v int8) *interface{} {
	r := Int8ToInterface(v)
	return &r
}

// Int8ToString converts int8 to string.
func Int8ToString(v int8) string {
	return strconv.FormatInt(int64(v), 10)
}

// Int8ToStringPtr converts int8 to *string.
func Int8ToStringPtr(v int8) *string {
	r := Int8ToString(v)
	return &r
}

// Int8ToBool converts int8 to bool.
func Int8ToBool(v int8) bool {
	return v != 0
}

// Int8ToBoolPtr converts int8 to *bool.
func Int8ToBoolPtr(v int8) *bool {
	r := Int8ToBool(v)
	return &r
}

// Int8ToFloat32 converts int8 to float32.
func Int8ToFloat32(v int8) float32 {
	return float32(v)
}

// Int8ToFloat32Ptr converts int8 to *float32.
func Int8ToFloat32Ptr(v int8) *float32 {
	r := Int8ToFloat32(v)
	return &r
}

// Int8ToFloat64 converts int8 to float64.
func Int8ToFloat64(v int8) float64 {
	return float64(v)
}

// Int8ToFloat64Ptr converts int8 to *float64.
func Int8ToFloat64Ptr(v int8) *float64 {
	r := Int8ToFloat64(v)
	return &r
}

// Int8ToInt converts int8 to int.
func Int8ToInt(v int8) int {
	return int(v)
}

// Int8ToIntPtr converts int8 to *int.
func Int8ToIntPtr(v int8) *int {
	r := Int8ToInt(v)
	return &r
}

// Int8ToInt8Ptr converts int8 to *int8.
func Int8ToInt8Ptr(v int8) *int8 {
	return &v
}

// Int8ToInt16 converts int8 to int16.
func Int8ToInt16(v int8) int16 {
	return int16(v)
}

// Int8ToInt16Ptr converts int8 to *int16.
func Int8ToInt16Ptr(v int8) *int16 {
	r := Int8ToInt16(v)
	return &r
}

// Int8ToInt32 converts int8 to int32.
func Int8ToInt32(v int8) int32 {
	return int32(v)
}

// Int8ToInt32Ptr converts int8 to *int32.
func Int8ToInt32Ptr(v int8) *int32 {
	r := Int8ToInt32(v)
	return &r
}

// Int8ToInt64 converts int8 to int64.
func Int8ToInt64(v int8) int64 {
	return int64(v)
}

// Int8ToInt64Ptr converts int8 to *int64.
func Int8ToInt64Ptr(v int8) *int64 {
	r := Int8ToInt64(v)
	return &r
}

// Int8ToUint converts int8 to uint.
func Int8ToUint(v int8) (uint, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint(v), nil
}

// Int8ToUintPtr converts int8 to *uint.
func Int8ToUintPtr(v int8) (*uint, error) {
	r, err := Int8ToUint(v)
	return &r, err
}

// Int8ToUint8 converts int8 to uint8.
func Int8ToUint8(v int8) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint8(v), nil
}

// Int8ToUint8Ptr converts int8 to *uint8.
func Int8ToUint8Ptr(v int8) (*uint8, error) {
	r, err := Int8ToUint8(v)
	return &r, err
}

// Int8ToUint16 converts int8 to uint16.
func Int8ToUint16(v int8) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint16(v), nil
}

// Int8ToUint16Ptr converts int8 to *uint16.
func Int8ToUint16Ptr(v int8) (*uint16, error) {
	r, err := Int8ToUint16(v)
	return &r, err
}

// Int8ToUint32 converts int8 to uint32.
func Int8ToUint32(v int8) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint32(v), nil
}

// Int8ToUint32Ptr converts int8 to *uint32.
func Int8ToUint32Ptr(v int8) (*uint32, error) {
	r, err := Int8ToUint32(v)
	return &r, err
}

// Int8ToUint64 converts int8 to uint64.
func Int8ToUint64(v int8) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint64(v), nil
}

// Int8ToUint64Ptr converts int8 to *uint64.
func Int8ToUint64Ptr(v int8) (*uint64, error) {
	r, err := Int8ToUint64(v)
	return &r, err
}

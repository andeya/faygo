package ameda

import (
	"math"
	"strconv"
)

// MaxInt returns max int number for current os.
func MaxInt() int {
	if Host32bit {
		return math.MaxInt32
	}
	return math.MaxInt64
}

// IntToInterface converts int to interface.
func IntToInterface(v int) interface{} {
	return v
}

// IntToInterfacePtr converts int to *interface.
func IntToInterfacePtr(v int) *interface{} {
	r := IntToInterface(v)
	return &r
}

// IntToString converts int to string.
func IntToString(v int) string {
	return strconv.Itoa(v)
}

// IntToStringPtr converts int to *string.
func IntToStringPtr(v int) *string {
	r := IntToString(v)
	return &r
}

// IntToBool converts int to bool.
func IntToBool(v int) bool {
	return v != 0
}

// IntToBoolPtr converts int to *bool.
func IntToBoolPtr(v int) *bool {
	r := IntToBool(v)
	return &r
}

// IntToFloat32 converts int to float32.
func IntToFloat32(v int) float32 {
	return float32(v)
}

// IntToFloat32Ptr converts int to *float32.
func IntToFloat32Ptr(v int) *float32 {
	r := IntToFloat32(v)
	return &r
}

// IntToFloat64 converts int to float64.
func IntToFloat64(v int) float64 {
	return float64(v)
}

// IntToFloat64Ptr converts int to *float64.
func IntToFloat64Ptr(v int) *float64 {
	r := IntToFloat64(v)
	return &r
}

// IntToIntPtr converts int to *int.
func IntToIntPtr(v int) *int {
	return &v
}

// IntToInt8 converts int to int8.
func IntToInt8(v int) (int8, error) {
	if v > math.MaxInt8 || v < math.MinInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// IntToInt8Ptr converts int to *int8.
func IntToInt8Ptr(v int) (*int8, error) {
	r, err := IntToInt8(v)
	return &r, err
}

// IntToInt16 converts int to int16.
func IntToInt16(v int) (int16, error) {
	if v > math.MaxInt16 || v < math.MinInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// IntToInt16Ptr converts int to *int16.
func IntToInt16Ptr(v int) (*int16, error) {
	r, err := IntToInt16(v)
	return &r, err
}

// IntToInt32 converts int to int32.
func IntToInt32(v int) (int32, error) {
	if Host64bit && (v > math.MaxInt32 || v < math.MinInt32) {
		return 0, errOverflowValue
	}
	return int32(v), nil
}

// IntToInt32Ptr converts int to *int32.
func IntToInt32Ptr(v int) (*int32, error) {
	r, err := IntToInt32(v)
	return &r, err
}

// IntToInt64 converts int to int64.
func IntToInt64(v int) int64 {
	return int64(v)
}

// IntToInt64Ptr converts int to *int64.
func IntToInt64Ptr(v int) *int64 {
	r := IntToInt64(v)
	return &r
}

// IntToUint converts int to uint.
func IntToUint(v int) (uint, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint(v), nil
}

// IntToUintPtr converts int to *uint.
func IntToUintPtr(v int) (*uint, error) {
	r, err := IntToUint(v)
	return &r, err
}

// IntToUint8 converts int to uint8.
func IntToUint8(v int) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// IntToUint8Ptr converts int to *uint8.
func IntToUint8Ptr(v int) (*uint8, error) {
	r, err := IntToUint8(v)
	return &r, err
}

// IntToUint16 converts int to uint16.
func IntToUint16(v int) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// IntToUint16Ptr converts int to *uint16.
func IntToUint16Ptr(v int) (*uint16, error) {
	r, err := IntToUint16(v)
	return &r, err
}

// IntToUint32 converts int to uint32.
func IntToUint32(v int) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if Host64bit && v > int(maxUint32) {
		return 0, errOverflowValue
	}
	return uint32(v), nil
}

// IntToUint32Ptr converts int to *uint32.
func IntToUint32Ptr(v int) (*uint32, error) {
	r, err := IntToUint32(v)
	return &r, err
}

// IntToUint64 converts int to uint64.
func IntToUint64(v int) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	return uint64(v), nil
}

// IntToUint64Ptr converts int to *uint64.
func IntToUint64Ptr(v int) (*uint64, error) {
	r, err := IntToUint64(v)
	return &r, err
}

package ameda

import (
	"strconv"
)

// BoolToInterface converts bool to interface.
func BoolToInterface(v bool) interface{} {
	return v
}

// BoolToInterfacePtr converts bool to *interface.
func BoolToInterfacePtr(v bool) *interface{} {
	r := BoolToInterface(v)
	return &r
}

// BoolToString converts bool to string.
func BoolToString(v bool) string {
	return strconv.FormatBool(v)
}

// BoolToStringPtr converts bool to *string.
func BoolToStringPtr(v bool) *string {
	r := BoolToString(v)
	return &r
}

// BoolToBoolPtr converts bool to *bool.
func BoolToBoolPtr(v bool) *bool {
	return &v
}

// BoolToFloat32 converts bool to float32.
func BoolToFloat32(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

// BoolToFloat32Ptr converts bool to *float32.
func BoolToFloat32Ptr(v bool) *float32 {
	r := BoolToFloat32(v)
	return &r
}

// BoolToFloat64 converts bool to float64.
func BoolToFloat64(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

// BoolToFloat64Ptr converts bool to *float64.
func BoolToFloat64Ptr(v bool) *float64 {
	r := BoolToFloat64(v)
	return &r
}

// BoolToInt converts bool to int.
func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// BoolToIntPtr converts bool to *int.
func BoolToIntPtr(v bool) *int {
	r := BoolToInt(v)
	return &r
}

// BoolToInt8 converts bool to int8.
func BoolToInt8(v bool) int8 {
	if v {
		return 1
	}
	return 0
}

// BoolToInt8Ptr converts bool to *int8.
func BoolToInt8Ptr(v bool) *int8 {
	r := BoolToInt8(v)
	return &r
}

// BoolToInt16 converts bool to int16.
func BoolToInt16(v bool) int16 {
	if v {
		return 1
	}
	return 0
}

// BoolToInt16Ptr converts bool to *int16.
func BoolToInt16Ptr(v bool) *int16 {
	r := BoolToInt16(v)
	return &r
}

// BoolToInt32 converts bool to int32.
func BoolToInt32(v bool) int32 {
	if v {
		return 1
	}
	return 0
}

// BoolToInt32Ptr converts bool to *int32.
func BoolToInt32Ptr(v bool) *int32 {
	r := BoolToInt32(v)
	return &r
}

// BoolToInt64 converts bool to int64.
func BoolToInt64(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

// BoolToInt64Ptr converts bool to *int64.
func BoolToInt64Ptr(v bool) *int64 {
	r := BoolToInt64(v)
	return &r
}

// BoolToUint converts bool to uint.
func BoolToUint(v bool) uint {
	if v {
		return 1
	}
	return 0
}

// BoolToUintPtr converts bool to *uint.
func BoolToUintPtr(v bool) *uint {
	r := BoolToUint(v)
	return &r
}

// BoolToUint8 converts bool to uint8.
func BoolToUint8(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}

// BoolToUint8Ptr converts bool to *uint8.
func BoolToUint8Ptr(v bool) *uint8 {
	r := BoolToUint8(v)
	return &r
}

// BoolToUint16 converts bool to uint16.
func BoolToUint16(v bool) uint16 {
	if v {
		return 1
	}
	return 0
}

// BoolToUint16Ptr converts bool to *uint16.
func BoolToUint16Ptr(v bool) *uint16 {
	r := BoolToUint16(v)
	return &r
}

// BoolToUint32 converts bool to uint32.
func BoolToUint32(v bool) uint32 {
	if v {
		return 1
	}
	return 0
}

// BoolToUint32Ptr converts bool to *uint32.
func BoolToUint32Ptr(v bool) *uint32 {
	r := BoolToUint32(v)
	return &r
}

// BoolToUint64 converts bool to uint64.
func BoolToUint64(v bool) uint64 {
	if v {
		return 1
	}
	return 0
}

// BoolToUint64Ptr converts bool to *uint64.
func BoolToUint64Ptr(v bool) *uint64 {
	r := BoolToUint64(v)
	return &r
}

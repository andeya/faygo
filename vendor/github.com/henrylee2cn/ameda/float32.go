package ameda

import (
	"fmt"
	"math"
)

// Float32ToInterface converts float32 to interface.
func Float32ToInterface(v float32) interface{} {
	return v
}

// Float32ToInterfacePtr converts float32 to *interface.
func Float32ToInterfacePtr(v float32) *interface{} {
	r := Float32ToInterface(v)
	return &r
}

// Float32ToString converts float32 to string.
func Float32ToString(v float32) string {
	return fmt.Sprintf("%f", v)
}

// Float32ToStringPtr converts float32 to *string.
func Float32ToStringPtr(v float32) *string {
	r := Float32ToString(v)
	return &r
}

// Float32ToBool converts float32 to bool.
func Float32ToBool(v float32) bool {
	return v != 0
}

// Float32ToBoolPtr converts float32 to *bool.
func Float32ToBoolPtr(v float32) *bool {
	r := Float32ToBool(v)
	return &r
}

// Float32ToFloat32Ptr converts float32 to *float32.
func Float32ToFloat32Ptr(v float32) *float32 {
	return &v
}

// Float32ToFloat64 converts float32 to float64.
func Float32ToFloat64(v float32) float64 {
	return float64(v)
}

// Float32ToFloat64Ptr converts float32 to *float64.
func Float32ToFloat64Ptr(v float32) *float64 {
	r := Float32ToFloat64(v)
	return &r
}

// Float32ToInt converts float32 to int.
func Float32ToInt(v float32) (int, error) {
	if Host64bit {
		if v > math.MaxInt64 || v < math.MinInt64 {
			return 0, errOverflowValue
		}
	} else {
		if v > math.MaxInt32 || v < math.MinInt32 {
			return 0, errOverflowValue
		}
	}
	return int(v), nil

}

// Float32ToInt8 converts float32 to int8.
func Float32ToInt8(v float32) (int8, error) {
	if v > math.MaxInt8 || v < math.MinInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Float32ToInt8Ptr converts float32 to *int8.
func Float32ToInt8Ptr(v float32) (*int8, error) {
	r, err := Float32ToInt8(v)
	return &r, err
}

// Float32ToInt16 converts float32 to int16.
func Float32ToInt16(v float32) (int16, error) {
	if v > math.MaxInt16 || v < math.MinInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// Float32ToInt16Ptr converts float32 to *int16.
func Float32ToInt16Ptr(v float32) (*int16, error) {
	r, err := Float32ToInt16(v)
	return &r, err
}

// Float32ToInt32 converts float32 to int32.
func Float32ToInt32(v float32) (int32, error) {
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, errOverflowValue
	}
	return int32(v), nil
}

// Float32ToInt32Ptr converts float32 to *int32.
func Float32ToInt32Ptr(v float32) (*int32, error) {
	r, err := Float32ToInt32(v)
	return &r, err
}

// Float32ToInt64 converts float32 to int64.
func Float32ToInt64(v float32) (int64, error) {
	if v > math.MaxInt64 || v < math.MinInt64 {
		return 0, errOverflowValue
	}
	return int64(v), nil
}

// Float32ToInt64Ptr converts float32 to *int64.
func Float32ToInt64Ptr(v float32) (*int64, error) {
	r, err := Float32ToInt64(v)
	return &r, err
}

// Float32ToUint converts float32 to uint.
func Float32ToUint(v float32) (uint, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if Host64bit {
		if v > math.MaxUint64 {
			return 0, errOverflowValue
		}
	} else {
		if v > math.MaxUint32 {
			return 0, errOverflowValue
		}
	}
	return uint(v), nil
}

// Float32ToUintPtr converts float32 to *uint.
func Float32ToUintPtr(v float32) (*uint, error) {
	r, err := Float32ToUint(v)
	return &r, err
}

// Float32ToUint8 converts float32 to uint8.
func Float32ToUint8(v float32) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Float32ToUint8Ptr converts float32 to *uint8.
func Float32ToUint8Ptr(v float32) (*uint8, error) {
	r, err := Float32ToUint8(v)
	return &r, err
}

// Float32ToUint16 converts float32 to uint16.
func Float32ToUint16(v float32) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// Float32ToUint16Ptr converts float32 to *uint16.
func Float32ToUint16Ptr(v float32) (*uint16, error) {
	r, err := Float32ToUint16(v)
	return &r, err
}

// Float32ToUint32 converts float32 to uint32.
func Float32ToUint32(v float32) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint32 {
		return 0, errOverflowValue
	}
	return uint32(v), nil
}

// Float32ToUint32Ptr converts float32 to *uint32.
func Float32ToUint32Ptr(v float32) (*uint32, error) {
	r, err := Float32ToUint32(v)
	return &r, err
}

// Float32ToUint64 converts float32 to uint64.
func Float32ToUint64(v float32) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint64 {
		return 0, errOverflowValue
	}
	return uint64(v), nil
}

// Float32ToUint64Ptr converts float32 to *uint64.
func Float32ToUint64Ptr(v float32) (*uint64, error) {
	r, err := Float32ToUint64(v)
	return &r, err
}

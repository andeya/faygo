package ameda

import (
	"fmt"
	"math"
)

// Float64ToInterface converts float64 to interface.
func Float64ToInterface(v float64) interface{} {
	return v
}

// Float64ToInterfacePtr converts float64 to *interface.
func Float64ToInterfacePtr(v float64) *interface{} {
	r := Float64ToInterface(v)
	return &r
}

// Float64ToString converts float64 to string.
func Float64ToString(v float64) string {
	return fmt.Sprintf("%f", v)
}

// Float64ToStringPtr converts float64 to *string.
func Float64ToStringPtr(v float64) *string {
	r := Float64ToString(v)
	return &r
}

// Float64ToBool converts float64 to bool.
func Float64ToBool(v float64) bool {
	return v != 0
}

// Float64ToBoolPtr converts float64 to *bool.
func Float64ToBoolPtr(v float64) *bool {
	r := Float64ToBool(v)
	return &r
}

// Float64ToFloat32 converts float64 to float32.
func Float64ToFloat32(v float64) (float32, error) {
	if v > math.MaxFloat32 || v < -math.MaxFloat32 {
		return 0, errOverflowValue
	}
	return float32(v), nil
}

// Float64ToFloat32Ptr converts float64 to *float32.
func Float64ToFloat32Ptr(v float64) (*float32, error) {
	r, err := Float64ToFloat32(v)
	return &r, err
}

// Float64ToFloat64Ptr converts float64 to *float64.
func Float64ToFloat64Ptr(v float64) *float64 {
	return &v
}

// Float64ToInt converts float64 to int.
func Float64ToInt(v float64) (int, error) {
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

// Float64ToInt8 converts float64 to int8.
func Float64ToInt8(v float64) (int8, error) {
	if v > math.MaxInt8 || v < math.MinInt8 {
		return 0, errOverflowValue
	}
	return int8(v), nil
}

// Float64ToInt8Ptr converts float64 to *int8.
func Float64ToInt8Ptr(v float64) (*int8, error) {
	r, err := Float64ToInt8(v)
	return &r, err
}

// Float64ToInt16 converts float64 to int16.
func Float64ToInt16(v float64) (int16, error) {
	if v > math.MaxInt16 || v < math.MinInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// Float64ToInt16Ptr converts float64 to *int16.
func Float64ToInt16Ptr(v float64) (*int16, error) {
	r, err := Float64ToInt16(v)
	return &r, err
}

// Float64ToInt32 converts float64 to int32.
func Float64ToInt32(v float64) (int32, error) {
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, errOverflowValue
	}
	return int32(v), nil
}

// Float64ToInt32Ptr converts float64 to *int32.
func Float64ToInt32Ptr(v float64) (*int32, error) {
	r, err := Float64ToInt32(v)
	return &r, err
}

// Float64ToInt64 converts float64 to int64.
func Float64ToInt64(v float64) (int64, error) {
	if v > math.MaxInt64 || v < math.MinInt64 {
		return 0, errOverflowValue
	}
	return int64(v), nil
}

// Float64ToInt64Ptr converts float64 to *int64.
func Float64ToInt64Ptr(v float64) (*int64, error) {
	r, err := Float64ToInt64(v)
	return &r, err
}

// Float64ToUint converts float64 to uint.
func Float64ToUint(v float64) (uint, error) {
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

// Float64ToUintPtr converts float64 to *uint.
func Float64ToUintPtr(v float64) (*uint, error) {
	r, err := Float64ToUint(v)
	return &r, err
}

// Float64ToUint8 converts float64 to uint8.
func Float64ToUint8(v float64) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// Float64ToUint8Ptr converts float64 to *uint8.
func Float64ToUint8Ptr(v float64) (*uint8, error) {
	r, err := Float64ToUint8(v)
	return &r, err
}

// Float64ToUint16 converts float64 to uint16.
func Float64ToUint16(v float64) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// Float64ToUint16Ptr converts float64 to *uint16.
func Float64ToUint16Ptr(v float64) (*uint16, error) {
	r, err := Float64ToUint16(v)
	return &r, err
}

// Float64ToUint32 converts float64 to uint32.
func Float64ToUint32(v float64) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint32 {
		return 0, errOverflowValue
	}
	return uint32(v), nil
}

// Float64ToUint32Ptr converts float64 to *uint32.
func Float64ToUint32Ptr(v float64) (*uint32, error) {
	r, err := Float64ToUint32(v)
	return &r, err
}

// Float64ToUint64 converts float64 to uint64.
func Float64ToUint64(v float64) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if v > math.MaxUint64 {
		return 0, errOverflowValue
	}
	return uint64(v), nil
}

// Float64ToUint64Ptr converts float64 to *uint64.
func Float64ToUint64Ptr(v float64) (*uint64, error) {
	r, err := Float64ToUint64(v)
	return &r, err
}

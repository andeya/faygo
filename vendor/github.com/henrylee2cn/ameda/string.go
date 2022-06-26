package ameda

import (
	"strconv"
)

// StringToInterface converts string to interface.
func StringToInterface(v string) interface{} {
	return v
}

// StringToInterfacePtr converts string to *interface.
func StringToInterfacePtr(v string) *interface{} {
	r := StringToInterface(v)
	return &r
}

// StringToStringPtr converts string to *string.
func StringToStringPtr(v string) *string {
	return &v
}

// StringToBool converts string to bool.
func StringToBool(v string, emptyAsFalse ...bool) (bool, error) {
	r, err := strconv.ParseBool(v)
	if err != nil {
		if !isEmptyAsZero(emptyAsFalse) {
			return false, err
		}
	}
	return r, nil
}

// StringToBoolPtr converts string to *bool.
func StringToBoolPtr(v string, emptyAsFalse ...bool) (*bool, error) {
	r, err := StringToBool(v, emptyAsFalse...)
	return &r, err
}

// StringToFloat32 converts string to float32.
func StringToFloat32(v string, emptyAsZero ...bool) (float32, error) {
	i, err := strconv.ParseFloat(v, 32)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return float32(i), nil
}

// StringToFloat32Ptr converts string to *float32.
func StringToFloat32Ptr(v string, emptyAsZero ...bool) (*float32, error) {
	r, err := StringToFloat32(v, emptyAsZero...)
	return &r, err
}

// StringToFloat64 converts string to float64.
func StringToFloat64(v string, emptyAsZero ...bool) (float64, error) {
	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return i, nil
}

// StringToFloat64Ptr converts string to *float64.
func StringToFloat64Ptr(v string, emptyAsZero ...bool) (*float64, error) {
	r, err := StringToFloat64(v, emptyAsZero...)
	return &r, err
}

// StringToInt converts string to int.
func StringToInt(v string, emptyAsZero ...bool) (int, error) {
	i, err := strconv.Atoi(v)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return i, nil
}

// StringToIntPtr converts string to *int.
func StringToIntPtr(v string, emptyAsZero ...bool) (*int, error) {
	r, err := StringToInt(v, emptyAsZero...)
	return &r, err
}

// StringToInt8 converts string to int8.
func StringToInt8(v string, emptyAsZero ...bool) (int8, error) {
	i, err := strconv.ParseInt(v, 10, 8)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return int8(i), nil
}

// StringToInt8Ptr converts string to *int8.
func StringToInt8Ptr(v string, emptyAsZero ...bool) (*int8, error) {
	r, err := StringToInt8(v, emptyAsZero...)
	return &r, err
}

// StringToInt16 converts string to int16.
func StringToInt16(v string, emptyAsZero ...bool) (int16, error) {
	i, err := strconv.ParseInt(v, 10, 16)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return int16(i), nil
}

// StringToInt16Ptr converts string to *int16.
func StringToInt16Ptr(v string, emptyAsZero ...bool) (*int16, error) {
	r, err := StringToInt16(v, emptyAsZero...)
	return &r, err
}

// StringToInt32 converts string to int32.
func StringToInt32(v string, emptyAsZero ...bool) (int32, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return int32(i), nil
}

// StringToInt32Ptr converts string to *int32.
func StringToInt32Ptr(v string, emptyAsZero ...bool) (*int32, error) {
	r, err := StringToInt32(v, emptyAsZero...)
	return &r, err
}

// StringToInt64 converts string to int64.
func StringToInt64(v string, emptyAsZero ...bool) (int64, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return i, nil
}

// StringToInt64Ptr converts string to *int64.
func StringToInt64Ptr(v string, emptyAsZero ...bool) (*int64, error) {
	r, err := StringToInt64(v, emptyAsZero...)
	return &r, err
}

// StringToUint converts string to uint.
func StringToUint(v string, emptyAsZero ...bool) (uint, error) {
	u, err := strconv.ParseUint(v, 10, strconv.IntSize)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return uint(u), nil
}

// StringToUintPtr converts string to *uint.
func StringToUintPtr(v string, emptyAsZero ...bool) (*uint, error) {
	r, err := StringToUint(v, emptyAsZero...)
	return &r, err
}

// StringToUint8 converts string to uint8.
func StringToUint8(v string, emptyAsZero ...bool) (uint8, error) {
	u, err := strconv.ParseUint(v, 10, 8)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return uint8(u), nil
}

// StringToUint8Ptr converts string to *uint8.
func StringToUint8Ptr(v string, emptyAsZero ...bool) (*uint8, error) {
	r, err := StringToUint8(v, emptyAsZero...)
	return &r, err
}

// StringToUint16 converts string to uint16.
func StringToUint16(v string, emptyAsZero ...bool) (uint16, error) {
	u, err := strconv.ParseUint(v, 10, 16)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return uint16(u), nil
}

// StringToUint16Ptr converts string to *uint16.
func StringToUint16Ptr(v string, emptyAsZero ...bool) (*uint16, error) {
	r, err := StringToUint16(v, emptyAsZero...)
	return &r, err
}

// StringToUint32 converts string to uint32.
func StringToUint32(v string, emptyAsZero ...bool) (uint32, error) {
	u, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return uint32(u), nil
}

// StringToUint32Ptr converts string to *uint32.
func StringToUint32Ptr(v string, emptyAsZero ...bool) (*uint32, error) {
	r, err := StringToUint32(v, emptyAsZero...)
	return &r, err
}

// StringToUint64 converts string to uint64.
func StringToUint64(v string, emptyAsZero ...bool) (uint64, error) {
	u, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		if !isEmptyAsZero(emptyAsZero) {
			return 0, err
		}
	}
	return u, nil
}

// StringToUint64Ptr converts string to *uint64.
func StringToUint64Ptr(v string, emptyAsZero ...bool) (*uint64, error) {
	r, err := StringToUint64(v, emptyAsZero...)
	return &r, err
}

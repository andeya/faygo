package ameda

import (
	"fmt"
	"reflect"
)

// InterfaceToInterfacePtr converts interface to *interface.
func InterfaceToInterfacePtr(i interface{}) *interface{} {
	return &i
}

// InterfaceToString converts interface to string.
func InterfaceToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

// InterfaceToStringPtr converts interface to *string.
func InterfaceToStringPtr(i interface{}) *string {
	v := InterfaceToString(i)
	return &v
}

// InterfaceToBool converts interface to bool.
// NOTE:
//  0 is false, other numbers are true
func InterfaceToBool(i interface{}, emptyAsFalse ...bool) (bool, error) {
	switch v := i.(type) {
	case bool:
		return v, nil
	case nil:
		return false, nil
	case float32:
		return Float32ToBool(v), nil
	case float64:
		return Float64ToBool(v), nil
	case int:
		return IntToBool(v), nil
	case int8:
		return Int8ToBool(v), nil
	case int16:
		return Int16ToBool(v), nil
	case int32:
		return Int32ToBool(v), nil
	case int64:
		return Int64ToBool(v), nil
	case uint:
		return UintToBool(v), nil
	case uint8:
		return Uint8ToBool(v), nil
	case uint16:
		return Uint16ToBool(v), nil
	case uint32:
		return Uint32ToBool(v), nil
	case uint64:
		return Uint64ToBool(v), nil
	case uintptr:
		return v != 0, nil
	case string:
		return StringToBool(v, emptyAsFalse...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return r.Bool(), nil
		case reflect.Invalid:
			return false, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToBool(r.Float()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToBool(r.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToBool(r.Uint()), nil
		case reflect.String:
			return StringToBool(r.String(), emptyAsFalse...)
		}
		if isEmptyAsZero(emptyAsFalse) {
			return !isZero(r), nil
		}
		return false, fmt.Errorf("cannot convert %#v of type %T to bool", i, i)
	}
}

// InterfaceToBoolPtr converts interface to *bool.
// NOTE:
//  0 is false, other numbers are true
func InterfaceToBoolPtr(i interface{}, emptyAsFalse ...bool) (*bool, error) {
	r, err := InterfaceToBool(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToFloat32 converts interface to float32.
func InterfaceToFloat32(i interface{}, emptyStringAsZero ...bool) (float32, error) {
	switch v := i.(type) {
	case bool:
		return BoolToFloat32(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToFloat32(v), nil
	case int8:
		return Int8ToFloat32(v), nil
	case int16:
		return Int16ToFloat32(v), nil
	case int32:
		return Int32ToFloat32(v), nil
	case int64:
		return Int64ToFloat32(v), nil
	case uint:
		return UintToFloat32(v), nil
	case uint8:
		return Uint8ToFloat32(v), nil
	case uint16:
		return Uint16ToFloat32(v), nil
	case uint32:
		return Uint32ToFloat32(v), nil
	case uint64:
		return Uint64ToFloat32(v), nil
	case uintptr:
		return UintToFloat32(uint(v)), nil
	case string:
		return StringToFloat32(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToFloat32(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32:
			return float32(r.Float()), nil
		case reflect.Float64:
			return Float64ToFloat32(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToFloat32(r.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToFloat32(r.Uint()), nil
		case reflect.String:
			return StringToFloat32(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToFloat32(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to float32", i, i)
	}
}

// InterfaceToFloat32Ptr converts interface to *float32.
func InterfaceToFloat32Ptr(i interface{}, emptyAsFalse ...bool) (*float32, error) {
	r, err := InterfaceToFloat32(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToFloat64 converts interface to float64.
func InterfaceToFloat64(i interface{}, emptyStringAsZero ...bool) (float64, error) {
	switch v := i.(type) {
	case bool:
		return BoolToFloat64(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToFloat64(v), nil
	case int8:
		return Int8ToFloat64(v), nil
	case int16:
		return Int16ToFloat64(v), nil
	case int32:
		return Int32ToFloat64(v), nil
	case int64:
		return Int64ToFloat64(v), nil
	case uint:
		return UintToFloat64(v), nil
	case uint8:
		return Uint8ToFloat64(v), nil
	case uint16:
		return Uint16ToFloat64(v), nil
	case uint32:
		return Uint32ToFloat64(v), nil
	case uint64:
		return Uint64ToFloat64(v), nil
	case uintptr:
		return UintToFloat64(uint(v)), nil
	case string:
		return StringToFloat64(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToFloat64(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return r.Float(), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToFloat64(r.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToFloat64(r.Uint()), nil
		case reflect.String:
			return StringToFloat64(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToFloat64(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to float64", i, i)
	}
}

// InterfaceToFloat64Ptr converts interface to *float64.
func InterfaceToFloat64Ptr(i interface{}, emptyAsFalse ...bool) (*float64, error) {
	r, err := InterfaceToFloat64(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToInt converts interface to int.
func InterfaceToInt(i interface{}, emptyStringAsZero ...bool) (int, error) {
	switch v := i.(type) {
	case bool:
		return BoolToInt(v), nil
	case nil:
		return 0, nil
	case int:
		return v, nil
	case int8:
		return Int8ToInt(v), nil
	case int16:
		return Int16ToInt(v), nil
	case int32:
		return Int32ToInt(v), nil
	case int64:
		return Int64ToInt(v)
	case uint:
		return UintToInt(v)
	case uint8:
		return Uint8ToInt(v), nil
	case uint16:
		return Uint16ToInt(v), nil
	case uint32:
		return Uint32ToInt(v), nil
	case uint64:
		return Uint64ToInt(v), nil
	case uintptr:
		return UintToInt(uint(v))
	case string:
		return StringToInt(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToInt(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToInt(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToInt(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToInt(r.Uint()), nil
		case reflect.String:
			return StringToInt(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToInt(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to int", i, i)
	}
}

// InterfaceToIntPtr converts interface to *float64.
func InterfaceToIntPtr(i interface{}, emptyAsFalse ...bool) (*int, error) {
	r, err := InterfaceToInt(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToInt8 converts interface to int8.
func InterfaceToInt8(i interface{}, emptyStringAsZero ...bool) (int8, error) {
	switch v := i.(type) {
	case bool:
		return BoolToInt8(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToInt8(v)
	case int8:
		return v, nil
	case int16:
		return Int16ToInt8(v)
	case int32:
		return Int32ToInt8(v)
	case int64:
		return Int64ToInt8(v)
	case uint:
		return UintToInt8(v)
	case uint8:
		return Uint8ToInt8(v)
	case uint16:
		return Uint16ToInt8(v)
	case uint32:
		return Uint32ToInt8(v)
	case uint64:
		return Uint64ToInt8(v)
	case uintptr:
		return UintToInt8(uint(v))
	case string:
		return StringToInt8(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToInt8(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToInt8(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToInt8(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToInt8(r.Uint())
		case reflect.String:
			return StringToInt8(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToInt8(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to int8", i, i)
	}
}

// InterfaceToInt8Ptr converts interface to *int8.
func InterfaceToInt8Ptr(i interface{}, emptyAsFalse ...bool) (*int8, error) {
	r, err := InterfaceToInt8(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToInt16 converts interface to int16.
func InterfaceToInt16(i interface{}, emptyStringAsZero ...bool) (int16, error) {
	switch v := i.(type) {
	case bool:
		return BoolToInt16(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToInt16(v)
	case int8:
		return Int8ToInt16(v), nil
	case int16:
		return v, nil
	case int32:
		return Int32ToInt16(v)
	case int64:
		return Int64ToInt16(v)
	case uint:
		return UintToInt16(v)
	case uint8:
		return Uint8ToInt16(v), nil
	case uint16:
		return Uint16ToInt16(v)
	case uint32:
		return Uint32ToInt16(v)
	case uint64:
		return Uint64ToInt16(v)
	case uintptr:
		return UintToInt16(uint(v))
	case string:
		return StringToInt16(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToInt16(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToInt16(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToInt16(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToInt16(r.Uint())
		case reflect.String:
			return StringToInt16(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToInt16(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to int16", i, i)
	}
}

// InterfaceToInt16Ptr converts interface to *int16.
func InterfaceToInt16Ptr(i interface{}, emptyAsFalse ...bool) (*int16, error) {
	r, err := InterfaceToInt16(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToInt32 converts interface to int32.
func InterfaceToInt32(i interface{}, emptyStringAsZero ...bool) (int32, error) {
	switch v := i.(type) {
	case bool:
		return BoolToInt32(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToInt32(v)
	case int8:
		return Int8ToInt32(v), nil
	case int16:
		return Int16ToInt32(v), nil
	case int32:
		return v, nil
	case int64:
		return Int64ToInt32(v)
	case uint:
		return UintToInt32(v)
	case uint8:
		return Uint8ToInt32(v), nil
	case uint16:
		return Uint16ToInt32(v), nil
	case uint32:
		return Uint32ToInt32(v)
	case uint64:
		return Uint64ToInt32(v)
	case uintptr:
		return UintToInt32(uint(v))
	case string:
		return StringToInt32(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToInt32(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToInt32(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToInt32(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToInt32(r.Uint())
		case reflect.String:
			return StringToInt32(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToInt32(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to int32", i, i)
	}
}

// InterfaceToInt32Ptr converts interface to *int32.
func InterfaceToInt32Ptr(i interface{}, emptyAsFalse ...bool) (*int32, error) {
	r, err := InterfaceToInt32(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToInt64 converts interface to int64.
func InterfaceToInt64(i interface{}, emptyStringAsZero ...bool) (int64, error) {
	switch v := i.(type) {
	case bool:
		return BoolToInt64(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToInt64(v), nil
	case int8:
		return Int8ToInt64(v), nil
	case int16:
		return Int16ToInt64(v), nil
	case int32:
		return Int32ToInt64(v), nil
	case int64:
		return v, nil
	case uint:
		return UintToInt64(v)
	case uint8:
		return Uint8ToInt64(v), nil
	case uint16:
		return Uint16ToInt64(v), nil
	case uint32:
		return Uint32ToInt64(v), nil
	case uint64:
		return Uint64ToInt64(v)
	case uintptr:
		return UintToInt64(uint(v))
	case string:
		return StringToInt64(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToInt64(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToInt64(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return r.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToInt64(r.Uint())
		case reflect.String:
			return StringToInt64(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToInt64(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to int64", i, i)
	}
}

// InterfaceToInt64Ptr converts interface to *int64.
func InterfaceToInt64Ptr(i interface{}, emptyAsFalse ...bool) (*int64, error) {
	r, err := InterfaceToInt64(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToUint converts interface to uint.
func InterfaceToUint(i interface{}, emptyStringAsZero ...bool) (uint, error) {
	switch v := i.(type) {
	case bool:
		return BoolToUint(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToUint(v)
	case int8:
		return Int8ToUint(v)
	case int16:
		return Int16ToUint(v)
	case int32:
		return Int32ToUint(v)
	case int64:
		return Int64ToUint(v)
	case uint:
		return v, nil
	case uint8:
		return Uint8ToUint(v), nil
	case uint16:
		return Uint16ToUint(v), nil
	case uint32:
		return Uint32ToUint(v), nil
	case uint64:
		return Uint64ToUint(v)
	case uintptr:
		return uint(v), nil
	case string:
		return StringToUint(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToUint(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToUint(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToUint(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToUint(r.Uint())
		case reflect.String:
			return StringToUint(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToUint(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint", i, i)
	}
}

// InterfaceToUintPtr converts interface to *uint.
func InterfaceToUintPtr(i interface{}, emptyAsFalse ...bool) (*uint, error) {
	r, err := InterfaceToUint(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToUint8 converts interface to uint8.
func InterfaceToUint8(i interface{}, emptyStringAsZero ...bool) (uint8, error) {
	switch v := i.(type) {
	case bool:
		return BoolToUint8(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToUint8(v)
	case int8:
		return Int8ToUint8(v)
	case int16:
		return Int16ToUint8(v)
	case int32:
		return Int32ToUint8(v)
	case int64:
		return Int64ToUint8(v)
	case uint:
		return UintToUint8(v)
	case uint8:
		return v, nil
	case uint16:
		return Uint16ToUint8(v)
	case uint32:
		return Uint32ToUint8(v)
	case uint64:
		return Uint64ToUint8(v)
	case uintptr:
		return UintToUint8(uint(v))
	case string:
		return StringToUint8(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToUint8(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToUint8(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToUint8(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToUint8(r.Uint())
		case reflect.String:
			return StringToUint8(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToUint8(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint8", i, i)
	}
}

// InterfaceToUint8Ptr converts interface to *uint8.
func InterfaceToUint8Ptr(i interface{}, emptyAsFalse ...bool) (*uint8, error) {
	r, err := InterfaceToUint8(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToUint16 converts interface to uint16.
func InterfaceToUint16(i interface{}, emptyStringAsZero ...bool) (uint16, error) {
	switch v := i.(type) {
	case bool:
		return BoolToUint16(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToUint16(v)
	case int8:
		return Int8ToUint16(v)
	case int16:
		return Int16ToUint16(v)
	case int32:
		return Int32ToUint16(v)
	case int64:
		return Int64ToUint16(v)
	case uint:
		return UintToUint16(v)
	case uint8:
		return Uint8ToUint16(v), nil
	case uint16:
		return v, nil
	case uint32:
		return Uint32ToUint16(v)
	case uint64:
		return Uint64ToUint16(v)
	case uintptr:
		return UintToUint16(uint(v))
	case string:
		return StringToUint16(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToUint16(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToUint16(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToUint16(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToUint16(r.Uint())
		case reflect.String:
			return StringToUint16(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToUint16(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint16", i, i)
	}
}

// InterfaceToUint16Ptr converts interface to *uint16.
func InterfaceToUint16Ptr(i interface{}, emptyAsFalse ...bool) (*uint16, error) {
	r, err := InterfaceToUint16(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToUint32 converts interface to uint32.
func InterfaceToUint32(i interface{}, emptyStringAsZero ...bool) (uint32, error) {
	switch v := i.(type) {
	case bool:
		return BoolToUint32(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToUint32(v)
	case int8:
		return Int8ToUint32(v)
	case int16:
		return Int16ToUint32(v)
	case int32:
		return Int32ToUint32(v)
	case int64:
		return Int64ToUint32(v)
	case uint:
		return UintToUint32(v)
	case uint8:
		return Uint8ToUint32(v), nil
	case uint16:
		return Uint16ToUint32(v), nil
	case uint32:
		return v, nil
	case uint64:
		return Uint64ToUint32(v)
	case uintptr:
		return UintToUint32(uint(v))
	case string:
		return StringToUint32(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToUint32(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToUint32(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToUint32(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Uint64ToUint32(r.Uint())
		case reflect.String:
			return StringToUint32(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToUint32(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint32", i, i)
	}
}

// InterfaceToUint32Ptr converts interface to *uint32.
func InterfaceToUint32Ptr(i interface{}, emptyAsFalse ...bool) (*uint32, error) {
	r, err := InterfaceToUint32(i, emptyAsFalse...)
	return &r, err
}

// InterfaceToUint64 converts interface to uint64.
func InterfaceToUint64(i interface{}, emptyStringAsZero ...bool) (uint64, error) {
	switch v := i.(type) {
	case bool:
		return BoolToUint64(v), nil
	case nil:
		return 0, nil
	case int:
		return IntToUint64(v)
	case int8:
		return Int8ToUint64(v)
	case int16:
		return Int16ToUint64(v)
	case int32:
		return Int32ToUint64(v)
	case int64:
		return Int64ToUint64(v)
	case uint:
		return UintToUint64(v), nil
	case uint8:
		return Uint8ToUint64(v), nil
	case uint16:
		return Uint16ToUint64(v), nil
	case uint32:
		return Uint32ToUint64(v), nil
	case uint64:
		return v, nil
	case uintptr:
		return UintToUint64(uint(v)), nil
	case string:
		return StringToUint64(v, emptyStringAsZero...)
	default:
		r := IndirectValue(reflect.ValueOf(i))
		switch r.Kind() {
		case reflect.Bool:
			return BoolToUint64(r.Bool()), nil
		case reflect.Invalid:
			return 0, nil
		case reflect.Float32, reflect.Float64:
			return Float64ToUint64(r.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Int64ToUint64(r.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return r.Uint(), nil
		case reflect.String:
			return StringToUint64(r.String(), emptyStringAsZero...)
		}
		if isEmptyAsZero(emptyStringAsZero) {
			return BoolToUint64(!isZero(r)), nil
		}
		return 0, fmt.Errorf("cannot convert %#v of type %T to uint64", i, i)
	}
}

// InterfaceToUint64Ptr converts interface to *uint64.
func InterfaceToUint64Ptr(i interface{}, emptyAsFalse ...bool) (*uint64, error) {
	r, err := InterfaceToUint64(i, emptyAsFalse...)
	return &r, err
}

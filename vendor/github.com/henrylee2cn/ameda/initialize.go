package ameda

import (
	"reflect"
)

// InitPointer initializes nil pointer with zero value.
func InitPointer(v reflect.Value) (done bool) {
	for {
		kind := v.Kind()
		if kind == reflect.Interface {
			v = v.Elem()
			continue
		}
		if kind != reflect.Ptr {
			return true
		}
		u := v.Elem()
		if u.IsValid() {
			v = u
			continue
		}
		if !v.CanSet() {
			return false
		}
		v2 := reflect.New(v.Type().Elem())
		v.Set(v2)
		v = v.Elem()
	}
}

// InitString initializes empty string pointer with def.
func InitString(p *string, def string) (done bool) {
	if p == nil {
		return false
	}
	if *p == "" {
		*p = def
	}
	return true
}

// InitBool initializes false bool pointer with def.
func InitBool(p *bool, def bool) (done bool) {
	if p == nil {
		return false
	}
	if *p == false {
		*p = def
	}
	return true
}

// InitByte initializes zero byte pointer with def.
func InitByte(p *byte, def byte) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt initializes zero int pointer with def.
func InitInt(p *int, def int) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt8 initializes zero int8 pointer with def.
func InitInt8(p *int8, def int8) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt16 initializes zero int16 pointer with def.
func InitInt16(p *int16, def int16) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt32 initializes zero int32 pointer with def.
func InitInt32(p *int32, def int32) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt64 initializes zero int64 pointer with def.
func InitInt64(p *int64, def int64) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint initializes zero uint pointer with def.
func InitUint(p *uint, def uint) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint8 initializes zero uint8 pointer with def.
func InitUint8(p *uint8, def uint8) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint16 initializes zero uint16 pointer with def.
func InitUint16(p *uint16, def uint16) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint32 initializes zero uint32 pointer with def.
func InitUint32(p *uint32, def uint32) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint64 initializes zero uint64 pointer with def.
func InitUint64(p *uint64, def uint64) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitFloat32 initializes zero float32 pointer with def.
func InitFloat32(p *float32, def float32) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitFloat64 initializes zero float64 pointer with def.
func InitFloat64(p *float64, def float64) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitSampleValue initialize the given type with some non-zero value( "?", $max_number, 0.1, true)
func InitSampleValue(t reflect.Type, maxNestingDeep int) reflect.Value {
	if maxNestingDeep <= 0 {
		maxNestingDeep = 10
	}
	ptrDepth := 0
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		ptrDepth++
	}
	v := reflect.New(t)
	v = initValue(v, 1, maxNestingDeep)
	return ReferenceValue(v, ptrDepth-1)
}

func initValue(v reflect.Value, curDeep int, maxDeep int) reflect.Value {
	InitPointer(v)
	if curDeep >= maxDeep {
		return v
	}
	var numPtr int
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
		numPtr++
	}
	switch v.Kind() {
	case reflect.Struct:
		curDeep++
		fieldNum := v.Type().NumField()
		for i := 0; i < fieldNum; i++ {
			e := v.Field(i)
			InitPointer(e)
			e.Set(initValue(e, curDeep, maxDeep))
		}
	case reflect.Slice:
		if v.Len() == 0 {
			e := reflect.New(v.Type().Elem())
			InitPointer(e)
			e = e.Elem()
			e = initValue(e, curDeep, maxDeep)
			v.Set(reflect.Append(v, e))
		}
	case reflect.Array:
		if v.Len() > 0 {
			e := reflect.New(v.Type().Elem())
			InitPointer(e)
			e = e.Elem()
			e = initValue(e, curDeep, maxDeep)
			v.Index(0).Set(reflect.Append(v, e))
		}
	case reflect.Map:
		if v.Len() == 0 {
			v.Set(reflect.MakeMap(v.Type()))
			k := reflect.New(v.Type().Key())
			InitPointer(k)
			k = k.Elem()
			k = initValue(k, curDeep, maxDeep)
			e := reflect.New(v.Type().Elem())
			InitPointer(e)
			e = e.Elem()
			e = initValue(e, curDeep, maxDeep)
			v.SetMapIndex(k, e)
		}
	case reflect.Int:
		if Host32bit {
			v.SetInt(-32)
		} else {
			v.SetInt(-64)
		}
	case reflect.Int8:
		v.SetInt(-8)
	case reflect.Int16:
		v.SetInt(-16)
	case reflect.Int32:
		v.SetInt(-32)
	case reflect.Int64:
		v.SetInt(-64)
	case reflect.Uint, reflect.Uintptr:
		if Host32bit {
			v.SetUint(32)
		} else {
			v.SetUint(64)
		}
	case reflect.Uint8:
		v.SetUint(8)
	case reflect.Uint16:
		v.SetUint(16)
	case reflect.Uint32:
		v.SetUint(32)
	case reflect.Uint64:
		v.SetUint(64)
	case reflect.Float32:
		v.SetFloat(-0.32)
	case reflect.Float64:
		v.SetFloat(-0.64)
	case reflect.Bool:
		v.SetBool(true)
	case reflect.String:
		v.SetString("?")
	default:
	}
	return ReferenceValue(v, numPtr)
}

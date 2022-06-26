package ameda

import (
	"reflect"
	"unsafe"
)

// UnsafeBytesToString convert []byte type to string type.
func UnsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeStringToBytes convert string type to []byte type.
// NOTE:
//  panic if modify the member value of the []byte.
func UnsafeStringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: sh.Len}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// IndirectValue gets the indirect value.
func IndirectValue(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}
	if v.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return v
	}
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

// DereferenceType dereference, get the underlying non-pointer type.
func DereferenceType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// DereferenceValue dereference and unpack interface,
// get the underlying non-pointer and non-interface value.
func DereferenceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DereferencePtrValue returns the underlying non-pointer type value.
func DereferencePtrValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// DereferenceInterfaceValue returns the value of the underlying type that implements the interface v.
func DereferenceInterfaceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DereferenceImplementType returns the underlying type of the value that implements the interface v.
func DereferenceImplementType(v reflect.Value) reflect.Type {
	return DereferenceType(DereferenceInterfaceValue(v).Type())
}

// DereferenceSlice convert []*T to []T.
func DereferenceSlice(v reflect.Value) reflect.Value {
	m := v.Len() - 1
	if m < 0 {
		return reflect.New(reflect.SliceOf(DereferenceType(v.Type().Elem()))).Elem()
	}
	s := make([]reflect.Value, m+1)
	for ; m >= 0; m-- {
		s[m] = DereferenceValue(v.Index(m))
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}

// ReferenceSlice convert []T to []*T, the ptrDepth is the count of '*'.
func ReferenceSlice(v reflect.Value, ptrDepth int) reflect.Value {
	if ptrDepth <= 0 {
		return v
	}
	m := v.Len() - 1
	if m < 0 {
		return reflect.New(reflect.SliceOf(ReferenceType(v.Type().Elem(), ptrDepth))).Elem()
	}
	s := make([]reflect.Value, m+1)
	for ; m >= 0; m-- {
		s[m] = ReferenceValue(v.Index(m), ptrDepth)
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}

// ReferenceType convert T to *T, the ptrDepth is the count of '*'.
func ReferenceType(t reflect.Type, ptrDepth int) reflect.Type {
	switch {
	case ptrDepth > 0:
		for ; ptrDepth > 0; ptrDepth-- {
			t = reflect.PtrTo(t)
		}
	case ptrDepth < 0:
		for ; ptrDepth < 0 && t.Kind() == reflect.Ptr; ptrDepth++ {
			t = t.Elem()
		}
	}
	return t
}

// ReferenceValue convert T to *T, the ptrDepth is the count of '*'.
func ReferenceValue(v reflect.Value, ptrDepth int) reflect.Value {
	switch {
	case ptrDepth > 0:
		for ; ptrDepth > 0; ptrDepth-- {
			vv := reflect.New(v.Type())
			vv.Elem().Set(v)
			v = vv
		}
	case ptrDepth < 0:
		for ; ptrDepth < 0 && v.Kind() == reflect.Ptr; ptrDepth++ {
			v = v.Elem()
		}
	}
	return v
}

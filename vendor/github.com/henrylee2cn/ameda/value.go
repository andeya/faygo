package ameda

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// Value go underlying type data
type Value struct {
	flag
	typPtr uintptr
	ptr    unsafe.Pointer
	_iPtr  unsafe.Pointer // avoid being GC
}

// ValueOf unpacks i to go underlying type data.
func ValueOf(i interface{}) Value {
	checkValueUsable()
	return newT(unsafe.Pointer(&i))
}

// ValueFrom gets go underlying type data from reflect.Value.
func ValueFrom(v reflect.Value) Value {
	return ValueFrom2(&v)
}

// ValueFrom2 gets go underlying type data from *reflect.Value.
func ValueFrom2(v *reflect.Value) Value {
	checkValueUsable()
	vv := newT(unsafe.Pointer(v))
	if v.CanAddr() {
		vv.flag |= flagAddr
	}
	return vv
}

//go:nocheckptr
func newT(iPtr unsafe.Pointer) Value {
	typPtr := *(*uintptr)(iPtr)
	return Value{
		typPtr: typPtr,
		flag:   getFlag(typPtr),
		ptr:    pointerElem(unsafe.Pointer(uintptr(iPtr) + ptrOffset)),
		_iPtr:  iPtr,
	}
}

// RuntimeTypeIDOf returns the underlying type ID in current runtime from interface object.
// NOTE:
//  *A and A returns the different runtime type ID;
//  It is 10 times performance of t.String().
//go:nocheckptr
func RuntimeTypeIDOf(i interface{}) uintptr {
	checkValueUsable()
	iPtr := unsafe.Pointer(&i)
	typPtr := *(*uintptr)(iPtr)
	return typPtr
}

// RuntimeTypeID returns the underlying type ID in current runtime from reflect.Type.
// NOTE:
//  *A and A returns the different runtime type ID;
//  It is 10 times performance of t.String().
//go:nocheckptr
func RuntimeTypeID(t reflect.Type) uintptr {
	checkValueUsable()
	typPtr := uintptrElem(uintptr(unsafe.Pointer(&t)) + ptrOffset)
	return typPtr
}

// RuntimeTypeID gets the underlying type ID in current runtime.
// NOTE:
//  *A and A gets the different runtime type ID;
//  It is 10 times performance of reflect.TypeOf(i).String().
//go:nocheckptr
func (v Value) RuntimeTypeID() uintptr {
	return v.typPtr
}

// Kind gets the reflect.Kind fastly.
func (v Value) Kind() reflect.Kind {
	return reflect.Kind(v.flag & flagKindMask)
}

// CanAddr reports whether the value's address can be obtained with Addr.
// Such values are called addressable. A value is addressable if it is
// an element of a slice, an element of an addressable array,
// a field of an addressable struct, or the result of dereferencing a pointer.
func (v Value) CanAddr() bool {
	return v.flag&flagAddr != 0
}

// Elem returns the Value that the interface i contains
// or that the pointer i points to.
//go:nocheckptr
func (v Value) Elem() Value {
	k := v.Kind()
	switch k {
	default:
		return v
	case reflect.Interface:
		return newT(v.ptr)
	case reflect.Ptr:
		flag2, typPtr2, has := typeUnderlying(v.flag, v.typPtr)
		if has {
			v.typPtr = typPtr2
			v.flag = flag2
			if v.Kind() == reflect.Ptr {
				v.ptr = pointerElem(v.ptr)
			}
		}
		return v
	}
}

// UnderlyingElem returns the underlying Value that the interface i contains
// or that the pointer i points to.
//go:nocheckptr
func (v Value) UnderlyingElem() Value {
	for kind := v.Kind(); kind == reflect.Ptr || kind == reflect.Interface; kind = v.Kind() {
		v = v.Elem()
	}
	return v
}

// Pointer gets the pointer of i.
// NOTE:
//  *T and T, gets diffrent pointer
//go:nocheckptr
func (v Value) Pointer() uintptr {
	switch v.Kind() {
	case reflect.Invalid:
		return 0
	case reflect.Slice:
		return uintptrElem(uintptr(v.ptr)) + sliceDataOffset
	default:
		return uintptr(v.ptr)
	}
}

// IsNil reports whether its argument i is nil.
//go:nocheckptr
func (v Value) IsNil() bool {
	return unsafe.Pointer(v.Pointer()) == nil
}

// FuncForPC returns a *Func describing the function that contains the
// given program counter address, or else nil.
//
// If pc represents multiple functions because of inlining, it returns
// the a *Func describing the innermost function, but with an entry
// of the outermost function.
//
// NOTE: Its kind must be a reflect.Func, otherwise it returns nil
//go:nocheckptr
func (v Value) FuncForPC() *runtime.Func {
	return runtime.FuncForPC(*(*uintptr)(v.ptr))
}

//go:nocheckptr
func typeUnderlying(flagVal flag, typPtr uintptr) (flag, uintptr, bool) {
	typPtr2 := uintptrElem(typPtr + elemOffset)
	if unsafe.Pointer(typPtr2) == nil {
		return flagVal, typPtr, false
	}
	tt := (*ptrType)(unsafe.Pointer(typPtr2))
	flagVal2 := flagVal&flagRO | flagIndir | flagAddr
	flagVal2 |= flag(tt.kind) & flagKindMask
	return flagVal2, typPtr2, true
}

//go:nocheckptr
func getFlag(typPtr uintptr) flag {
	if unsafe.Pointer(typPtr) == nil {
		return 0
	}
	return *(*flag)(unsafe.Pointer(typPtr + kindOffset))
}

//go:nocheckptr
func uintptrElem(ptr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr))
}

//go:nocheckptr
func pointerElem(p unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(p)
}

var errValueUsable error

func init() {
	if errValueUsable == nil {
		_, errValueUsable = checkGoVersion(runtime.Version())
	}
}

func checkGoVersion(goVer string) (string, error) {
	var rs []rune
	for _, r := range strings.TrimPrefix(goVer, "go") {
		if r >= '0' && r <= '9' || r == '.' {
			rs = append(rs, r)
		} else {
			break
		}
	}
	goVersion := strings.TrimRight(string(rs), ".")
	a, err := StringsToInts(strings.Split(goVersion, "."))
	if err != nil {
		return goVersion, err
	}
	if a[0] != 1 || a[1] < 9 {
		return goVersion, fmt.Errorf("required 1.9≤go<2.0, but current version is go" + goVersion)
	}
	return goVersion, nil
}

func checkValueUsable() {
	if errValueUsable != nil {
		panic(errValueUsable)
	}
}

var (
	e         = emptyInterface{typ: new(rtype)}
	ptrOffset = func() uintptr {
		return unsafe.Offsetof(e.word)
	}()
	kindOffset = func() uintptr {
		return unsafe.Offsetof(e.typ.kind)
	}()
	elemOffset = func() uintptr {
		return unsafe.Offsetof(new(ptrType).elem)
	}()
	sliceDataOffset = func() uintptr {
		return unsafe.Offsetof(new(reflect.SliceHeader).Data)
	}()
	// valueFlagOffset = func() uintptr {
	// 	t := reflect.TypeOf(reflect.Value{})
	// 	s, ok := t.FieldByName("flag")
	// 	if !ok {
	// 		errValueUsable = errors.New("not found reflect.Value.flag field")
	// 		return 0
	// 	}
	// 	return s.Offset
	// }()
)

// NOTE: The following definitions must be consistent with those in the standard package!!!

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)

type (
	// reflectValue struct {
	// 	typ *rtype
	// 	ptr unsafe.Pointer
	// 	flag
	// }
	emptyInterface struct {
		typ  *rtype
		word unsafe.Pointer
	}
	rtype struct {
		size       uintptr
		ptrdata    uintptr  // number of bytes in the type that can contain pointers
		hash       uint32   // hash of type; avoids computation in hash tables
		tflag      tflag    // extra type information flags
		align      uint8    // alignment of variable with this type
		fieldAlign uint8    // alignment of struct field with this type
		kind       uint8    // enumeration for C
		alg        *typeAlg // algorithm table
		gcdata     *byte    // garbage collection data
		str        nameOff  // string form
		ptrToThis  typeOff  // type for pointer to this type, may be zero
	}
	ptrType struct {
		rtype
		elem *rtype // pointer element (pointed at) type
	}
	typeAlg struct {
		hash  func(unsafe.Pointer, uintptr) uintptr
		equal func(unsafe.Pointer, unsafe.Pointer) bool
	}
	nameOff int32 // offset to a name
	typeOff int32 // offset to an *rtype
	flag    uintptr
	tflag   uint8
)

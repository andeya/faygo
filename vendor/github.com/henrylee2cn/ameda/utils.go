package ameda

import (
	"errors"
	"math"
	"reflect"
	"strconv"
)

const (
	Host64bit = strconv.IntSize == 64
	Host32bit = ^uint(0)>>32 == 0
)

var (
	errNegativeValue = errors.New("contains negative value")
	errOverflowValue = errors.New("contains overflow value")
)

var (
	maxUint32 = uint32(math.MaxUint32)
	maxInt64  = int64(math.MaxInt64)
)

func isEmptyAsZero(emptyAsZero []bool) bool {
	return len(emptyAsZero) > 0 && emptyAsZero[0]
}

func getFromIndex(length int, fromIndex ...int) int {
	if len(fromIndex) > 0 {
		return fixIndex(length, fromIndex[0], true)
	}
	return 0
}

func fixRange(length, start int, end ...int) (fixedStart, fixedEnd int, ok bool) {
	fixedStart = fixIndex(length, start, true)
	if fixedStart == length {
		return
	}
	fixedEnd = length
	if len(end) > 0 {
		fixedEnd = fixIndex(length, end[0], true)
	}
	if fixedEnd-fixedStart <= 0 {
		return
	}
	ok = true
	return
}

func fixIndex(length int, idx int, canLen bool) int {
	if idx < 0 {
		idx = length + idx
		if idx < 0 {
			return 0
		}
		return idx
	}
	if idx >= length {
		if canLen {
			return length
		}
		return length - 1
	}
	return idx
}

// isZero reports whether v is the zero value for its type.
// It panics if the argument is invalid.
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		// This should never happens, but will act as a safeguard for
		// later, as a default value doesn't makes sense here.
		panic(&reflect.ValueError{"reflect.Value.IsZero", v.Kind()})
	}
}

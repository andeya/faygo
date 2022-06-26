package ameda

// OneUint32 try to return the first element, otherwise return zero value.
func OneUint32(u []uint32) uint32 {
	if len(u) > 0 {
		return u[0]
	}
	return 0
}

// Uint32sCopy creates a copy of the uint32 slice.
func Uint32sCopy(u []uint32) []uint32 {
	b := make([]uint32, len(u))
	copy(b, u)
	return b
}

// Uint32sToInterfaces converts uint32 slice to interface slice.
func Uint32sToInterfaces(u []uint32) []interface{} {
	r := make([]interface{}, len(u))
	for k, v := range u {
		r[k] = Uint32ToInterface(v)
	}
	return r
}

// Uint32sToStrings converts uint32 slice to string slice.
func Uint32sToStrings(u []uint32) []string {
	r := make([]string, len(u))
	for k, v := range u {
		r[k] = Uint32ToString(v)
	}
	return r
}

// Uint32sToBools converts uint32 slice to bool slice.
// NOTE:
//  0 is false, everything else is true
func Uint32sToBools(u []uint32) []bool {
	r := make([]bool, len(u))
	for k, v := range u {
		r[k] = Uint32ToBool(v)
	}
	return r
}

// Uint32sToFloat32s converts uint32 slice to float32 slice.
func Uint32sToFloat32s(u []uint32) []float32 {
	r := make([]float32, len(u))
	for k, v := range u {
		r[k] = Uint32ToFloat32(v)
	}
	return r
}

// Uint32sToFloat64s converts uint32 slice to float64 slice.
func Uint32sToFloat64s(u []uint32) []float64 {
	r := make([]float64, len(u))
	for k, v := range u {
		r[k] = Uint32ToFloat64(v)
	}
	return r
}

// Uint32sToInts converts uint32 slice to int slice.
func Uint32sToInts(u []uint32) []int {
	r := make([]int, len(u))
	for k, v := range u {
		r[k] = Uint32ToInt(v)
	}
	return r
}

// Uint32sToInt8s converts uint32 slice to int8 slice.
func Uint32sToInt8s(u []uint32) ([]int8, error) {
	var err error
	r := make([]int8, len(u))
	for k, v := range u {
		r[k], err = Uint32ToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint32sToInt16s converts uint32 slice to int16 slice.
func Uint32sToInt16s(u []uint32) ([]int16, error) {
	var err error
	r := make([]int16, len(u))
	for k, v := range u {
		r[k], err = Uint32ToInt16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint32sToInt32s converts uint32 slice to int32 slice.
func Uint32sToInt32s(u []uint32) ([]int32, error) {
	var err error
	r := make([]int32, len(u))
	for k, v := range u {
		r[k], err = Uint32ToInt32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint32sToInt64s converts uint32 slice to int64 slice.
func Uint32sToInt64s(u []uint32) []int64 {
	r := make([]int64, len(u))
	for k, v := range u {
		r[k] = Uint32ToInt64(v)
	}
	return r
}

// Uint32sToUints converts uint32 slice to uint slice.
func Uint32sToUints(u []uint32) []uint {
	r := make([]uint, len(u))
	for k, v := range u {
		r[k] = Uint32ToUint(v)
	}
	return r
}

// Uint32sToUint8s converts uint32 slice to uint8 slice.
func Uint32sToUint8s(u []uint32) ([]uint8, error) {
	var err error
	r := make([]uint8, len(u))
	for k, v := range u {
		r[k], err = Uint32ToUint8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint32sToUint16s converts uint32 slice to uint16 slice.
func Uint32sToUint16s(u []uint32) ([]uint16, error) {
	var err error
	r := make([]uint16, len(u))
	for k, v := range u {
		r[k], err = Uint32ToUint16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint32sToUint64s converts uint32 slice to uint64 slice.
func Uint32sToUint64s(u []uint32) []uint64 {
	r := make([]uint64, len(u))
	for k, v := range u {
		r[k] = Uint32ToUint64(v)
	}
	return r
}

// Uint32sCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Uint32sCopyWithin(u []uint32, target, start int, end ...int) {
	target = fixIndex(len(u), target, true)
	if target == len(u) {
		return
	}
	sub := Uint32sSlice(u, start, end...)
	for k, v := range sub {
		u[target+k] = v
	}
}

// Uint32sEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func Uint32sEvery(u []uint32, fn func(u []uint32, k int, v uint32) bool) bool {
	for k, v := range u {
		if !fn(u, k, v) {
			return false
		}
	}
	return true
}

// Uint32sFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Uint32sFill(u []uint32, value uint32, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(u), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		u[k] = value
	}
}

// Uint32sFilter creates a new slice with all elements that pass the test implemented by the provided function.
func Uint32sFilter(u []uint32, fn func(u []uint32, k int, v uint32) bool) []uint32 {
	ret := make([]uint32, 0)
	for k, v := range u {
		if fn(u, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Uint32sFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func Uint32sFind(u []uint32, fn func(u []uint32, k int, v uint32) bool) (k int, v uint32) {
	for k, v := range u {
		if fn(u, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// Uint32sIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint32sIncludes(u []uint32, valueToFind uint32, fromIndex ...int) bool {
	return Uint32sIndexOf(u, valueToFind, fromIndex...) > -1
}

// Uint32sIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint32sIndexOf(u []uint32, searchElement uint32, fromIndex ...int) int {
	idx := getFromIndex(len(u), fromIndex...)
	for k, v := range u[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// Uint32sLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint32sLastIndexOf(u []uint32, searchElement uint32, fromIndex ...int) int {
	idx := getFromIndex(len(u), fromIndex...)
	for k := len(u) - 1; k >= idx; k-- {
		if searchElement == u[k] {
			return k
		}
	}
	return -1
}

// Uint32sMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func Uint32sMap(u []uint32, fn func(u []uint32, k int, v uint32) uint32) []uint32 {
	ret := make([]uint32, len(u))
	for k, v := range u {
		ret[k] = fn(u, k, v)
	}
	return ret
}

// Uint32sPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func Uint32sPop(u *[]uint32) (uint32, bool) {
	a := *u
	if len(a) == 0 {
		return 0, false
	}
	lastIndex := len(a) - 1
	last := a[lastIndex]
	a = a[:lastIndex]
	*u = a[:len(a):len(a)]
	return last, true
}

// Uint32sPush adds one or more elements to the end of an slice and returns the new length of the slice.
func Uint32sPush(u *[]uint32, element ...uint32) int {
	*u = append(*u, element...)
	return len(*u)
}

// Uint32sPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func Uint32sPushDistinct(u []uint32, element ...uint32) []uint32 {
L:
	for _, v := range element {
		for _, vv := range u {
			if vv == v {
				continue L
			}
		}
		u = append(u, v)
	}
	return u
}

// Uint32sReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Uint32sReduce(
	u []uint32,
	fn func(u []uint32, k int, v, accumulator uint32) uint32, initialValue ...uint32,
) uint32 {
	if len(u) == 0 {
		return 0
	}
	start := 0
	acc := u[start]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		start += 1
	}
	for k := start; k < len(u); k++ {
		acc = fn(u, k, u[k], acc)
	}
	return acc
}

// Uint32sReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Uint32sReduceRight(
	u []uint32,
	fn func(u []uint32, k int, v, accumulator uint32) uint32, initialValue ...uint32,
) uint32 {
	if len(u) == 0 {
		return 0
	}
	end := len(u) - 1
	acc := u[end]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		end -= 1
	}
	for k := end; k >= 0; k-- {
		acc = fn(u, k, u[k], acc)
	}
	return acc
}

// Uint32sReverse reverses an slice in place.
func Uint32sReverse(u []uint32) {
	first := 0
	last := len(u) - 1
	for first < last {
		u[first], u[last] = u[last], u[first]
		first++
		last--
	}
}

// Uint32sShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func Uint32sShift(u *[]uint32) (uint32, bool) {
	a := *u
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*u = a[:len(a):len(a)]
	return first, true
}

// Uint32sSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func Uint32sSlice(u []uint32, begin int, end ...int) []uint32 {
	fixedStart, fixedEnd, ok := fixRange(len(u), begin, end...)
	if !ok {
		return []uint32{}
	}
	return Uint32sCopy(u[fixedStart:fixedEnd])
}

// Uint32sSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func Uint32sSome(u []uint32, fn func(u []uint32, k int, v uint32) bool) bool {
	for k, v := range u {
		if fn(u, k, v) {
			return true
		}
	}
	return false
}

// Uint32sSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func Uint32sSplice(u *[]uint32, start, deleteCount int, items ...uint32) {
	a := *u
	if deleteCount < 0 {
		deleteCount = 0
	}
	start, end, _ := fixRange(len(a), start, start+1+deleteCount)
	deleteCount = end - start - 1
	for k := 0; k < len(items); k++ {
		if deleteCount > 0 {
			// replace
			a[start] = items[k]
			deleteCount--
			start++
		} else {
			// insert
			lastSlice := Uint32sCopy(a[start:])
			items = items[k:]
			a = append(a[:start], items...)
			a = append(a[:start+len(items)], lastSlice...)
			*u = a[:len(a):len(a)]
			return
		}
	}
	if deleteCount > 0 {
		a = append(a[:start], a[start+1+deleteCount:]...)
	}
	*u = a[:len(a):len(a)]
}

// Uint32sUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func Uint32sUnshift(u *[]uint32, element ...uint32) int {
	*u = append(element, *u...)
	return len(*u)
}

// Uint32sUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func Uint32sUnshiftDistinct(u *[]uint32, element ...uint32) int {
	a := *u
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[uint32]bool, len(element))
	r := make([]uint32, 0, len(a)+len(element))
L:
	for _, v := range element {
		if m[v] {
			continue
		}
		m[v] = true
		for _, vv := range a {
			if vv == v {
				continue L
			}
		}
		r = append(r, v)
	}
	r = append(r, a...)
	*u = r[:len(r):len(r)]
	return len(r)
}

// Uint32sRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func Uint32sRemoveFirst(p *[]uint32, elements ...uint32) int {
	a := *p
	m := make(map[interface{}]struct{}, len(elements))
	for _, element := range elements {
		if _, ok := m[element]; ok {
			continue
		}
		m[element] = struct{}{}
		for k, v := range a {
			if v == element {
				a = append(a[:k], a[k+1:]...)
				break
			}
		}
	}
	n := len(a)
	*p = a[:n:n]
	return n
}

// Uint32sRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func Uint32sRemoveEvery(p *[]uint32, elements ...uint32) int {
	a := *p
	m := make(map[interface{}]struct{}, len(elements))
	for _, element := range elements {
		if _, ok := m[element]; ok {
			continue
		}
		m[element] = struct{}{}
		for i := 0; i < len(a); i++ {
			if a[i] == element {
				a = append(a[:i], a[i+1:]...)
				i--
			}
		}
	}
	n := len(a)
	*p = a[:n:n]
	return n
}

// Uint32sConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func Uint32sConcat(u ...[]uint32) []uint32 {
	var totalLen int
	for _, v := range u {
		totalLen += len(v)
	}
	ret := make([]uint32, totalLen)
	dst := ret
	for _, v := range u {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// Uint32sIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func Uint32sIntersect(u ...[]uint32) (intersectCount map[uint32]int) {
	if len(u) == 0 {
		return nil
	}
	for _, v := range u {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[uint32]int, len(u))
	for k, v := range u {
		counts[k] = uint32sDistinct(v, nil)
	}
	intersectCount = counts[0]
L:
	for k, v := range intersectCount {
		for _, c := range counts[1:] {
			v2 := c[k]
			if v2 == 0 {
				delete(intersectCount, k)
				continue L
			}
			if v > v2 {
				v = v2
			}
		}
		intersectCount[k] = v
	}
	return intersectCount
}

// Uint32sDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func Uint32sDistinct(i *[]uint32, changeSlice bool) (distinctCount map[uint32]int) {
	if !changeSlice {
		return uint32sDistinct(*i, nil)
	}
	a := (*i)[:0]
	distinctCount = uint32sDistinct(*i, &a)
	n := len(distinctCount)
	*i = a[:n:n]
	return distinctCount
}

func uint32sDistinct(src []uint32, dst *[]uint32) map[uint32]int {
	m := make(map[uint32]int, len(src))
	if dst == nil {
		for _, v := range src {
			n := m[v]
			m[v] = n + 1
		}
	} else {
		a := *dst
		for _, v := range src {
			n := m[v]
			m[v] = n + 1
			if n == 0 {
				a = append(a, v)
			}
		}
		*dst = a
	}
	return m
}

// Uint32SetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint32SetUnion(set1, set2 []uint32, others ...[]uint32) []uint32 {
	m := make(map[uint32]struct{}, len(set1)+len(set2))
	r := make([]uint32, 0, len(m))
	for _, set := range append([][]uint32{set1, set2}, others...) {
		for _, v := range set {
			_, ok := m[v]
			if ok {
				continue
			}
			r = append(r, v)
			m[v] = struct{}{}
		}
	}
	return r
}

// Uint32SetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint32SetIntersect(set1, set2 []uint32, others ...[]uint32) []uint32 {
	sets := append([][]uint32{set2}, others...)
	setsCount := make([]map[uint32]int, len(sets))
	for k, v := range sets {
		setsCount[k] = uint32sDistinct(v, nil)
	}
	m := make(map[uint32]struct{}, len(set1))
	r := make([]uint32, 0, len(m))
L:
	for _, v := range set1 {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		for _, m2 := range setsCount {
			if m2[v] == 0 {
				continue L
			}
		}
		r = append(r, v)
	}
	return r
}

// Uint32SetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint32SetDifference(set1, set2 []uint32, others ...[]uint32) []uint32 {
	m := make(map[uint32]struct{}, len(set1))
	r := make([]uint32, 0, len(set1))
	sets := append([][]uint32{set2}, others...)
	for _, v := range sets {
		inter := Uint32SetIntersect(set1, v)
		for _, v := range inter {
			m[v] = struct{}{}
		}
	}
	for _, v := range set1 {
		if _, ok := m[v]; !ok {
			r = append(r, v)
			m[v] = struct{}{}
		}
	}
	return r
}

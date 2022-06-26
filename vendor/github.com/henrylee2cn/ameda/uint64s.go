package ameda

// OneUint64 try to return the first element, otherwise return zero value.
func OneUint64(u []uint64) uint64 {
	if len(u) > 0 {
		return u[0]
	}
	return 0
}

// Uint64sCopy creates a copy of the uint64 slice.
func Uint64sCopy(u []uint64) []uint64 {
	b := make([]uint64, len(u))
	copy(b, u)
	return b
}

// Uint64sToInterfaces converts uint64 slice to interface slice.
func Uint64sToInterfaces(u []uint64) []interface{} {
	r := make([]interface{}, len(u))
	for k, v := range u {
		r[k] = Uint64ToInterface(v)
	}
	return r
}

// Uint64sToStrings converts uint64 slice to string slice.
func Uint64sToStrings(u []uint64) []string {
	r := make([]string, len(u))
	for k, v := range u {
		r[k] = Uint64ToString(v)
	}
	return r
}

// Uint64sToBools converts uint64 slice to bool slice.
// NOTE:
//  0 is false, everything else is true
func Uint64sToBools(u []uint64) []bool {
	r := make([]bool, len(u))
	for k, v := range u {
		r[k] = Uint64ToBool(v)
	}
	return r
}

// Uint64sToFloat32s converts uint64 slice to float32 slice.
func Uint64sToFloat32s(u []uint64) []float32 {
	r := make([]float32, len(u))
	for k, v := range u {
		r[k] = Uint64ToFloat32(v)
	}
	return r
}

// Uint64sToFloat64s converts uint64 slice to float64 slice.
func Uint64sToFloat64s(u []uint64) []float64 {
	r := make([]float64, len(u))
	for k, v := range u {
		r[k] = Uint64ToFloat64(v)
	}
	return r
}

// Uint64sToInts converts uint64 slice to int slice.
func Uint64sToInts(u []uint64) []int {
	r := make([]int, len(u))
	for k, v := range u {
		r[k] = Uint64ToInt(v)
	}
	return r
}

// Uint64sToInt8s converts uint64 slice to int8 slice.
func Uint64sToInt8s(u []uint64) ([]int8, error) {
	var err error
	r := make([]int8, len(u))
	for k, v := range u {
		r[k], err = Uint64ToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sToInt16s converts uint64 slice to int16 slice.
func Uint64sToInt16s(u []uint64) ([]int16, error) {
	var err error
	r := make([]int16, len(u))
	for k, v := range u {
		r[k], err = Uint64ToInt16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sToInt32s converts uint64 slice to int32 slice.
func Uint64sToInt32s(u []uint64) ([]int32, error) {
	var err error
	r := make([]int32, len(u))
	for k, v := range u {
		r[k], err = Uint64ToInt32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sToInt64s converts uint64 slice to int64 slice.
func Uint64sToInt64s(u []uint64) ([]int64, error) {
	var err error
	r := make([]int64, len(u))
	for k, v := range u {
		r[k], err = Uint64ToInt64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sToUints converts uint64 slice to uint slice.
func Uint64sToUints(u []uint64) ([]uint, error) {
	var err error
	r := make([]uint, len(u))
	for k, v := range u {
		r[k], err = Uint64ToUint(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sToUint8s converts uint64 slice to uint8 slice.
func Uint64sToUint8s(u []uint64) ([]uint8, error) {
	var err error
	r := make([]uint8, len(u))
	for k, v := range u {
		r[k], err = Uint64ToUint8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sToUint16s converts uint64 slice to uint16 slice.
func Uint64sToUint16s(u []uint64) ([]uint16, error) {
	var err error
	r := make([]uint16, len(u))
	for k, v := range u {
		r[k], err = Uint64ToUint16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sToUint32s converts uint64 slice to uint32 slice.
func Uint64sToUint32s(u []uint64) ([]uint32, error) {
	var err error
	r := make([]uint32, len(u))
	for k, v := range u {
		r[k], err = Uint64ToUint32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint64sCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Uint64sCopyWithin(u []uint64, target, start int, end ...int) {
	target = fixIndex(len(u), target, true)
	if target == len(u) {
		return
	}
	sub := Uint64sSlice(u, start, end...)
	for k, v := range sub {
		u[target+k] = v
	}
}

// Uint64sEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func Uint64sEvery(u []uint64, fn func(u []uint64, k int, v uint64) bool) bool {
	for k, v := range u {
		if !fn(u, k, v) {
			return false
		}
	}
	return true
}

// Uint64sFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Uint64sFill(u []uint64, value uint64, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(u), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		u[k] = value
	}
}

// Uint64sFilter creates a new slice with all elements that pass the test implemented by the provided function.
func Uint64sFilter(u []uint64, fn func(u []uint64, k int, v uint64) bool) []uint64 {
	ret := make([]uint64, 0)
	for k, v := range u {
		if fn(u, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Uint64sFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func Uint64sFind(u []uint64, fn func(u []uint64, k int, v uint64) bool) (k int, v uint64) {
	for k, v := range u {
		if fn(u, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// Uint64sIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint64sIncludes(u []uint64, valueToFind uint64, fromIndex ...int) bool {
	return Uint64sIndexOf(u, valueToFind, fromIndex...) > -1
}

// Uint64sIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint64sIndexOf(u []uint64, searchElement uint64, fromIndex ...int) int {
	idx := getFromIndex(len(u), fromIndex...)
	for k, v := range u[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// Uint64sLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint64sLastIndexOf(u []uint64, searchElement uint64, fromIndex ...int) int {
	idx := getFromIndex(len(u), fromIndex...)
	for k := len(u) - 1; k >= idx; k-- {
		if searchElement == u[k] {
			return k
		}
	}
	return -1
}

// Uint64sMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func Uint64sMap(u []uint64, fn func(u []uint64, k int, v uint64) uint64) []uint64 {
	ret := make([]uint64, len(u))
	for k, v := range u {
		ret[k] = fn(u, k, v)
	}
	return ret
}

// Uint64sPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func Uint64sPop(u *[]uint64) (uint64, bool) {
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

// Uint64sPush adds one or more elements to the end of an slice and returns the new length of the slice.
func Uint64sPush(u *[]uint64, element ...uint64) int {
	*u = append(*u, element...)
	return len(*u)
}

// Uint64sPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func Uint64sPushDistinct(u []uint64, element ...uint64) []uint64 {
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

// Uint64sReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Uint64sReduce(
	u []uint64,
	fn func(u []uint64, k int, v, accumulator uint64) uint64, initialValue ...uint64,
) uint64 {
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

// Uint64sReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Uint64sReduceRight(
	u []uint64,
	fn func(u []uint64, k int, v, accumulator uint64) uint64, initialValue ...uint64,
) uint64 {
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

// Uint64sReverse reverses an slice in place.
func Uint64sReverse(u []uint64) {
	first := 0
	last := len(u) - 1
	for first < last {
		u[first], u[last] = u[last], u[first]
		first++
		last--
	}
}

// Uint64sShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func Uint64sShift(u *[]uint64) (uint64, bool) {
	a := *u
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*u = a[:len(a):len(a)]
	return first, true
}

// Uint64sSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func Uint64sSlice(u []uint64, begin int, end ...int) []uint64 {
	fixedStart, fixedEnd, ok := fixRange(len(u), begin, end...)
	if !ok {
		return []uint64{}
	}
	return Uint64sCopy(u[fixedStart:fixedEnd])
}

// Uint64sSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func Uint64sSome(u []uint64, fn func(u []uint64, k int, v uint64) bool) bool {
	for k, v := range u {
		if fn(u, k, v) {
			return true
		}
	}
	return false
}

// Uint64sSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func Uint64sSplice(u *[]uint64, start, deleteCount int, items ...uint64) {
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
			lastSlice := Uint64sCopy(a[start:])
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

// Uint64sUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func Uint64sUnshift(u *[]uint64, element ...uint64) int {
	*u = append(element, *u...)
	return len(*u)
}

// Uint64sUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func Uint64sUnshiftDistinct(u *[]uint64, element ...uint64) int {
	a := *u
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[uint64]bool, len(element))
	r := make([]uint64, 0, len(a)+len(element))
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

// Uint64sRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func Uint64sRemoveFirst(p *[]uint64, elements ...uint64) int {
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

// Uint64sRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func Uint64sRemoveEvery(p *[]uint64, elements ...uint64) int {
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

// Uint64sConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func Uint64sConcat(u ...[]uint64) []uint64 {
	var totalLen int
	for _, v := range u {
		totalLen += len(v)
	}
	ret := make([]uint64, totalLen)
	dst := ret
	for _, v := range u {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// Uint64sIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func Uint64sIntersect(u ...[]uint64) (intersectCount map[uint64]int) {
	if len(u) == 0 {
		return nil
	}
	for _, v := range u {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[uint64]int, len(u))
	for k, v := range u {
		counts[k] = uint64sDistinct(v, nil)
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

// Uint64sDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func Uint64sDistinct(i *[]uint64, changeSlice bool) (distinctCount map[uint64]int) {
	if !changeSlice {
		return uint64sDistinct(*i, nil)
	}
	a := (*i)[:0]
	distinctCount = uint64sDistinct(*i, &a)
	n := len(distinctCount)
	*i = a[:n:n]
	return distinctCount
}

func uint64sDistinct(src []uint64, dst *[]uint64) map[uint64]int {
	m := make(map[uint64]int, len(src))
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

// Uint64SetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint64SetUnion(set1, set2 []uint64, others ...[]uint64) []uint64 {
	m := make(map[uint64]struct{}, len(set1)+len(set2))
	r := make([]uint64, 0, len(m))
	for _, set := range append([][]uint64{set1, set2}, others...) {
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

// Uint64SetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint64SetIntersect(set1, set2 []uint64, others ...[]uint64) []uint64 {
	sets := append([][]uint64{set2}, others...)
	setsCount := make([]map[uint64]int, len(sets))
	for k, v := range sets {
		setsCount[k] = uint64sDistinct(v, nil)
	}
	m := make(map[uint64]struct{}, len(set1))
	r := make([]uint64, 0, len(m))
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

// Uint64SetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint64SetDifference(set1, set2 []uint64, others ...[]uint64) []uint64 {
	m := make(map[uint64]struct{}, len(set1))
	r := make([]uint64, 0, len(set1))
	sets := append([][]uint64{set2}, others...)
	for _, v := range sets {
		inter := Uint64SetIntersect(set1, v)
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

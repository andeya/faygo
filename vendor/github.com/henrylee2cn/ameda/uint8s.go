package ameda

// OneUint8 try to return the first element, otherwise return zero value.
func OneUint8(u []uint8) uint8 {
	if len(u) > 0 {
		return u[0]
	}
	return 0
}

// Uint8sCopy creates a copy of the uint8 slice.
func Uint8sCopy(u []uint8) []uint8 {
	b := make([]uint8, len(u))
	copy(b, u)
	return b
}

// Uint8sToInterfaces converts uint8 slice to interface slice.
func Uint8sToInterfaces(u []uint8) []interface{} {
	r := make([]interface{}, len(u))
	for k, v := range u {
		r[k] = Uint8ToInterface(v)
	}
	return r
}

// Uint8sToStrings converts uint8 slice to string slice.
func Uint8sToStrings(u []uint8) []string {
	r := make([]string, len(u))
	for k, v := range u {
		r[k] = Uint8ToString(v)
	}
	return r
}

// Uint8sToBools converts uint8 slice to bool slice.
// NOTE:
//  0 is false, everything else is true
func Uint8sToBools(u []uint8) []bool {
	r := make([]bool, len(u))
	for k, v := range u {
		r[k] = Uint8ToBool(v)
	}
	return r
}

// Uint8sToFloat32s converts uint8 slice to float32 slice.
func Uint8sToFloat32s(u []uint8) []float32 {
	r := make([]float32, len(u))
	for k, v := range u {
		r[k] = Uint8ToFloat32(v)
	}
	return r
}

// Uint8sToFloat64s converts uint8 slice to float64 slice.
func Uint8sToFloat64s(u []uint8) []float64 {
	r := make([]float64, len(u))
	for k, v := range u {
		r[k] = Uint8ToFloat64(v)
	}
	return r
}

// Uint8sToInts converts uint8 slice to int slice.
func Uint8sToInts(u []uint8) []int {
	r := make([]int, len(u))
	for k, v := range u {
		r[k] = Uint8ToInt(v)
	}
	return r
}

// Uint8sToInt8s converts uint8 slice to int8 slice.
func Uint8sToInt8s(u []uint8) ([]int8, error) {
	var err error
	r := make([]int8, len(u))
	for k, v := range u {
		r[k], err = Uint8ToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Uint8sToInt16s converts uint8 slice to int16 slice.
func Uint8sToInt16s(u []uint8) []int16 {
	r := make([]int16, len(u))
	for k, v := range u {
		r[k] = Uint8ToInt16(v)
	}
	return r
}

// Uint8sToInt32s converts uint8 slice to int32 slice.
func Uint8sToInt32s(u []uint8) []int32 {
	r := make([]int32, len(u))
	for k, v := range u {
		r[k] = Uint8ToInt32(v)
	}
	return r
}

// Uint8sToInt64s converts uint8 slice to int64 slice.
func Uint8sToInt64s(u []uint8) []int64 {
	r := make([]int64, len(u))
	for k, v := range u {
		r[k] = Uint8ToInt64(v)
	}
	return r
}

// Uint8sToUints converts uint8 slice to uint slice.
func Uint8sToUints(u []uint8) []uint {
	r := make([]uint, len(u))
	for k, v := range u {
		r[k] = Uint8ToUint(v)
	}
	return r
}

// Uint8sToUint16s converts uint8 slice to uint16 slice.
func Uint8sToUint16s(u []uint8) []uint16 {
	r := make([]uint16, len(u))
	for k, v := range u {
		r[k] = Uint8ToUint16(v)
	}
	return r
}

// Uint8sToUint32s converts uint8 slice to uint32 slice.
func Uint8sToUint32s(u []uint8) []uint32 {
	r := make([]uint32, len(u))
	for k, v := range u {
		r[k] = Uint8ToUint32(v)
	}
	return r
}

// Uint8sToUint64s converts uint8 slice to uint64 slice.
func Uint8sToUint64s(u []uint8) []uint64 {
	r := make([]uint64, len(u))
	for k, v := range u {
		r[k] = Uint8ToUint64(v)
	}
	return r
}

// Uint8sCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Uint8sCopyWithin(u []uint8, target, start int, end ...int) {
	target = fixIndex(len(u), target, true)
	if target == len(u) {
		return
	}
	sub := Uint8sSlice(u, start, end...)
	for k, v := range sub {
		u[target+k] = v
	}
}

// Uint8sEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func Uint8sEvery(u []uint8, fn func(u []uint8, k int, v uint8) bool) bool {
	for k, v := range u {
		if !fn(u, k, v) {
			return false
		}
	}
	return true
}

// Uint8sFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Uint8sFill(u []uint8, value uint8, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(u), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		u[k] = value
	}
}

// Uint8sFilter creates a new slice with all elements that pass the test implemented by the provided function.
func Uint8sFilter(u []uint8, fn func(u []uint8, k int, v uint8) bool) []uint8 {
	ret := make([]uint8, 0)
	for k, v := range u {
		if fn(u, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Uint8sFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func Uint8sFind(u []uint8, fn func(u []uint8, k int, v uint8) bool) (k int, v uint8) {
	for k, v := range u {
		if fn(u, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// Uint8sIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint8sIncludes(u []uint8, valueToFind uint8, fromIndex ...int) bool {
	return Uint8sIndexOf(u, valueToFind, fromIndex...) > -1
}

// Uint8sIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint8sIndexOf(u []uint8, searchElement uint8, fromIndex ...int) int {
	idx := getFromIndex(len(u), fromIndex...)
	for k, v := range u[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// Uint8sLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Uint8sLastIndexOf(u []uint8, searchElement uint8, fromIndex ...int) int {
	idx := getFromIndex(len(u), fromIndex...)
	for k := len(u) - 1; k >= idx; k-- {
		if searchElement == u[k] {
			return k
		}
	}
	return -1
}

// Uint8sMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func Uint8sMap(u []uint8, fn func(u []uint8, k int, v uint8) uint8) []uint8 {
	ret := make([]uint8, len(u))
	for k, v := range u {
		ret[k] = fn(u, k, v)
	}
	return ret
}

// Uint8sPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func Uint8sPop(u *[]uint8) (uint8, bool) {
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

// Uint8sPush adds one or more elements to the end of an slice and returns the new length of the slice.
func Uint8sPush(u *[]uint8, element ...uint8) int {
	*u = append(*u, element...)
	return len(*u)
}

// Uint8sPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func Uint8sPushDistinct(u []uint8, element ...uint8) []uint8 {
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

// Uint8sReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Uint8sReduce(u []uint8,
	fn func(u []uint8, k int, v, accumulator uint8) uint8, initialValue ...uint8,
) uint8 {
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

// Uint8sReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Uint8sReduceRight(u []uint8,
	fn func(u []uint8, k int, v, accumulator uint8) uint8, initialValue ...uint8,
) uint8 {
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

// Uint8sReverse reverses an slice in place.
func Uint8sReverse(u []uint8) {
	first := 0
	last := len(u) - 1
	for first < last {
		u[first], u[last] = u[last], u[first]
		first++
		last--
	}
}

// Uint8sShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func Uint8sShift(u *[]uint8) (uint8, bool) {
	a := *u
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*u = a[:len(a):len(a)]
	return first, true
}

// Uint8sSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func Uint8sSlice(u []uint8, begin int, end ...int) []uint8 {
	fixedStart, fixedEnd, ok := fixRange(len(u), begin, end...)
	if !ok {
		return []uint8{}
	}
	return Uint8sCopy(u[fixedStart:fixedEnd])
}

// Uint8sSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func Uint8sSome(u []uint8, fn func(u []uint8, k int, v uint8) bool) bool {
	for k, v := range u {
		if fn(u, k, v) {
			return true
		}
	}
	return false
}

// Uint8sSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func Uint8sSplice(u *[]uint8, start, deleteCount int, items ...uint8) {
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
			lastSlice := Uint8sCopy(a[start:])
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

// Uint8sUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func Uint8sUnshift(u *[]uint8, element ...uint8) int {
	*u = append(element, *u...)
	return len(*u)
}

// Uint8sUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func Uint8sUnshiftDistinct(u *[]uint8, element ...uint8) int {
	a := *u
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[uint8]bool, len(element))
	r := make([]uint8, 0, len(a)+len(element))
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

// Uint8sRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func Uint8sRemoveFirst(p *[]uint8, elements ...uint8) int {
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

// Uint8sRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func Uint8sRemoveEvery(p *[]uint8, elements ...uint8) int {
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

// Uint8sConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func Uint8sConcat(u ...[]uint8) []uint8 {
	var totalLen int
	for _, v := range u {
		totalLen += len(v)
	}
	ret := make([]uint8, totalLen)
	dst := ret
	for _, v := range u {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// Uint8sIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func Uint8sIntersect(u ...[]uint8) (intersectCount map[uint8]int) {
	if len(u) == 0 {
		return nil
	}
	for _, v := range u {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[uint8]int, len(u))
	for k, v := range u {
		counts[k] = uint8sDistinct(v, nil)
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

// Uint8sDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func Uint8sDistinct(i *[]uint8, changeSlice bool) (distinctCount map[uint8]int) {
	if !changeSlice {
		return uint8sDistinct(*i, nil)
	}
	a := (*i)[:0]
	distinctCount = uint8sDistinct(*i, &a)
	n := len(distinctCount)
	*i = a[:n:n]
	return distinctCount
}

func uint8sDistinct(src []uint8, dst *[]uint8) map[uint8]int {
	m := make(map[uint8]int, len(src))
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

// Uint8SetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint8SetUnion(set1, set2 []uint8, others ...[]uint8) []uint8 {
	m := make(map[uint8]struct{}, len(set1)+len(set2))
	r := make([]uint8, 0, len(m))
	for _, set := range append([][]uint8{set1, set2}, others...) {
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

// Uint8SetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint8SetIntersect(set1, set2 []uint8, others ...[]uint8) []uint8 {
	sets := append([][]uint8{set2}, others...)
	setsCount := make([]map[uint8]int, len(sets))
	for k, v := range sets {
		setsCount[k] = uint8sDistinct(v, nil)
	}
	m := make(map[uint8]struct{}, len(set1))
	r := make([]uint8, 0, len(m))
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

// Uint8SetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func Uint8SetDifference(set1, set2 []uint8, others ...[]uint8) []uint8 {
	m := make(map[uint8]struct{}, len(set1))
	r := make([]uint8, 0, len(set1))
	sets := append([][]uint8{set2}, others...)
	for _, v := range sets {
		inter := Uint8SetIntersect(set1, v)
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

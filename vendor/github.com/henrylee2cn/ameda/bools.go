package ameda

// OneBool try to return the first element, otherwise return zero value.
func OneBool(b []bool) bool {
	if len(b) > 0 {
		return b[0]
	}
	return false
}

// BoolsCopy creates a copy of the bool slice.
func BoolsCopy(b []bool) []bool {
	r := make([]bool, len(b))
	copy(r, b)
	return r
}

// BoolsToInterfaces converts int8 slice to interface slice.
func BoolsToInterfaces(b []bool) []interface{} {
	r := make([]interface{}, len(b))
	for k, v := range b {
		r[k] = v
	}
	return r
}

// BoolsToStrings converts int8 slice to string slice.
func BoolsToStrings(b []bool) []string {
	r := make([]string, len(b))
	for k, v := range b {
		r[k] = BoolToString(v)
	}
	return r
}

// BoolsToFloat32s converts int8 slice to float32 slice.
func BoolsToFloat32s(b []bool) []float32 {
	r := make([]float32, len(b))
	for k, v := range b {
		r[k] = BoolToFloat32(v)
	}
	return r
}

// BoolsToFloat64s converts int8 slice to float64 slice.
func BoolsToFloat64s(b []bool) []float64 {
	r := make([]float64, len(b))
	for k, v := range b {
		r[k] = BoolToFloat64(v)
	}
	return r
}

// BoolsToInts converts int8 slice to int slice.
func BoolsToInts(b []bool) []int {
	r := make([]int, len(b))
	for k, v := range b {
		r[k] = BoolToInt(v)
	}
	return r
}

// BoolsToInt16s converts int8 slice to int16 slice.
func BoolsToInt16s(b []bool) []int16 {
	r := make([]int16, len(b))
	for k, v := range b {
		r[k] = BoolToInt16(v)
	}
	return r
}

// BoolsToInt32s converts int8 slice to int32 slice.
func BoolsToInt32s(b []bool) []int32 {
	r := make([]int32, len(b))
	for k, v := range b {
		r[k] = BoolToInt32(v)
	}
	return r
}

// BoolsToInt64s converts int8 slice to int64 slice.
func BoolsToInt64s(b []bool) []int64 {
	r := make([]int64, len(b))
	for k, v := range b {
		r[k] = BoolToInt64(v)
	}
	return r
}

// BoolsToUints converts bool slice to uint slice.
func BoolsToUints(b []bool) []uint {
	r := make([]uint, len(b))
	for k, v := range b {
		r[k] = BoolToUint(v)
	}
	return r
}

// BoolsToUint8s converts bool slice to uint8 slice.
func BoolsToUint8s(b []bool) []uint8 {
	r := make([]uint8, len(b))
	for k, v := range b {
		r[k] = BoolToUint8(v)
	}
	return r
}

// BoolsToUint16s converts bool slice to uint16 slice.
func BoolsToUint16s(b []bool) []uint16 {
	r := make([]uint16, len(b))
	for k, v := range b {
		r[k] = BoolToUint16(v)
	}
	return r
}

// BoolsToUint32s converts bool slice to uint32 slice.
func BoolsToUint32s(b []bool) []uint32 {
	r := make([]uint32, len(b))
	for k, v := range b {
		r[k] = BoolToUint32(v)
	}
	return r
}

// BoolsToUint64s converts bool slice to uint64 slice.
func BoolsToUint64s(b []bool) []uint64 {
	r := make([]uint64, len(b))
	for k, v := range b {
		r[k] = BoolToUint64(v)
	}
	return r
}

// BoolsCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func BoolsCopyWithin(b []bool, target, start int, end ...int) {
	target = fixIndex(len(b), target, true)
	if target == len(b) {
		return
	}
	sub := BoolsSlice(b, start, end...)
	for k, v := range sub {
		b[target+k] = v
	}
}

// BoolsEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func BoolsEvery(b []bool, fn func(b []bool, k int, v bool) bool) bool {
	for k, v := range b {
		if !fn(b, k, v) {
			return false
		}
	}
	return true
}

// BoolsFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func BoolsFill(b []bool, value bool, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(b), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		b[k] = value
	}
}

// BoolsFilter creates a new slice with all elements that pass the test implemented by the provided function.
func BoolsFilter(b []bool, fn func(b []bool, k int, v bool) bool) []bool {
	ret := make([]bool, 0, 16)
	for k, v := range b {
		if fn(b, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// BoolsFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func BoolsFind(b []bool, fn func(b []bool, k int, v bool) bool) (k int, v bool) {
	for k, v := range b {
		if fn(b, k, v) {
			return k, v
		}
	}
	return -1, false
}

// BoolsIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func BoolsIncludes(b []bool, valueToFind bool, fromIndex ...int) bool {
	return BoolsIndexOf(b, valueToFind, fromIndex...) > -1
}

// BoolsIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func BoolsIndexOf(b []bool, searchElement bool, fromIndex ...int) int {
	idx := getFromIndex(len(b), fromIndex...)
	for k, v := range b[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// BoolsLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func BoolsLastIndexOf(b []bool, searchElement bool, fromIndex ...int) int {
	idx := getFromIndex(len(b), fromIndex...)
	for k := len(b) - 1; k >= idx; k-- {
		if searchElement == b[k] {
			return k
		}
	}
	return -1
}

// BoolsMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func BoolsMap(b []bool, fn func(b []bool, k int, v bool) bool) []bool {
	ret := make([]bool, len(b))
	for k, v := range b {
		ret[k] = fn(b, k, v)
	}
	return ret
}

// BoolsPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func BoolsPop(b *[]bool) (elem bool, found bool) {
	a := *b
	if len(a) == 0 {
		return false, false
	}
	lastIndex := len(a) - 1
	last := a[lastIndex]
	a = a[:lastIndex]
	*b = a[:len(a):len(a)]
	return last, true
}

// BoolsPush adds one or more elements to the end of an slice and returns the new length of the slice.
func BoolsPush(b *[]bool, element ...bool) int {
	*b = append(*b, element...)
	return len(*b)
}

// BoolsPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func BoolsPushDistinct(b []bool, element ...bool) []bool {
L:
	for _, v := range element {
		for _, vv := range b {
			if vv == v {
				continue L
			}
		}
		b = append(b, v)
	}
	return b
}

// BoolsReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func BoolsReduce(
	b []bool,
	fn func(b []bool, k int, v, accumulator bool) bool, initialValue ...bool,
) bool {
	if len(b) == 0 {
		return false
	}
	start := 0
	acc := b[start]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		start += 1
	}
	for k := start; k < len(b); k++ {
		acc = fn(b, k, b[k], acc)
	}
	return acc
}

// BoolsReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func BoolsReduceRight(
	b []bool,
	fn func(b []bool, k int, v, accumulator bool) bool, initialValue ...bool,
) bool {
	if len(b) == 0 {
		return false
	}
	end := len(b) - 1
	acc := b[end]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		end -= 1
	}
	for k := end; k >= 0; k-- {
		acc = fn(b, k, b[k], acc)
	}
	return acc
}

// BoolsReverse reverses an slice in place.
func BoolsReverse(b []bool) {
	first := 0
	last := len(b) - 1
	for first < last {
		b[first], b[last] = b[last], b[first]
		first++
		last--
	}
}

// BoolsShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func BoolsShift(b *[]bool) (element bool, found bool) {
	a := *b
	if len(a) == 0 {
		return false, false
	}
	first := a[0]
	a = a[1:]
	*b = a[:len(a):len(a)]
	return first, true
}

// BoolsSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func BoolsSlice(b []bool, begin int, end ...int) []bool {
	fixedStart, fixedEnd, ok := fixRange(len(b), begin, end...)
	if !ok {
		return []bool{}
	}
	return BoolsCopy(b[fixedStart:fixedEnd])
}

// BoolsSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func BoolsSome(b []bool, fn func(b []bool, k int, v bool) bool) bool {
	for k, v := range b {
		if fn(b, k, v) {
			return true
		}
	}
	return false
}

// BoolsSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func BoolsSplice(b *[]bool, start, deleteCount int, items ...bool) {
	a := *b
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
			lastSlice := BoolsCopy(a[start:])
			items = items[k:]
			a = append(a[:start], items...)
			a = append(a[:start+len(items)], lastSlice...)
			*b = a[:len(a):len(a)]
			return
		}
	}
	if deleteCount > 0 {
		a = append(a[:start], a[start+1+deleteCount:]...)
	}
	*b = a[:len(a):len(a)]
}

// BoolsUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func BoolsUnshift(b *[]bool, element ...bool) int {
	*b = append(element, *b...)
	return len(*b)
}

// BoolsUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func BoolsUnshiftDistinct(b *[]bool, element ...bool) int {
	a := *b
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[bool]bool, len(element))
	r := make([]bool, 0, len(a)+len(element))
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
	*b = r[:len(r):len(r)]
	return len(r)
}

// BoolsRemoveFirst removes the first matched element from the slice,
// and returns the new length of the slice.
func BoolsRemoveFirst(p *[]bool, elements ...bool) int {
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

// BoolsRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func BoolsRemoveEvery(p *[]bool, elements ...bool) int {
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

// BoolsConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func BoolsConcat(b ...[]bool) []bool {
	var totalLen int
	for _, v := range b {
		totalLen += len(v)
	}
	ret := make([]bool, totalLen)
	dst := ret
	for _, v := range b {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// BoolsIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func BoolsIntersect(b ...[]bool) (intersectCount map[bool]int) {
	if len(b) == 0 {
		return nil
	}
	for _, v := range b {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[bool]int, len(b))
	for k, v := range b {
		counts[k] = boolsDistinct(v, nil)
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

// BoolsDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func BoolsDistinct(b *[]bool, changeSlice bool) (distinctCount map[bool]int) {
	if !changeSlice {
		return boolsDistinct(*b, nil)
	}
	a := (*b)[:0]
	distinctCount = boolsDistinct(*b, &a)
	n := len(distinctCount)
	*b = a[:n:n]
	return distinctCount
}

func boolsDistinct(src []bool, dst *[]bool) map[bool]int {
	m := make(map[bool]int, len(src))
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

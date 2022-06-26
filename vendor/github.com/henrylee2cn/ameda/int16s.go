package ameda

// OneInt16 try to return the first element, otherwise return zero value.
func OneInt16(i []int16) int16 {
	if len(i) > 0 {
		return i[0]
	}
	return 0
}

// Int16sCopy creates a copy of the int16 slice.
func Int16sCopy(i []int16) []int16 {
	b := make([]int16, len(i))
	copy(b, i)
	return b
}

// Int16sToInterfaces converts int16 slice to interface slice.
func Int16sToInterfaces(i []int16) []interface{} {
	r := make([]interface{}, len(i))
	for k, v := range i {
		r[k] = v
	}
	return r
}

// Int16sToStrings converts int16 slice to string slice.
func Int16sToStrings(i []int16) []string {
	r := make([]string, len(i))
	for k, v := range i {
		r[k] = Int16ToString(v)
	}
	return r
}

// Int16sToBools converts int16 slice to bool slice.
// NOTE:
//  0 is false, everything else is true
func Int16sToBools(i []int16) []bool {
	r := make([]bool, len(i))
	for k, v := range i {
		r[k] = Int16ToBool(v)
	}
	return r
}

// Int16sToFloat32s converts int16 slice to float32 slice.
func Int16sToFloat32s(i []int16) []float32 {
	r := make([]float32, len(i))
	for k, v := range i {
		r[k] = Int16ToFloat32(v)
	}
	return r
}

// Int16sToFloat64s converts int16 slice to float64 slice.
func Int16sToFloat64s(i []int16) []float64 {
	r := make([]float64, len(i))
	for k, v := range i {
		r[k] = Int16ToFloat64(v)
	}
	return r
}

// Int16sToInts converts int16 slice to int slice.
func Int16sToInts(i []int16) []int {
	r := make([]int, len(i))
	for k, v := range i {
		r[k] = Int16ToInt(v)
	}
	return r
}

// Int16sToInt8s converts int16 slice to int8 slice.
func Int16sToInt8s(i []int16) ([]int8, error) {
	var err error
	r := make([]int8, len(i))
	for k, v := range i {
		r[k], err = Int16ToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Int16sToInt32s converts int16 slice to int32 slice.
func Int16sToInt32s(i []int16) []int32 {
	r := make([]int32, len(i))
	for k, v := range i {
		r[k] = Int16ToInt32(v)
	}
	return r
}

// Int16sToInt64s converts int16 slice to int64 slice.
func Int16sToInt64s(i []int16) []int64 {
	r := make([]int64, len(i))
	for k, v := range i {
		r[k] = Int16ToInt64(v)
	}
	return r
}

// Int16sToUints converts int16 slice to uint slice.
func Int16sToUints(i []int16) ([]uint, error) {
	var err error
	r := make([]uint, len(i))
	for k, v := range i {
		r[k], err = Int16ToUint(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Int16sToUint8s converts int16 slice to uint8 slice.
func Int16sToUint8s(i []int16) ([]uint8, error) {
	var err error
	r := make([]uint8, len(i))
	for k, v := range i {
		r[k], err = Int16ToUint8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Int16sToUint16s converts int16 slice to uint16 slice.
func Int16sToUint16s(i []int16) ([]uint16, error) {
	var err error
	r := make([]uint16, len(i))
	for k, v := range i {
		r[k], err = Int16ToUint16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Int16sToUint32s converts int16 slice to uint32 slice.
func Int16sToUint32s(i []int16) ([]uint32, error) {
	var err error
	r := make([]uint32, len(i))
	for k, v := range i {
		r[k], err = Int16ToUint32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Int16sToUint64s converts int16 slice to uint64 slice.
func Int16sToUint64s(i []int16) ([]uint64, error) {
	var err error
	r := make([]uint64, len(i))
	for k, v := range i {
		r[k], err = Int16ToUint64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Int16sCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Int16sCopyWithin(i []int16, target, start int, end ...int) {
	target = fixIndex(len(i), target, true)
	if target == len(i) {
		return
	}
	sub := Int16sSlice(i, start, end...)
	for k, v := range sub {
		i[target+k] = v
	}
}

// Int16sEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func Int16sEvery(i []int16, fn func(i []int16, k int, v int16) bool) bool {
	for k, v := range i {
		if !fn(i, k, v) {
			return false
		}
	}
	return true
}

// Int16sFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Int16sFill(i []int16, value int16, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(i), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		i[k] = value
	}
}

// Int16sFilter creates a new slice with all elements that pass the test implemented by the provided function.
func Int16sFilter(i []int16, fn func(i []int16, k int, v int16) bool) []int16 {
	ret := make([]int16, 0)
	for k, v := range i {
		if fn(i, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Int16sFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func Int16sFind(i []int16, fn func(i []int16, k int, v int16) bool) (k int, v int16) {
	for k, v := range i {
		if fn(i, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// Int16sIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Int16sIncludes(i []int16, valueToFind int16, fromIndex ...int) bool {
	return Int16sIndexOf(i, valueToFind, fromIndex...) > -1
}

// Int16sIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Int16sIndexOf(i []int16, searchElement int16, fromIndex ...int) int {
	idx := getFromIndex(len(i), fromIndex...)
	for k, v := range i[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// Int16sLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Int16sLastIndexOf(i []int16, searchElement int16, fromIndex ...int) int {
	idx := getFromIndex(len(i), fromIndex...)
	for k := len(i) - 1; k >= idx; k-- {
		if searchElement == i[k] {
			return k
		}
	}
	return -1
}

// Int16sMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func Int16sMap(i []int16, fn func(i []int16, k int, v int16) int16) []int16 {
	ret := make([]int16, len(i))
	for k, v := range i {
		ret[k] = fn(i, k, v)
	}
	return ret
}

// Int16sPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func Int16sPop(i *[]int16) (int16, bool) {
	a := *i
	if len(a) == 0 {
		return 0, false
	}
	lastIndex := len(a) - 1
	last := a[lastIndex]
	a = a[:lastIndex]
	*i = a[:len(a):len(a)]
	return last, true
}

// Int16sPush adds one or more elements to the end of an slice and returns the new length of the slice.
func Int16sPush(i *[]int16, element ...int16) int {
	*i = append(*i, element...)
	return len(*i)
}

// Int16sPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func Int16sPushDistinct(i []int16, element ...int16) []int16 {
L:
	for _, v := range element {
		for _, vv := range i {
			if vv == v {
				continue L
			}
		}
		i = append(i, v)
	}
	return i
}

// Int16sReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Int16sReduce(i []int16,
	fn func(i []int16, k int, v, accumulator int16) int16, initialValue ...int16,
) int16 {
	if len(i) == 0 {
		return 0
	}
	start := 0
	acc := i[start]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		start += 1
	}
	for k := start; k < len(i); k++ {
		acc = fn(i, k, i[k], acc)
	}
	return acc
}

// Int16sReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Int16sReduceRight(i []int16,
	fn func(i []int16, k int, v, accumulator int16) int16, initialValue ...int16,
) int16 {
	if len(i) == 0 {
		return 0
	}
	end := len(i) - 1
	acc := i[end]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		end -= 1
	}
	for k := end; k >= 0; k-- {
		acc = fn(i, k, i[k], acc)
	}
	return acc
}

// Int16sReverse reverses an slice in place.
func Int16sReverse(i []int16) {
	first := 0
	last := len(i) - 1
	for first < last {
		i[first], i[last] = i[last], i[first]
		first++
		last--
	}
}

// Int16sShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func Int16sShift(i *[]int16) (int16, bool) {
	a := *i
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*i = a[:len(a):len(a)]
	return first, true
}

// Int16sSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func Int16sSlice(i []int16, begin int, end ...int) []int16 {
	fixedStart, fixedEnd, ok := fixRange(len(i), begin, end...)
	if !ok {
		return []int16{}
	}
	return Int16sCopy(i[fixedStart:fixedEnd])
}

// Int16sSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func Int16sSome(i []int16, fn func(i []int16, k int, v int16) bool) bool {
	for k, v := range i {
		if fn(i, k, v) {
			return true
		}
	}
	return false
}

// Int16sSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func Int16sSplice(i *[]int16, start, deleteCount int, items ...int16) {
	a := *i
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
			lastSlice := Int16sCopy(a[start:])
			items = items[k:]
			a = append(a[:start], items...)
			a = append(a[:start+len(items)], lastSlice...)
			*i = a[:len(a):len(a)]
			return
		}
	}
	if deleteCount > 0 {
		a = append(a[:start], a[start+1+deleteCount:]...)
	}
	*i = a[:len(a):len(a)]
}

// Int16sUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func Int16sUnshift(i *[]int16, element ...int16) int {
	*i = append(element, *i...)
	return len(*i)
}

// Int16sUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func Int16sUnshiftDistinct(i *[]int16, element ...int16) int {
	a := *i
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[int16]bool, len(element))
	r := make([]int16, 0, len(a)+len(element))
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
	*i = r[:len(r):len(r)]
	return len(r)
}

// Int16sRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func Int16sRemoveFirst(p *[]int16, elements ...int16) int {
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

// Int16sRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func Int16sRemoveEvery(p *[]int16, elements ...int16) int {
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

// Int16sConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func Int16sConcat(i ...[]int16) []int16 {
	var totalLen int
	for _, v := range i {
		totalLen += len(v)
	}
	ret := make([]int16, totalLen)
	dst := ret
	for _, v := range i {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// Int16sIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func Int16sIntersect(i ...[]int16) (intersectCount map[int16]int) {
	if len(i) == 0 {
		return nil
	}
	for _, v := range i {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[int16]int, len(i))
	for k, v := range i {
		counts[k] = int16sDistinct(v, nil)
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

// Int16sDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func Int16sDistinct(i *[]int16, changeSlice bool) (distinctCount map[int16]int) {
	if !changeSlice {
		return int16sDistinct(*i, nil)
	}
	a := (*i)[:0]
	distinctCount = int16sDistinct(*i, &a)
	n := len(distinctCount)
	*i = a[:n:n]
	return distinctCount
}

func int16sDistinct(src []int16, dst *[]int16) map[int16]int {
	m := make(map[int16]int, len(src))
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

// Int16SetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func Int16SetUnion(set1, set2 []int16, others ...[]int16) []int16 {
	m := make(map[int16]struct{}, len(set1)+len(set2))
	r := make([]int16, 0, len(m))
	for _, set := range append([][]int16{set1, set2}, others...) {
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

// Int16SetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func Int16SetIntersect(set1, set2 []int16, others ...[]int16) []int16 {
	sets := append([][]int16{set2}, others...)
	setsCount := make([]map[int16]int, len(sets))
	for k, v := range sets {
		setsCount[k] = int16sDistinct(v, nil)
	}
	m := make(map[int16]struct{}, len(set1))
	r := make([]int16, 0, len(m))
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

// Int16SetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func Int16SetDifference(set1, set2 []int16, others ...[]int16) []int16 {
	m := make(map[int16]struct{}, len(set1))
	r := make([]int16, 0, len(set1))
	sets := append([][]int16{set2}, others...)
	for _, v := range sets {
		inter := Int16SetIntersect(set1, v)
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

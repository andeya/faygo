package ameda

// OneInt try to return the first element, otherwise return zero value.
func OneInt(i []int) int {
	if len(i) > 0 {
		return i[0]
	}
	return 0
}

// IntsCopy creates a copy of the int slice.
func IntsCopy(i []int) []int {
	b := make([]int, len(i))
	copy(b, i)
	return b
}

// IntsToInterfaces converts int slice to interface slice.
func IntsToInterfaces(i []int) []interface{} {
	r := make([]interface{}, len(i))
	for k, v := range i {
		r[k] = IntToInterface(v)
	}
	return r
}

// IntsToStrings converts int slice to string slice.
func IntsToStrings(i []int) []string {
	r := make([]string, len(i))
	for k, v := range i {
		r[k] = IntToString(v)
	}
	return r
}

// IntsToBools converts int slice to bool slice.
// NOTE:
//  0 is false, everything else is true
func IntsToBools(i []int) []bool {
	r := make([]bool, len(i))
	for k, v := range i {
		r[k] = IntToBool(v)
	}
	return r
}

// IntsToFloat32s converts int slice to float32 slice.
func IntsToFloat32s(i []int) []float32 {
	r := make([]float32, len(i))
	for k, v := range i {
		r[k] = IntToFloat32(v)
	}
	return r
}

// IntsToFloat64s converts int slice to float64 slice.
func IntsToFloat64s(i []int) []float64 {
	r := make([]float64, len(i))
	for k, v := range i {
		r[k] = IntToFloat64(v)
	}
	return r
}

// IntsToInt8s converts int slice to int8 slice.
func IntsToInt8s(i []int) ([]int8, error) {
	var err error
	r := make([]int8, len(i))
	for k, v := range i {
		r[k], err = IntToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsToInt16s converts int slice to int16 slice.
func IntsToInt16s(i []int) ([]int16, error) {
	var err error
	r := make([]int16, len(i))
	for k, v := range i {
		r[k], err = IntToInt16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsToInt32s converts int slice to int32 slice.
func IntsToInt32s(i []int) ([]int32, error) {
	var err error
	r := make([]int32, len(i))
	for k, v := range i {
		r[k], err = IntToInt32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsToInt64s converts int slice to int64 slice.
func IntsToInt64s(i []int) []int64 {
	r := make([]int64, len(i))
	for k, v := range i {
		r[k] = IntToInt64(v)
	}
	return r
}

// IntsToUints converts int slice to uint slice.
func IntsToUints(i []int) ([]uint, error) {
	var err error
	r := make([]uint, len(i))
	for k, v := range i {
		r[k], err = IntToUint(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsToUint8s converts int slice to uint8 slice.
func IntsToUint8s(i []int) ([]uint8, error) {
	var err error
	r := make([]uint8, len(i))
	for k, v := range i {
		r[k], err = IntToUint8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsToUint16s converts int slice to uint16 slice.
func IntsToUint16s(i []int) ([]uint16, error) {
	var err error
	r := make([]uint16, len(i))
	for k, v := range i {
		r[k], err = IntToUint16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsToUint32s converts int slice to uint32 slice.
func IntsToUint32s(i []int) ([]uint32, error) {
	var err error
	r := make([]uint32, len(i))
	for k, v := range i {
		r[k], err = IntToUint32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsToUint64s converts int slice to uint64 slice.
func IntsToUint64s(i []int) ([]uint64, error) {
	var err error
	r := make([]uint64, len(i))
	for k, v := range i {
		r[k], err = IntToUint64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// IntsCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func IntsCopyWithin(i []int, target, start int, end ...int) {
	target = fixIndex(len(i), target, true)
	if target == len(i) {
		return
	}
	sub := IntsSlice(i, start, end...)
	for k, v := range sub {
		i[target+k] = v
	}
}

// IntsEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func IntsEvery(i []int, fn func(i []int, k int, v int) bool) bool {
	for k, v := range i {
		if !fn(i, k, v) {
			return false
		}
	}
	return true
}

// IntsFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func IntsFill(i []int, value int, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(i), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		i[k] = value
	}
}

// IntsFilter creates a new slice with all elements that pass the test implemented by the provided function.
func IntsFilter(i []int, fn func(i []int, k int, v int) bool) []int {
	ret := make([]int, 0)
	for k, v := range i {
		if fn(i, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// IntsFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func IntsFind(i []int, fn func(i []int, k int, v int) bool) (k int, v int) {
	for k, v := range i {
		if fn(i, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// IntsIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func IntsIncludes(i []int, valueToFind int, fromIndex ...int) bool {
	return IntsIndexOf(i, valueToFind, fromIndex...) > -1
}

// IntsIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func IntsIndexOf(i []int, searchElement int, fromIndex ...int) int {
	idx := getFromIndex(len(i), fromIndex...)
	for k, v := range i[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// IntsLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func IntsLastIndexOf(i []int, searchElement int, fromIndex ...int) int {
	idx := getFromIndex(len(i), fromIndex...)
	for k := len(i) - 1; k >= idx; k-- {
		if searchElement == i[k] {
			return k
		}
	}
	return -1
}

// IntsMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func IntsMap(i []int, fn func(i []int, k int, v int) int) []int {
	ret := make([]int, len(i))
	for k, v := range i {
		ret[k] = fn(i, k, v)
	}
	return ret
}

// IntsPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func IntsPop(i *[]int) (int, bool) {
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

// IntsPush adds one or more elements to the end of an slice and returns the new length of the slice.
func IntsPush(i *[]int, element ...int) int {
	*i = append(*i, element...)
	return len(*i)
}

// IntsPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func IntsPushDistinct(i []int, element ...int) []int {
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

// IntsReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func IntsReduce(i []int, fn func(i []int, k int, v, accumulator int) int, initialValue ...int) int {
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

// IntsReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func IntsReduceRight(i []int, fn func(i []int, k int, v, accumulator int) int, initialValue ...int) int {
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

// IntsReverse reverses an slice in place.
func IntsReverse(i []int) {
	first := 0
	last := len(i) - 1
	for first < last {
		i[first], i[last] = i[last], i[first]
		first++
		last--
	}
}

// IntsShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func IntsShift(i *[]int) (int, bool) {
	a := *i
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*i = a[:len(a):len(a)]
	return first, true
}

// IntsSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func IntsSlice(i []int, begin int, end ...int) []int {
	fixedStart, fixedEnd, ok := fixRange(len(i), begin, end...)
	if !ok {
		return []int{}
	}
	return IntsCopy(i[fixedStart:fixedEnd])
}

// IntsSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func IntsSome(i []int, fn func(i []int, k int, v int) bool) bool {
	for k, v := range i {
		if fn(i, k, v) {
			return true
		}
	}
	return false
}

// IntsSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func IntsSplice(i *[]int, start, deleteCount int, items ...int) {
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
			lastSlice := IntsCopy(a[start:])
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

// IntsUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func IntsUnshift(i *[]int, element ...int) int {
	*i = append(element, *i...)
	return len(*i)
}

// IntsUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func IntsUnshiftDistinct(i *[]int, element ...int) int {
	a := *i
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[int]bool, len(element))
	r := make([]int, 0, len(a)+len(element))
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

// IntsRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func IntsRemoveFirst(p *[]int, elements ...int) int {
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

// IntsRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func IntsRemoveEvery(p *[]int, elements ...int) int {
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

// IntsConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func IntsConcat(i ...[]int) []int {
	var totalLen int
	for _, v := range i {
		totalLen += len(v)
	}
	ret := make([]int, totalLen)
	dst := ret
	for _, v := range i {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// IntsIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func IntsIntersect(i ...[]int) (intersectCount map[int]int) {
	if len(i) == 0 {
		return nil
	}
	for _, v := range i {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[int]int, len(i))
	for k, v := range i {
		counts[k] = intsDistinct(v, nil)
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

// IntsDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func IntsDistinct(i *[]int, changeSlice bool) (distinctCount map[int]int) {
	if !changeSlice {
		return intsDistinct(*i, nil)
	}
	a := (*i)[:0]
	distinctCount = intsDistinct(*i, &a)
	n := len(distinctCount)
	*i = a[:n:n]
	return distinctCount
}

func intsDistinct(src []int, dst *[]int) map[int]int {
	m := make(map[int]int, len(src))
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

// IntSetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func IntSetUnion(set1, set2 []int, others ...[]int) []int {
	m := make(map[int]struct{}, len(set1)+len(set2))
	r := make([]int, 0, len(m))
	for _, set := range append([][]int{set1, set2}, others...) {
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

// IntSetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func IntSetIntersect(set1, set2 []int, others ...[]int) []int {
	sets := append([][]int{set2}, others...)
	setsCount := make([]map[int]int, len(sets))
	for k, v := range sets {
		setsCount[k] = intsDistinct(v, nil)
	}
	m := make(map[int]struct{}, len(set1))
	r := make([]int, 0, len(m))
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

// IntSetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func IntSetDifference(set1, set2 []int, others ...[]int) []int {
	m := make(map[int]struct{}, len(set1))
	r := make([]int, 0, len(set1))
	sets := append([][]int{set2}, others...)
	for _, v := range sets {
		inter := IntSetIntersect(set1, v)
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

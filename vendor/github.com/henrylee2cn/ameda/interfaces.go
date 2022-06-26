package ameda

// OneInterface try to return the first element, otherwise return zero value.
func OneInterface(i []interface{}) interface{} {
	if len(i) > 0 {
		return i[0]
	}
	return nil
}

// InterfacesCopy creates a copy of the interface slice.
func InterfacesCopy(i []interface{}) []interface{} {
	b := make([]interface{}, len(i))
	copy(b, i)
	return b
}

// InterfacesToStrings converts interface slice to string slice.
func InterfacesToStrings(i []interface{}) []string {
	r := make([]string, len(i))
	for k, v := range i {
		r[k] = InterfaceToString(v)
	}
	return r
}

// InterfacesToBools converts interface slice to bool slice.
// NOTE:
//  0 is false, other numbers are true
func InterfacesToBools(i []interface{}) ([]bool, error) {
	var err error
	r := make([]bool, len(i))
	for k, v := range i {
		r[k], err = InterfaceToBool(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToFloat32s converts interface slice to float32 slice.
func InterfacesToFloat32s(i []interface{}) ([]float32, error) {
	var err error
	r := make([]float32, len(i))
	for k, v := range i {
		r[k], err = InterfaceToFloat32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToFloat64s converts interface slice to float64 slice.
func InterfacesToFloat64s(i []interface{}) ([]float64, error) {
	var err error
	r := make([]float64, len(i))
	for k, v := range i {
		r[k], err = InterfaceToFloat64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToInts converts interface slice to int slice.
func InterfacesToInts(i []interface{}) ([]int, error) {
	var err error
	r := make([]int, len(i))
	for k, v := range i {
		r[k], err = InterfaceToInt(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToInt8s converts interface slice to int8 slice.
func InterfacesToInt8s(i []interface{}) ([]int8, error) {
	var err error
	r := make([]int8, len(i))
	for k, v := range i {
		r[k], err = InterfaceToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToInt16s converts interface slice to int16 slice.
func InterfacesToInt16s(i []interface{}) ([]int16, error) {
	var err error
	r := make([]int16, len(i))
	for k, v := range i {
		r[k], err = InterfaceToInt16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToInt32s converts interface slice to int32 slice.
func InterfacesToInt32s(i []interface{}) ([]int32, error) {
	var err error
	r := make([]int32, len(i))
	for k, v := range i {
		r[k], err = InterfaceToInt32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToInt64s converts interface slice to int64 slice.
func InterfacesToInt64s(i []interface{}) ([]int64, error) {
	var err error
	r := make([]int64, len(i))
	for k, v := range i {
		r[k], err = InterfaceToInt64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToUints converts interface slice to uint slice.
func InterfacesToUints(i []interface{}) ([]uint, error) {
	var err error
	r := make([]uint, len(i))
	for k, v := range i {
		r[k], err = InterfaceToUint(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToUint8s converts interface slice to uint8 slice.
func InterfacesToUint8s(i []interface{}) ([]uint8, error) {
	var err error
	r := make([]uint8, len(i))
	for k, v := range i {
		r[k], err = InterfaceToUint8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToUint16s converts interface slice to uint16 slice.
func InterfacesToUint16s(i []interface{}) ([]uint16, error) {
	var err error
	r := make([]uint16, len(i))
	for k, v := range i {
		r[k], err = InterfaceToUint16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToUint32s converts interface slice to uint32 slice.
func InterfacesToUint32s(i []interface{}) ([]uint32, error) {
	var err error
	r := make([]uint32, len(i))
	for k, v := range i {
		r[k], err = InterfaceToUint32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesToUint64s converts interface slice to uint64 slice.
func InterfacesToUint64s(i []interface{}) ([]uint64, error) {
	var err error
	r := make([]uint64, len(i))
	for k, v := range i {
		r[k], err = InterfaceToUint64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// InterfacesCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func InterfacesCopyWithin(i []interface{}, target, start int, end ...int) {
	target = fixIndex(len(i), target, true)
	if target == len(i) {
		return
	}
	sub := InterfacesSlice(i, start, end...)
	for k, v := range sub {
		i[target+k] = v
	}
}

// InterfacesEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func InterfacesEvery(i []interface{}, fn func(i []interface{}, k int, v interface{}) bool) bool {
	for k, v := range i {
		if !fn(i, k, v) {
			return false
		}
	}
	return true
}

// InterfacesFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func InterfacesFill(i []interface{}, value []interface{}, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(i), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		i[k] = value
	}
}

// InterfacesFilter creates a new slice with all elements that pass the test implemented by the provided function.
func InterfacesFilter(i []interface{}, fn func(i []interface{}, k int, v interface{}) bool) []interface{} {
	ret := make([]interface{}, 0)
	for k, v := range i {
		if fn(i, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// InterfacesFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func InterfacesFind(i []interface{}, fn func(i []interface{}, k int, v interface{}) bool) (k int, v interface{}) {
	for k, v := range i {
		if fn(i, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// InterfacesIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func InterfacesIncludes(i []interface{}, valueToFind int64, fromIndex ...int) bool {
	return InterfacesIndexOf(i, valueToFind, fromIndex...) > -1
}

// InterfacesIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func InterfacesIndexOf(i []interface{}, searchElement int64, fromIndex ...int) int {
	idx := getFromIndex(len(i), fromIndex...)
	for k, v := range i[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// InterfacesLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func InterfacesLastIndexOf(i []interface{}, searchElement int64, fromIndex ...int) int {
	idx := getFromIndex(len(i), fromIndex...)
	for k := len(i) - 1; k >= idx; k-- {
		if searchElement == i[k] {
			return k
		}
	}
	return -1
}

// InterfacesMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func InterfacesMap(i []interface{}, fn func(i []interface{}, k int, v interface{}) int64) []int64 {
	ret := make([]int64, len(i))
	for k, v := range i {
		ret[k] = fn(i, k, v)
	}
	return ret
}

// InterfacesPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func InterfacesPop(i *[]interface{}) (interface{}, bool) {
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

// InterfacesPush adds one or more elements to the end of an slice and returns the new length of the slice.
func InterfacesPush(i *[]interface{}, element ...interface{}) int {
	*i = append(*i, element...)
	return len(*i)
}

// InterfacesPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func InterfacesPushDistinct(i []interface{}, element ...interface{}) []interface{} {
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

// InterfacesReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func InterfacesReduce(
	i []interface{},
	fn func(i []interface{}, k int, v, accumulator interface{}) interface{}, initialValue ...interface{},
) interface{} {
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

// InterfacesReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func InterfacesReduceRight(
	i []interface{},
	fn func(i []interface{}, k int, v, accumulator interface{}) interface{}, initialValue ...interface{},
) interface{} {
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

// InterfacesReverse reverses an slice in place.
func InterfacesReverse(i []interface{}) {
	first := 0
	last := len(i) - 1
	for first < last {
		i[first], i[last] = i[last], i[first]
		first++
		last--
	}
}

// InterfacesShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func InterfacesShift(i *[]interface{}) (interface{}, bool) {
	a := *i
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*i = a[:len(a):len(a)]
	return first, true
}

// InterfacesSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func InterfacesSlice(i []interface{}, begin int, end ...int) []interface{} {
	fixedStart, fixedEnd, ok := fixRange(len(i), begin, end...)
	if !ok {
		return []interface{}{}
	}
	return InterfacesCopy(i[fixedStart:fixedEnd])
}

// InterfacesSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func InterfacesSome(i []interface{}, fn func(i []interface{}, k int, v interface{}) bool) bool {
	for k, v := range i {
		if fn(i, k, v) {
			return true
		}
	}
	return false
}

// InterfacesSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func InterfacesSplice(i *[]interface{}, start, deleteCount int, items ...interface{}) {
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
			lastSlice := InterfacesCopy(a[start:])
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

// InterfacesUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func InterfacesUnshift(i *[]interface{}, element ...interface{}) int {
	*i = append(element, *i...)
	return len(*i)
}

// InterfacesUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func InterfacesUnshiftDistinct(i *[]interface{}, element ...interface{}) int {
	a := *i
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[interface{}]bool, len(element))
	r := make([]interface{}, 0, len(a)+len(element))
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

// InterfacesRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func InterfacesRemoveFirst(p *[]interface{}, elements ...interface{}) int {
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

// InterfacesRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func InterfacesRemoveEvery(p *[]interface{}, elements ...interface{}) int {
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

// InterfacesConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func InterfacesConcat(i ...[]interface{}) []interface{} {
	var totalLen int
	for _, v := range i {
		totalLen += len(v)
	}
	ret := make([]interface{}, totalLen)
	dst := ret
	for _, v := range i {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// InterfacesIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func InterfacesIntersect(i ...[]interface{}) (intersectCount map[interface{}]int) {
	if len(i) == 0 {
		return nil
	}
	for _, v := range i {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[interface{}]int, len(i))
	for k, v := range i {
		counts[k] = interfacesDistinct(v, nil)
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

// InterfacesDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func InterfacesDistinct(i *[]interface{}, changeSlice bool) (distinctCount map[interface{}]int) {
	if !changeSlice {
		return interfacesDistinct(*i, nil)
	}
	a := (*i)[:0]
	distinctCount = interfacesDistinct(*i, &a)
	n := len(distinctCount)
	*i = a[:n:n]
	return distinctCount
}

func interfacesDistinct(src []interface{}, dst *[]interface{}) map[interface{}]int {
	m := make(map[interface{}]int, len(src))
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

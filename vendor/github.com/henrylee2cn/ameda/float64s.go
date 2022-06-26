package ameda

// OneFloat64 try to return the first element, otherwise return zero value.
func OneFloat64(f []float64) float64 {
	if len(f) > 0 {
		return f[0]
	}
	return 0
}

// Float64sCopy creates a copy of the float64 slice.
func Float64sCopy(f []float64) []float64 {
	b := make([]float64, len(f))
	copy(b, f)
	return b
}

// Float64sToInterfaces converts float64 slice to interface slice.
func Float64sToInterfaces(f []float64) []interface{} {
	r := make([]interface{}, len(f))
	for k, v := range f {
		r[k] = Float64ToInterface(v)
	}
	return r
}

// Float64sToStrings converts float64 slice to string slice.
func Float64sToStrings(f []float64) []string {
	r := make([]string, len(f))
	for k, v := range f {
		r[k] = Float64ToString(v)
	}
	return r
}

// Float64sToBools converts float64 slice to bool slice.
// NOTE:
//  0 is false, everything else is true
func Float64sToBools(f []float64) []bool {
	r := make([]bool, len(f))
	for k, v := range f {
		r[k] = Float64ToBool(v)
	}
	return r
}

// Float64sToFloat32s converts float64 slice to float32 slice.
func Float64sToFloat32s(f []float64) ([]float32, error) {
	var err error
	r := make([]float32, len(f))
	for k, v := range f {
		r[k], err = Float64ToFloat32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToInts converts float64 slice to int slice.
func Float64sToInts(f []float64) ([]int, error) {
	var err error
	r := make([]int, len(f))
	for k, v := range f {
		r[k], err = Float64ToInt(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToInt8s converts float64 slice to int8 slice.
func Float64sToInt8s(f []float64) ([]int8, error) {
	var err error
	r := make([]int8, len(f))
	for k, v := range f {
		r[k], err = Float64ToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToInt16s converts float64 slice to int16 slice.
func Float64sToInt16s(f []float64) ([]int16, error) {
	var err error
	r := make([]int16, len(f))
	for k, v := range f {
		r[k], err = Float64ToInt16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToInt32s converts float64 slice to int32 slice.
func Float64sToInt32s(f []float64) ([]int32, error) {
	var err error
	r := make([]int32, len(f))
	for k, v := range f {
		r[k], err = Float64ToInt32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToInt64s converts float64 slice to int64 slice.
func Float64sToInt64s(f []float64) ([]int64, error) {
	var err error
	r := make([]int64, len(f))
	for k, v := range f {
		r[k], err = Float64ToInt64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToUints converts float64 slice to uint slice.
func Float64sToUints(f []float64) ([]uint, error) {
	var err error
	r := make([]uint, len(f))
	for k, v := range f {
		r[k], err = Float64ToUint(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToUint8s converts float64 slice to uint8 slice.
func Float64sToUint8s(f []float64) ([]uint8, error) {
	var err error
	r := make([]uint8, len(f))
	for k, v := range f {
		r[k], err = Float64ToUint8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToUint16s converts float64 slice to uint16 slice.
func Float64sToUint16s(f []float64) ([]uint16, error) {
	var err error
	r := make([]uint16, len(f))
	for k, v := range f {
		r[k], err = Float64ToUint16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToUint32s converts float64 slice to uint32 slice.
func Float64sToUint32s(f []float64) ([]uint32, error) {
	var err error
	r := make([]uint32, len(f))
	for k, v := range f {
		r[k], err = Float64ToUint32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sToUint64s converts float64 slice to uint64 slice.
func Float64sToUint64s(f []float64) ([]uint64, error) {
	var err error
	r := make([]uint64, len(f))
	for k, v := range f {
		r[k], err = Float64ToUint64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float64sCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Float64sCopyWithin(f []float64, target, start int, end ...int) {
	target = fixIndex(len(f), target, true)
	if target == len(f) {
		return
	}
	sub := Float64sSlice(f, start, end...)
	for k, v := range sub {
		f[target+k] = v
	}
}

// Float64sEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func Float64sEvery(f []float64, fn func(f []float64, k int, v float64) bool) bool {
	for k, v := range f {
		if !fn(f, k, v) {
			return false
		}
	}
	return true
}

// Float64sFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Float64sFill(f []float64, value float64, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(f), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		f[k] = value
	}
}

// Float64sFilter creates a new slice with all elements that pass the test implemented by the provided function.
func Float64sFilter(f []float64, fn func(f []float64, k int, v float64) bool) []float64 {
	ret := make([]float64, 0)
	for k, v := range f {
		if fn(f, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Float64sFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func Float64sFind(f []float64, fn func(f []float64, k int, v float64) bool) (k int, v float64) {
	for k, v := range f {
		if fn(f, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// Float64sIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Float64sIncludes(f []float64, valueToFind float64, fromIndex ...int) bool {
	return Float64sIndexOf(f, valueToFind, fromIndex...) > -1
}

// Float64sIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Float64sIndexOf(f []float64, searchElement float64, fromIndex ...int) int {
	idx := getFromIndex(len(f), fromIndex...)
	for k, v := range f[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// Float64sLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Float64sLastIndexOf(f []float64, searchElement float64, fromIndex ...int) int {
	idx := getFromIndex(len(f), fromIndex...)
	for k := len(f) - 1; k >= idx; k-- {
		if searchElement == f[k] {
			return k
		}
	}
	return -1
}

// Float64sMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func Float64sMap(f []float64, fn func(f []float64, k int, v float64) float64) []float64 {
	ret := make([]float64, len(f))
	for k, v := range f {
		ret[k] = fn(f, k, v)
	}
	return ret
}

// Float64sPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func Float64sPop(f *[]float64) (float64, bool) {
	a := *f
	if len(a) == 0 {
		return 0, false
	}
	lastIndex := len(a) - 1
	last := a[lastIndex]
	a = a[:lastIndex]
	*f = a[:len(a):len(a)]
	return last, true
}

// Float64sPush adds one or more elements to the end of an slice and returns the new length of the slice.
func Float64sPush(f *[]float64, element ...float64) int {
	*f = append(*f, element...)
	return len(*f)
}

// Float64sPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func Float64sPushDistinct(f []float64, element ...float64) []float64 {
L:
	for _, v := range element {
		for _, vv := range f {
			if vv == v {
				continue L
			}
		}
		f = append(f, v)
	}
	return f
}

// Float64sReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Float64sReduce(
	f []float64,
	fn func(f []float64, k int, v, accumulator float64) float64, initialValue ...float64,
) float64 {
	if len(f) == 0 {
		return 0
	}
	start := 0
	acc := f[start]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		start += 1
	}
	for k := start; k < len(f); k++ {
		acc = fn(f, k, f[k], acc)
	}
	return acc
}

// Float64sReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Float64sReduceRight(
	f []float64,
	fn func(f []float64, k int, v, accumulator float64) float64, initialValue ...float64,
) float64 {
	if len(f) == 0 {
		return 0
	}
	end := len(f) - 1
	acc := f[end]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		end -= 1
	}
	for k := end; k >= 0; k-- {
		acc = fn(f, k, f[k], acc)
	}
	return acc
}

// Float64sReverse reverses an slice in place.
func Float64sReverse(f []float64) {
	first := 0
	last := len(f) - 1
	for first < last {
		f[first], f[last] = f[last], f[first]
		first++
		last--
	}
}

// Float64sShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func Float64sShift(f *[]float64) (float64, bool) {
	a := *f
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*f = a[:len(a):len(a)]
	return first, true
}

// Float64sSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func Float64sSlice(f []float64, begin int, end ...int) []float64 {
	fixedStart, fixedEnd, ok := fixRange(len(f), begin, end...)
	if !ok {
		return []float64{}
	}
	return Float64sCopy(f[fixedStart:fixedEnd])
}

// Float64sSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func Float64sSome(f []float64, fn func(f []float64, k int, v float64) bool) bool {
	for k, v := range f {
		if fn(f, k, v) {
			return true
		}
	}
	return false
}

// Float64sSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func Float64sSplice(f *[]float64, start, deleteCount int, items ...float64) {
	a := *f
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
			lastSlice := Float64sCopy(a[start:])
			items = items[k:]
			a = append(a[:start], items...)
			a = append(a[:start+len(items)], lastSlice...)
			*f = a[:len(a):len(a)]
			return
		}
	}
	if deleteCount > 0 {
		a = append(a[:start], a[start+1+deleteCount:]...)
	}
	*f = a[:len(a):len(a)]
}

// Float64sUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func Float64sUnshift(f *[]float64, element ...float64) int {
	*f = append(element, *f...)
	return len(*f)
}

// Float64sUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func Float64sUnshiftDistinct(f *[]float64, element ...float64) int {
	a := *f
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[float64]bool, len(element))
	r := make([]float64, 0, len(a)+len(element))
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
	*f = r[:len(r):len(r)]
	return len(r)
}

// Float64sRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func Float64sRemoveFirst(p *[]float64, elements ...float64) int {
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

// Float64sRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func Float64sRemoveEvery(p *[]float64, elements ...float64) int {
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

// Float64sIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func Float64sIntersect(f ...[]float64) (intersectCount map[float64]int) {
	if len(f) == 0 {
		return nil
	}
	for _, v := range f {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[float64]int, len(f))
	for k, v := range f {
		counts[k] = float64sDistinct(v, nil)
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

// Float64sDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func Float64sDistinct(f *[]float64, changeSlice bool) (distinctCount map[float64]int) {
	if !changeSlice {
		return float64sDistinct(*f, nil)
	}
	a := (*f)[:0]
	distinctCount = float64sDistinct(*f, &a)
	n := len(distinctCount)
	*f = a[:n:n]
	return distinctCount
}

func float64sDistinct(src []float64, dst *[]float64) map[float64]int {
	m := make(map[float64]int, len(src))
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

// Float64SetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func Float64SetUnion(set1, set2 []float64, others ...[]float64) []float64 {
	m := make(map[float64]struct{}, len(set1)+len(set2))
	r := make([]float64, 0, len(m))
	for _, set := range append([][]float64{set1, set2}, others...) {
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

// Float64SetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func Float64SetIntersect(set1, set2 []float64, others ...[]float64) []float64 {
	sets := append([][]float64{set2}, others...)
	setsCount := make([]map[float64]int, len(sets))
	for k, v := range sets {
		setsCount[k] = float64sDistinct(v, nil)
	}
	m := make(map[float64]struct{}, len(set1))
	r := make([]float64, 0, len(m))
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

// Float64SetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func Float64SetDifference(set1, set2 []float64, others ...[]float64) []float64 {
	m := make(map[float64]struct{}, len(set1))
	r := make([]float64, 0, len(set1))
	sets := append([][]float64{set2}, others...)
	for _, v := range sets {
		inter := Float64SetIntersect(set1, v)
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

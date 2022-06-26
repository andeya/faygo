package ameda

// OneFloat32 try to return the first element, otherwise return zero value.
func OneFloat32(f []float32) float32 {
	if len(f) > 0 {
		return f[0]
	}
	return 0
}

// Float32sCopy creates a copy of the float32 slice.
func Float32sCopy(f []float32) []float32 {
	b := make([]float32, len(f))
	copy(b, f)
	return b
}

// Float32sToInterfaces converts float32 slice to interface slice.
func Float32sToInterfaces(f []float32) []interface{} {
	r := make([]interface{}, len(f))
	for k, v := range f {
		r[k] = Float32ToInterface(v)
	}
	return r
}

// Float32sToStrings converts float32 slice to string slice.
func Float32sToStrings(f []float32) []string {
	r := make([]string, len(f))
	for k, v := range f {
		r[k] = Float32ToString(v)
	}
	return r
}

// Float32sToBools converts float32 slice to bool slice.
// NOTE:
//  0 is false, everything else is true
func Float32sToBools(f []float32) []bool {
	r := make([]bool, len(f))
	for k, v := range f {
		r[k] = Float32ToBool(v)
	}
	return r
}

// Float32sToFloat64s converts float32 slice to float64 slice.
func Float32sToFloat64s(f []float32) []float64 {
	r := make([]float64, len(f))
	for k, v := range f {
		r[k] = Float32ToFloat64(v)
	}
	return r
}

// Float32sToInts converts float32 slice to int slice.
func Float32sToInts(f []float32) ([]int, error) {
	var err error
	r := make([]int, len(f))
	for k, v := range f {
		r[k], err = Float32ToInt(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToInt8s converts float32 slice to int8 slice.
func Float32sToInt8s(f []float32) ([]int8, error) {
	var err error
	r := make([]int8, len(f))
	for k, v := range f {
		r[k], err = Float32ToInt8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToInt16s converts float32 slice to int16 slice.
func Float32sToInt16s(f []float32) ([]int16, error) {
	var err error
	r := make([]int16, len(f))
	for k, v := range f {
		r[k], err = Float32ToInt16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToInt32s converts float32 slice to int32 slice.
func Float32sToInt32s(f []float32) ([]int32, error) {
	var err error
	r := make([]int32, len(f))
	for k, v := range f {
		r[k], err = Float32ToInt32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToInt64s converts float32 slice to int64 slice.
func Float32sToInt64s(f []float32) ([]int64, error) {
	var err error
	r := make([]int64, len(f))
	for k, v := range f {
		r[k], err = Float32ToInt64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToUints converts float32 slice to uint slice.
func Float32sToUints(f []float32) ([]uint, error) {
	var err error
	r := make([]uint, len(f))
	for k, v := range f {
		r[k], err = Float32ToUint(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToUint8s converts float32 slice to uint8 slice.
func Float32sToUint8s(f []float32) ([]uint8, error) {
	var err error
	r := make([]uint8, len(f))
	for k, v := range f {
		r[k], err = Float32ToUint8(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToUint16s converts float32 slice to uint16 slice.
func Float32sToUint16s(f []float32) ([]uint16, error) {
	var err error
	r := make([]uint16, len(f))
	for k, v := range f {
		r[k], err = Float32ToUint16(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToUint32s converts float32 slice to uint32 slice.
func Float32sToUint32s(f []float32) ([]uint32, error) {
	var err error
	r := make([]uint32, len(f))
	for k, v := range f {
		r[k], err = Float32ToUint32(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sToUint64s converts float32 slice to uint64 slice.
func Float32sToUint64s(f []float32) ([]uint64, error) {
	var err error
	r := make([]uint64, len(f))
	for k, v := range f {
		r[k], err = Float32ToUint64(v)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// Float32sCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Float32sCopyWithin(f []float32, target, start int, end ...int) {
	target = fixIndex(len(f), target, true)
	if target == len(f) {
		return
	}
	sub := Float32sSlice(f, start, end...)
	for k, v := range sub {
		f[target+k] = v
	}
}

// Float32sEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func Float32sEvery(f []float32, fn func(f []float32, k int, v float32) bool) bool {
	for k, v := range f {
		if !fn(f, k, v) {
			return false
		}
	}
	return true
}

// Float32sFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Float32sFill(f []float32, value float32, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(f), start, end...)
	if !ok {
		return
	}
	for k := fixedStart; k < fixedEnd; k++ {
		f[k] = value
	}
}

// Float32sFilter creates a new slice with all elements that pass the test implemented by the provided function.
func Float32sFilter(f []float32, fn func(f []float32, k int, v float32) bool) []float32 {
	ret := make([]float32, 0)
	for k, v := range f {
		if fn(f, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Float32sFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func Float32sFind(f []float32, fn func(f []float32, k int, v float32) bool) (k int, v float32) {
	for k, v := range f {
		if fn(f, k, v) {
			return k, v
		}
	}
	return -1, 0
}

// Float32sIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Float32sIncludes(f []float32, valueToFind float32, fromIndex ...int) bool {
	return Float32sIndexOf(f, valueToFind, fromIndex...) > -1
}

// Float32sIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Float32sIndexOf(f []float32, searchElement float32, fromIndex ...int) int {
	idx := getFromIndex(len(f), fromIndex...)
	for k, v := range f[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// Float32sLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func Float32sLastIndexOf(f []float32, searchElement float32, fromIndex ...int) int {
	idx := getFromIndex(len(f), fromIndex...)
	for k := len(f) - 1; k >= idx; k-- {
		if searchElement == f[k] {
			return k
		}
	}
	return -1
}

// Float32sMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func Float32sMap(f []float32, fn func(f []float32, k int, v float32) float32) []float32 {
	ret := make([]float32, len(f))
	for k, v := range f {
		ret[k] = fn(f, k, v)
	}
	return ret
}

// Float32sPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func Float32sPop(f *[]float32) (float32, bool) {
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

// Float32sPush adds one or more elements to the end of an slice and returns the new length of the slice.
func Float32sPush(f *[]float32, element ...float32) int {
	*f = append(*f, element...)
	return len(*f)
}

// Float32sPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func Float32sPushDistinct(f []float32, element ...float32) []float32 {
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

// Float32sReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Float32sReduce(
	f []float32,
	fn func(f []float32, k int, v, accumulator float32) float32, initialValue ...float32,
) float32 {
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

// Float32sReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func Float32sReduceRight(
	f []float32,
	fn func(f []float32, k int, v, accumulator float32) float32, initialValue ...float32,
) float32 {
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

// Float32sReverse reverses an slice in place.
func Float32sReverse(f []float32) {
	first := 0
	last := len(f) - 1
	for first < last {
		f[first], f[last] = f[last], f[first]
		first++
		last--
	}
}

// Float32sShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func Float32sShift(f *[]float32) (float32, bool) {
	a := *f
	if len(a) == 0 {
		return 0, false
	}
	first := a[0]
	a = a[1:]
	*f = a[:len(a):len(a)]
	return first, true
}

// Float32sSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func Float32sSlice(f []float32, begin int, end ...int) []float32 {
	fixedStart, fixedEnd, ok := fixRange(len(f), begin, end...)
	if !ok {
		return []float32{}
	}
	return Float32sCopy(f[fixedStart:fixedEnd])
}

// Float32sSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func Float32sSome(f []float32, fn func(f []float32, k int, v float32) bool) bool {
	for k, v := range f {
		if fn(f, k, v) {
			return true
		}
	}
	return false
}

// Float32sSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func Float32sSplice(f *[]float32, start, deleteCount int, items ...float32) {
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
			lastSlice := Float32sCopy(a[start:])
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

// Float32sUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func Float32sUnshift(f *[]float32, element ...float32) int {
	*f = append(element, *f...)
	return len(*f)
}

// Float32sUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func Float32sUnshiftDistinct(f *[]float32, element ...float32) int {
	a := *f
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[float32]bool, len(element))
	r := make([]float32, 0, len(a)+len(element))
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

// Float32sRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func Float32sRemoveFirst(p *[]float32, elements ...float32) int {
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

// Float32sRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func Float32sRemoveEvery(p *[]float32, elements ...float32) int {
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

// Float32sIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func Float32sIntersect(f ...[]float32) (intersectCount map[float32]int) {
	if len(f) == 0 {
		return nil
	}
	for _, v := range f {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[float32]int, len(f))
	for k, v := range f {
		counts[k] = float32sDistinct(v, nil)
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

// Float32sDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func Float32sDistinct(f *[]float32, changeSlice bool) (distinctCount map[float32]int) {
	if !changeSlice {
		return float32sDistinct(*f, nil)
	}
	a := (*f)[:0]
	distinctCount = float32sDistinct(*f, &a)
	n := len(distinctCount)
	*f = a[:n:n]
	return distinctCount
}

func float32sDistinct(src []float32, dst *[]float32) map[float32]int {
	m := make(map[float32]int, len(src))
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

// Float32SetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func Float32SetUnion(set1, set2 []float32, others ...[]float32) []float32 {
	m := make(map[float32]struct{}, len(set1)+len(set2))
	r := make([]float32, 0, len(m))
	for _, set := range append([][]float32{set1, set2}, others...) {
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

// Float32SetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func Float32SetIntersect(set1, set2 []float32, others ...[]float32) []float32 {
	sets := append([][]float32{set2}, others...)
	setsCount := make([]map[float32]int, len(sets))
	for k, v := range sets {
		setsCount[k] = float32sDistinct(v, nil)
	}
	m := make(map[float32]struct{}, len(set1))
	r := make([]float32, 0, len(m))
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

// Float32SetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func Float32SetDifference(set1, set2 []float32, others ...[]float32) []float32 {
	m := make(map[float32]struct{}, len(set1))
	r := make([]float32, 0, len(set1))
	sets := append([][]float32{set2}, others...)
	for _, v := range sets {
		inter := Float32SetIntersect(set1, v)
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

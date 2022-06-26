package ameda

import (
	"strings"
)

// OneString try to return the first element, otherwise return zero value.
func OneString(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

// StringsCopy creates a copy of the string slice.
func StringsCopy(s []string) []string {
	b := make([]string, len(s))
	copy(b, s)
	return b
}

// StringsToInterfaces converts string slice to interface slice.
func StringsToInterfaces(s []string) []interface{} {
	r := make([]interface{}, len(s))
	for k, v := range s {
		r[k] = StringToInterface(v)
	}
	return r
}

// StringsToBools converts string slice to bool slice.
func StringsToBools(s []string, emptyAsZero ...bool) ([]bool, error) {
	var err error
	r := make([]bool, len(s))
	for k, v := range s {
		r[k], err = StringToBool(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToFloat32s converts string slice to float32 slice.
func StringsToFloat32s(s []string, emptyAsZero ...bool) ([]float32, error) {
	var err error
	r := make([]float32, len(s))
	for k, v := range s {
		r[k], err = StringToFloat32(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToFloat64s converts string slice to float64 slice.
func StringsToFloat64s(s []string, emptyAsZero ...bool) ([]float64, error) {
	var err error
	r := make([]float64, len(s))
	for k, v := range s {
		r[k], err = StringToFloat64(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToInts converts string slice to int slice.
func StringsToInts(s []string, emptyAsZero ...bool) ([]int, error) {
	var err error
	r := make([]int, len(s))
	for k, v := range s {
		r[k], err = StringToInt(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToInt8s converts string slice to int8 slice.
func StringsToInt8s(s []string, emptyAsZero ...bool) ([]int8, error) {
	var err error
	r := make([]int8, len(s))
	for k, v := range s {
		r[k], err = StringToInt8(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToInt16s converts string slice to int16 slice.
func StringsToInt16s(s []string, emptyAsZero ...bool) ([]int16, error) {
	var err error
	r := make([]int16, len(s))
	for k, v := range s {
		r[k], err = StringToInt16(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToInt32s converts string slice to int32 slice.
func StringsToInt32s(s []string, emptyAsZero ...bool) ([]int32, error) {
	var err error
	r := make([]int32, len(s))
	for k, v := range s {
		r[k], err = StringToInt32(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToInt64s converts string slice to int64 slice.
func StringsToInt64s(s []string, emptyAsZero ...bool) ([]int64, error) {
	var err error
	r := make([]int64, len(s))
	for k, v := range s {
		r[k], err = StringToInt64(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToUints converts string slice to uint slice.
func StringsToUints(s []string, emptyAsZero ...bool) ([]uint, error) {
	var err error
	r := make([]uint, len(s))
	for k, v := range s {
		r[k], err = StringToUint(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToUint8s converts string slice to uint8 slice.
func StringsToUint8s(s []string, emptyAsZero ...bool) ([]uint8, error) {
	var err error
	r := make([]uint8, len(s))
	for k, v := range s {
		r[k], err = StringToUint8(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToUint16s converts string slice to uint16 slice.
func StringsToUint16s(s []string, emptyAsZero ...bool) ([]uint16, error) {
	var err error
	r := make([]uint16, len(s))
	for k, v := range s {
		r[k], err = StringToUint16(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToUint32s converts string slice to uint32 slice.
func StringsToUint32s(s []string, emptyAsZero ...bool) ([]uint32, error) {
	var err error
	r := make([]uint32, len(s))
	for k, v := range s {
		r[k], err = StringToUint32(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsToUint64s converts string slice to uint64 slice.
func StringsToUint64s(s []string, emptyAsZero ...bool) ([]uint64, error) {
	var err error
	r := make([]uint64, len(s))
	for k, v := range s {
		r[k], err = StringToUint64(v, emptyAsZero...)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// StringsCopyWithin copies part of an slice to another location in the current slice.
// @target
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func StringsCopyWithin(s []string, target, start int, end ...int) {
	target = fixIndex(len(s), target, true)
	if target == len(s) {
		return
	}
	sub := StringsSlice(s, start, end...)
	for i, v := range sub {
		s[target+i] = v
	}
}

// StringsEvery tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice will return true for any condition!
func StringsEvery(s []string, fn func(s []string, k int, v string) bool) bool {
	for k, v := range s {
		if !fn(s, k, v) {
			return false
		}
	}
	return true
}

// StringsFill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//  Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
// @start
//  Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
// @end
//  Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//  If negative, end will be counted from the end.
//  If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func StringsFill(s []string, value string, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(s), start, end...)
	if !ok {
		return
	}
	for i := fixedStart; i < fixedEnd; i++ {
		s[i] = value
	}
}

// StringsFilter creates a new slice with all elements that pass the test implemented by the provided function.
func StringsFilter(s []string, fn func(s []string, k int, v string) bool) []string {
	ret := make([]string, 0)
	for k, v := range s {
		if fn(s, k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// StringsFind returns the key-value of the first element in the provided slice that satisfies the provided testing function.
// NOTE:
//  If not found, k = -1
func StringsFind(s []string, fn func(s []string, k int, v string) bool) (k int, v string) {
	for k, v := range s {
		if fn(s, k, v) {
			return k, v
		}
	}
	return -1, ""
}

// StringsIncludes determines whether an slice includes a certain value among its entries.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func StringsIncludes(s []string, valueToFind string, fromIndex ...int) bool {
	return StringsIndexOf(s, valueToFind, fromIndex...) > -1
}

// StringsIndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func StringsIndexOf(s []string, searchElement string, fromIndex ...int) int {
	idx := getFromIndex(len(s), fromIndex...)
	for k, v := range s[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// StringsJoin concatenates the elements of s to create a single string. The separator string
// sep is placed between elements in the resulting string.
func StringsJoin(s []string, sep string) string {
	return strings.Join(s, sep)
}

// StringsLastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//  The index to start the search at. Defaults to 0.
func StringsLastIndexOf(s []string, searchElement string, fromIndex ...int) int {
	idx := getFromIndex(len(s), fromIndex...)
	for i := len(s) - 1; i >= idx; i-- {
		if searchElement == s[i] {
			return i
		}
	}
	return -1
}

// StringsMap creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func StringsMap(s []string, fn func(s []string, k int, v string) string) []string {
	ret := make([]string, len(s))
	for k, v := range s {
		ret[k] = fn(s, k, v)
	}
	return ret
}

// StringsPop removes the last element from an slice and returns that element.
// This method changes the length of the slice.
func StringsPop(s *[]string) (string, bool) {
	a := *s
	if len(a) == 0 {
		return "", false
	}
	lastIndex := len(a) - 1
	last := a[lastIndex]
	a = a[:lastIndex]
	*s = a[:len(a):len(a)]
	return last, true
}

// StringsPush adds one or more elements to the end of an slice and returns the new length of the slice.
func StringsPush(s *[]string, element ...string) int {
	*s = append(*s, element...)
	return len(*s)
}

// StringsPushDistinct adds one or more new elements that do not exist in the current slice at the end.
func StringsPushDistinct(s []string, element ...string) []string {
L:
	for _, v := range element {
		for _, vv := range s {
			if vv == v {
				continue L
			}
		}
		s = append(s, v)
	}
	return s
}

// StringsReduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func StringsReduce(s []string, fn func(s []string, k int, v, accumulator string) string, initialValue ...string) string {
	if len(s) == 0 {
		return ""
	}
	start := 0
	acc := s[start]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		start += 1
	}
	for i := start; i < len(s); i++ {
		acc = fn(s, i, s[i], acc)
	}
	return acc
}

// StringsReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//  The accumulator accumulates callback's return values.
//  It is the accumulated value previously returned in the last invocation of the callback—or initialValue,
//  if it was supplied (see below).
// @initialValue
//  A value to use as the first argument to the first call of the callback.
//  If no initialValue is supplied, the first element in the slice will be used and skipped.
func StringsReduceRight(s []string, fn func(s []string, k int, v, accumulator string) string, initialValue ...string) string {
	if len(s) == 0 {
		return ""
	}
	end := len(s) - 1
	acc := s[end]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		end -= 1
	}
	for i := end; i >= 0; i-- {
		acc = fn(s, i, s[i], acc)
	}
	return acc
}

// StringsReverse reverses an slice in place.
func StringsReverse(s []string) {
	first := 0
	last := len(s) - 1
	for first < last {
		s[first], s[last] = s[last], s[first]
		first++
		last--
	}
}

// StringsShift removes the first element from an slice and returns that removed element.
// This method changes the length of the slice.
func StringsShift(s *[]string) (string, bool) {
	a := *s
	if len(a) == 0 {
		return "", false
	}
	first := a[0]
	a = a[1:]
	*s = a[:len(a):len(a)]
	return first, true
}

// StringsSlice returns a copy of a portion of an slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func StringsSlice(s []string, begin int, end ...int) []string {
	fixedStart, fixedEnd, ok := fixRange(len(s), begin, end...)
	if !ok {
		return []string{}
	}
	return StringsCopy(s[fixedStart:fixedEnd])
}

// StringsSome tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//  Calling this method on an empty slice returns false for any condition!
func StringsSome(s []string, fn func(s []string, k int, v string) bool) bool {
	for k, v := range s {
		if fn(s, k, v) {
			return true
		}
	}
	return false
}

// StringsSplice changes the contents of an slice by removing or replacing
// existing elements and/or adding new elements in place.
func StringsSplice(s *[]string, start, deleteCount int, items ...string) {
	a := *s
	if deleteCount < 0 {
		deleteCount = 0
	}
	start, end, _ := fixRange(len(a), start, start+1+deleteCount)
	deleteCount = end - start - 1
	for i := 0; i < len(items); i++ {
		if deleteCount > 0 {
			// replace
			a[start] = items[i]
			deleteCount--
			start++
		} else {
			// insert
			lastSlice := StringsCopy(a[start:])
			items = items[i:]
			a = append(a[:start], items...)
			a = append(a[:start+len(items)], lastSlice...)
			*s = a[:len(a):len(a)]
			return
		}
	}
	if deleteCount > 0 {
		a = append(a[:start], a[start+1+deleteCount:]...)
	}
	*s = a[:len(a):len(a)]
}

// StringsUnshift adds one or more elements to the beginning of an slice and returns the new length of the slice.
func StringsUnshift(s *[]string, element ...string) int {
	*s = append(element, *s...)
	return len(*s)
}

// StringsUnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func StringsUnshiftDistinct(s *[]string, element ...string) int {
	a := *s
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[string]bool, len(element))
	r := make([]string, 0, len(a)+len(element))
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
	*s = r[:len(r):len(r)]
	return len(r)
}

// StringsRemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func StringsRemoveFirst(p *[]string, elements ...string) int {
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

// StringsRemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func StringsRemoveEvery(p *[]string, elements ...string) int {
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

// StringsConcat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func StringsConcat(s ...[]string) []string {
	var totalLen int
	for _, v := range s {
		totalLen += len(v)
	}
	ret := make([]string, totalLen)
	dst := ret
	for _, v := range s {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// StringsIntersect calculates intersection of two or more slices,
// and returns the count of each element.
func StringsIntersect(s ...[]string) (intersectCount map[string]int) {
	if len(s) == 0 {
		return nil
	}
	for _, v := range s {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[string]int, len(s))
	for k, v := range s {
		counts[k] = stringsDistinct(v, nil)
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

// StringsDistinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func StringsDistinct(s *[]string, changeSlice bool) (distinctCount map[string]int) {
	if !changeSlice {
		return stringsDistinct(*s, nil)
	}
	a := (*s)[:0]
	distinctCount = stringsDistinct(*s, &a)
	n := len(distinctCount)
	*s = a[:n:n]
	return distinctCount
}

func stringsDistinct(src []string, dst *[]string) map[string]int {
	m := make(map[string]int, len(src))
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

// StringSetUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func StringSetUnion(set1, set2 []string, others ...[]string) []string {
	m := make(map[string]struct{}, len(set1)+len(set2))
	r := make([]string, 0, len(m))
	for _, set := range append([][]string{set1, set2}, others...) {
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

// StringSetIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func StringSetIntersect(set1, set2 []string, others ...[]string) []string {
	sets := append([][]string{set2}, others...)
	setsCount := make([]map[string]int, len(sets))
	for k, v := range sets {
		setsCount[k] = stringsDistinct(v, nil)
	}
	m := make(map[string]struct{}, len(set1))
	r := make([]string, 0, len(m))
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

// StringSetDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func StringSetDifference(set1, set2 []string, others ...[]string) []string {
	m := make(map[string]struct{}, len(set1))
	r := make([]string, 0, len(set1))
	sets := append([][]string{set2}, others...)
	for _, v := range sets {
		inter := StringSetIntersect(set1, v)
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

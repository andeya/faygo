package ameda

// FormatUintByDict convert num into corresponding string according to dict.
func FormatUintByDict(dict []byte, num uint64) string {
	var base = uint64(len(dict))
	if base == 0 {
		return ""
	}
	var str []byte
	for {
		tmp := make([]byte, len(str)+1)
		tmp[0] = dict[num%base]
		copy(tmp[1:], str)
		str = tmp
		num = num / base
		if num == 0 {
			break
		}
	}
	return string(str)
}

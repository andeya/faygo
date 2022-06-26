package ameda

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

// ParseUintByDict convert numStr into corresponding uint64 according to dict.
func ParseUintByDict(dict []byte, numStr string) (uint64, error) {
	if len(dict) == 0 {
		return 0, errors.New("dict is empty")
	}
	base := float64(len(dict))
	len := len(numStr)
	var number float64
	for i := 0; i < len; i++ {
		char := numStr[i : i+1]
		pos := bytes.IndexAny(dict, char)
		if pos == -1 {
			return 0, fmt.Errorf("found a char not included in the dict: %q", char)
		}
		number = math.Pow(base, float64(len-i-1))*float64(pos) + number
	}
	return uint64(number), nil
}

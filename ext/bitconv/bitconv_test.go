package bitconv

import (
	"math"
	"testing"
)

func TestFormat(t *testing.T) {
	t.Log(New().Format(math.MaxInt64))
}

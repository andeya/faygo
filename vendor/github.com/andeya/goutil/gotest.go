package goutil

import (
	"os"
	"strings"
)

// IsGoTest returns whether the current process is a test.
func IsGoTest() bool {
	return isGoTest
}

var isGoTest bool

func init() {
	isGoTest = checkGoTestEnv()
}

const IS_GO_TEST = "IS_GO_TEST"

func checkGoTestEnv() bool {
	for _, arg := range os.Args[1:] {
		for _, s := range []string{
			"-test.timeout=",
			"-test.timeout",
			"-test.run=",
			"-test.run",
			"-test.bench=",
			"-test.bench",
			"-test.v=",
			"-test.v",
		} {
			if strings.HasPrefix(arg, s) || arg == s {
				return true
			}
		}
	}
	env := os.Getenv(IS_GO_TEST)
	return env == "1" || env == "true"
	// return strings.HasSuffix(os.Args[0], ".test")
}

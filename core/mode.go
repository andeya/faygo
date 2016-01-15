// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"io"
	"os"

	"github.com/henrylee2cn/thinkgo/core/binding"
	"github.com/mattn/go-colorable"
)

const ENV_THINKGO_MODE = "THINKGO_MODE"

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
	TestMode    string = "test"
)
const (
	debugCode   = iota
	releaseCode = iota
	testCode    = iota
)

var DefaultWriter io.Writer = colorable.NewColorableStdout()
var ginMode int = debugCode
var modeName string = DebugMode

// func init() {
// 	mode := os.Getenv(ENV_THINKGO_MODE)
// 	if len(mode) == 0 {
// 		SetMode(DebugMode)
// 	} else {
// 		SetMode(mode)
// 	}
// }

func SetMode(mode string) {
	if len(mode) == 0 {
		mode = os.Getenv(ENV_THINKGO_MODE)
	}
	switch mode {
	case DebugMode:
		ginMode = debugCode
	case ReleaseMode:
		ginMode = releaseCode
	case TestMode:
		ginMode = testCode
	default:
		panic("thinkgo mode unknown: " + mode)
	}
	modeName = mode
}

func DisableBindValidation() {
	binding.Validator = nil
}

func Mode() string {
	return modeName
}

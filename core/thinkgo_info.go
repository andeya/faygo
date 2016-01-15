// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

// 框架信息
const (
	NAME    = "ThinkGo"
	VERSION = "v0.10"
	AUTHOR  = "Henrylee2cn"
)

func ThinkGoInfo() string {
	var s string
	s += "[THINKGO-INFO] @NAME:		" + NAME + "\n"
	s += "[THINKGO-INFO] @VERSION:	" + VERSION + "\n"
	s += "[THINKGO-INFO] @AUTHOR:		" + AUTHOR + "\n"
	return s
}

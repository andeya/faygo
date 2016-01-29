// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package application

func init() {
	ModulePrepare(&Module{
		Name:        "Admin",
		Description: "后台管理模块",
	}).Router(&IndexController{})
}

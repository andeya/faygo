// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"net/http"
	"path"
)

type (
	// 控制器接口
	Controller interface {
		AutoInit(method string, ctx *Context, name string, callfunc string, module *Module) Controller
	}
	// 基础控制器
	BaseController struct {
		// 请求方法
		method string
		// 请求上下文
		*Context
		// 控制器名称
		name string
		// 本次调用的函数名称
		callfunc string
		// 所属模块
		Module *Module
		// 模板路径
		tplpath string
		// html模板变量
		Data H
	}
)

// 自动初始化
func (this *BaseController) AutoInit(method string, ctx *Context, name string, callfunc string, module *Module) Controller {
	this.method = method
	this.Context = ctx
	this.name = name
	this.callfunc = callfunc
	this.Module = module

	if method == "GET" {
		this.tplpath = path.Join(APP_PACKAGE,
			module.RouterGroup.BasePath(),
			VIEW_PACKAGE,
			module.Themes.Cur,
			name,
			callfunc) + ThinkGo.TplSuffex

		this.Data = H{
			// 定义模板中"__PUBLIC__"静态文件前缀
			"__PUBLIC__": path.Join(PUBLIC_PREFIX, module.RouterGroup.BasePath(), module.Themes.Cur),
		}
	}

	return this
}

// HTML renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (this *BaseController) HTML(code ...int) {
	if len(code) == 0 {
		code = append(code, http.StatusOK)
	}

	this.Context.HTML(code[0], this.tplpath, this.Data)
}

// IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
// It also sets the Content-Type as "application/json".
// WARNING: we recommend to use this only for development propuses since printing pretty JSON is
// more CPU and bandwidth consuming. Use Context.JSON() instead.
func (this *BaseController) IndentedJSON(obj interface{}, code ...int) {
	if len(code) == 0 {
		code = append(code, http.StatusOK)
	}
	this.Context.IndentedJSON(code[0], obj)
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (this *BaseController) JSON(obj interface{}, code ...int) {
	if len(code) == 0 {
		code = append(code, http.StatusOK)
	}
	this.Context.JSON(code[0], obj)
}

// XML serializes the given struct as XML into the response body.
// It also sets the Content-Type as "application/xml".
func (this *BaseController) XML(obj interface{}, code ...int) {
	if len(code) == 0 {
		code = append(code, http.StatusOK)
	}
	this.Context.XML(code[0], obj)
}

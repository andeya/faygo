// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"fmt"
	"path"
	"regexp"
)

type Think struct {
	*Echo
	// 模块列表
	Modules *Modules
	// 模板引擎
	*Template
	// 配置信息
	Config
	// 框架信息
	Author  string
	Version string
}

// 重要配置，涉及项目架构，请勿修改
const (
	// 模块应用目录名
	APP_PACKAGE = "application"
	// 视图文件目录名
	VIEW_PACKAGE = "view"
	// 公共目录
	COMMON_PACKAGE = "common"
	// 资源文件目录名
	PUBLIC_PACKAGE = "__public__"
	// 访问的资源文件url前缀
	PUBLIC_PREFIX = "/public"
	// 上传根目录名
	UPLOADS_PACKAGE = "uploads"
)

// 全局运行实例
var (
	ThinkGo = newThinkGo()
)

func newThinkGo() *Think {
	t := &Think{
		// 业务数据
		Echo:    New().Group("/").Echo(),
		Modules: newModules(),
		Config:  getConfig(),
		// 框架信息
		Author:  AUTHOR,
		Version: VERSION,
	}

	log := t.Logger()
	log.SetPrefix("TG")

	t.Echo.SetDebug(t.Config.Debug)
	t.Echo.SetLogLevel(t.Config.LogLevel)
	t.Template = NewRender()
	t.Template.Delims(t.Config.TplLeft, t.Config.TplRight)
	t.Template.SetBasepath(APP_PACKAGE)
	t.Template.SetSuffix(t.Config.TplSuffix)
	t.Template.SetDebug(t.Config.Debug)
	t.Template.Prepare()
	t.Echo.SetRenderer(t.Template)
	t.servedir()
	if t.Echo.Debug() {
		for k, v := range t.Map() {
			t.logger.Notice("	%-25s --> %-25s", k, v)
		}
	}
	// t.Echo.SetBinder(b)
	// t.Echo.SetHTTPErrorHandler(HTTPErrorHandler)
	// t.Echo.SetLogOutput(w io.Writer)
	// t.Echo.SetHTTPErrorHandler(h HTTPErrorHandler)
	return t
}

func (this *Think) servedir() {
	this.Echo.Favicon("deploy/favicon/favicon.ico")
	this.Echo.ServeDir("/uploads", UPLOADS_PACKAGE)
	this.Echo.ServeDir("/common", APP_PACKAGE+"/"+COMMON_PACKAGE+"/"+VIEW_PACKAGE+"/"+PUBLIC_PACKAGE)

	var re = regexp.MustCompile(APP_PACKAGE + "(/[^/]+)/" + VIEW_PACKAGE + "(/[^/]+)/" + PUBLIC_PACKAGE)
	for _, p := range WalkRelDirs(APP_PACKAGE, "/"+PUBLIC_PACKAGE) {
		a := re.FindStringSubmatch(p)
		if len(a) == 3 {
			// public/[模块]/[主题]/
			if a[1] == "/home" {
				a[1] = "/"
			}
			this.Echo.ServeDir(path.Join(PUBLIC_PREFIX, a[1], a[2]), p)
		}
	}
}

func (this *Think) Run() {
	this.Echo.Run(fmt.Sprintf("%s:%d", this.Config.HttpAddr, this.Config.HttpPort))
}

func (this *Think) Use(m ...Middleware) *Think {
	this.Echo.Use(m...)
	return this
}

// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"fmt"
	"regexp"
)

// 重要配置，涉及项目架构，请勿修改
const (
	// 模块应用目录名
	APP_PACKAGE = "application"
	// 视图文件目录名
	VIEW_PACKAGE = "view"
	// 资源文件目录名
	PUBLIC_PACKAGE = "__public__"
	// 访问的资源文件url前缀
	PUBLIC_PREFIX = "/public"
	// 上传根目录名
	UPLOADS_PACKAGE = "uploads"
)

// 全局运行实例
var ThinkGo = func() *Think {
	// 打印框架信息
	fmt.Println(ThinkGoInfo())

	config := readConfig()

	// 设置运行模式，当传入参数为空时，使用环境变量 "THINKGO_MODE"
	SetMode(config.RunMode)

	t := &Think{
		// 框架信息
		Author:  AUTHOR,
		Version: VERSION,
		// 业务数据
		Modules: newModules(),
		Addons:  newAddons(),
		Engine:  newEngine(),
		Config:  config,
	}

	return t
}()

type Think struct {
	*Engine
	// 模块列表
	Modules *Modules
	// 插件列表
	Addons *Addons
	// 配置信息
	Config
	// 框架信息
	Author  string
	Version string
}

// 供main函数调用
func (this *Think) Use(middleware ...HandlerFunc) *Think {
	this.Engine.Use(middleware...)
	return this
}

func (this *Think) TemplateFuncs(funcMap map[string]interface{}) *Think {
	this.Engine.HTMLTemplateFuncs(funcMap)
	return this
}

func (this *Think) Run() (err error) {
	// 配置必须的静态文件服务器
	this.Engine.StaticFile("/favicon.ico", "deploy/favicon/favicon.ico")
	this.Engine.StaticFile("/uploads", UPLOADS_PACKAGE)

	var re = regexp.MustCompile(APP_PACKAGE + "(/[^/]+)/" + VIEW_PACKAGE + "(/[^/]+)/" + PUBLIC_PACKAGE)
	for _, p := range WalkRelDirs(APP_PACKAGE, "/"+PUBLIC_PACKAGE) {
		a := re.FindStringSubmatch(p)
		if len(a) < 3 {
			continue
		}
		this.Engine.Static(PUBLIC_PREFIX+a[1]+a[2], p)
	}

	// 设置模板定界符
	this.Engine.HTMLTemplateDelims(this.Config.TplLeft, this.Config.TplRight)
	// 遍历模板文件
	tplFiles := WalkRelFiles(APP_PACKAGE, this.Config.TplSuffex)
	// 添加模板
	this.Engine.LoadHTMLFiles(tplFiles...)

	// 开启服务
	err = this.Engine.Run(fmt.Sprintf("%s:%d", this.Config.HttpAddr, this.Config.HttpPort))
	return
}

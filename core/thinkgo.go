// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"fmt"
	"regexp"

	"github.com/henrylee2cn/thinkgo/conf"
	"github.com/henrylee2cn/thinkgo/core/template"
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

// 供内部调用
var App = func() *Think {
	// 打印框架信息
	fmt.Println(ThinkGoInfo())

	// 设置运行模式，当传入参数为空时，使用环境变量 "THINKGO_MODE"
	SetMode(conf.THINKGO_MODE)

	return &Think{
		// 应用的基本信息
		Name:          conf.NAME,
		Description:   conf.DESCRIPTION,
		Version:       conf.VERSION,
		Company:       conf.COMPANY,
		Developer:     conf.DEVELOPER,
		CopyrightYear: conf.COPYRIGHT_YEAR,
		Contact:       conf.CONTACT,

		// 业务数据
		Modules:   newModules(),
		Addons:    newAddons(),
		Engine:    newEngine(),
		tplDelims: [2]string{"{{{", "}}}"},
		tplFuncs:  map[string]interface{}{},
		tplSuffex: ".html",
	}
}()

type Think struct {
	*Engine
	// 模块列表
	Modules *Modules
	// 插件列表
	Addons *Addons
	// 模板定界符
	tplDelims [2]string
	// 模板函数
	tplFuncs template.FuncMap
	// 模板后缀名
	tplSuffex string

	// 应用名称
	Name string
	// 应用描述
	Description string
	// 应用版本号
	Version string
	// 公司名称
	Company string
	// 开发者
	Developer string
	// 版权年限
	CopyrightYear string
	// 联系方式
	Contact string
}

// 供main函数调用
func ThinkGo() *Think {
	return App
}

// 供main函数调用
func ThinkGoDefault() *Think {
	App.Engine.Use(Recovery(), Logger())
	return App
}

func (this *Think) TemplateDelims(left, right string) *Think {
	this.tplDelims = [2]string{left, right}
	return this
}

func (this *Think) TemplateFuncs(funcMap map[string]interface{}) *Think {
	this.Engine.HTMLTemplateFuncs(funcMap)
	return this
}

func (this *Think) TemplateSuffex(ext string) {
	this.tplSuffex = ext
}

func (this *Think) Run(addr ...string) (err error) {
	// 配置必须的静态文件服务器
	this.Engine.StaticFile("/favicon.ico", "common/deploy/favicon/favicon.ico")
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
	this.Engine.HTMLTemplateDelims(this.tplDelims[0], this.tplDelims[1])

	// 遍历模板文件
	tplFiles := WalkRelFiles(APP_PACKAGE, this.tplSuffex)

	// 添加模板
	App.Engine.LoadHTMLFiles(tplFiles...)

	// 开启服务
	err = this.Engine.Run(addr...)
	return
}

func (this *Think) String() string {
	var s string
	s += "[APP] Name:		" + this.Name + "\n"
	s += "[APP] Description:	" + this.Description + "\n"
	s += "[APP] Version:		" + this.Version + "\n"
	s += "[APP] Company:		" + this.Company + "\n"
	s += "[APP] Developer:	" + this.Developer + "\n"
	s += "[APP] CopyrightYear:	" + this.CopyrightYear + "\n"
	s += "[APP] Contact:		" + this.Contact + "\n"
	return s
}

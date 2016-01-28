// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/henrylee2cn/thinkgo/core/template"
)

type (
	// 控制器接口
	Controller interface {
		AutoInit(ctx *Context) Controller
	}
	// 基础控制器
	BaseController struct {
		// 请求上下文
		*Context
		// 子模板
		sectionTpl map[string]string
	}
)

// 自动初始化
func (this *BaseController) AutoInit(ctx *Context) Controller {
	this.Context = ctx
	return this
}

func (this *BaseController) Render(code ...int) error {
	if len(code) == 0 {
		code = append(code, http.StatusOK)
	}
	return this.Context.Render(code[0], this.Context.Path(), this.Context.GetAll())
}

func (this *BaseController) RenderLayout(layoutName string, code ...int) error {
	if len(code) == 0 {
		code = append(code, http.StatusOK)
	}
	render := this.Echo().Render
	for k, v := range this.sectionTpl {
		sectionBytes := bytes.NewBufferString("")
		render(sectionBytes, v, this.Context.GetAll())
		sectionContent, _ := ioutil.ReadAll(sectionBytes)
		this.Set(k, template.HTML(sectionContent))
	}
	return this.Context.Render(code[0], layoutName, this.Context.GetAll())
}

func (this *BaseController) SetSection(position string, sectionName ...string) {
	if len(sectionName) == 0 {
		sectionName = append(sectionName, this.Path())
	}
	if this.sectionTpl == nil {
		this.sectionTpl = make(map[string]string)
	}
	this.sectionTpl[position] = sectionName[0]
}

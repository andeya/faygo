// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// 应用模块
type Module struct {
	id          string
	Name        string
	Description string
	*Themes
	*Group
}

var (
	Modules = map[string]*Module{}
	re      = regexp.MustCompile("^[/]?([a-zA-Z0-9_]+)([\\./\\?])?")
)

// 创建模块
// 自动设置default主题
// 文件名作为id，且文件名应与模块目录名、包名保存一致
func NewModule(description string, mw ...Middleware) *Module {
	m := &Module{
		Description: description,
		Themes:      &Themes{},
	}

	// 设置默认主题
	m.Themes.Set(&Theme{
		Name:        "default",
		Description: "default",
		Src:         map[string]string{},
	})

	// 设置id
	_, file, _, _ := runtime.Caller(1)
	m.id = strings.TrimSuffix(filepath.Base(file), ".go")

	// 设置Name
	m.Name = strings.Title(m.id)

	// 生成url前缀
	prefix := "/" + m.id

	// 创建分组并修改请求路径c.path "/[模块]/[控制器]/[操作]"为"/[模块]/[主题]/[控制器]/[操作]"
	m.Group = ThinkGo.Echo.Group(
		prefix,
		func(c *Context) error {
			// 补全主题字段
			p := strings.Split(c.Path(), "/:")[0]
			p = path.Join(prefix, m.Themes.Cur().Name, strings.TrimPrefix(p, prefix))
			c.SetPath(p)
			// 静态文件前缀
			c.Set("__PUBLIC__", path.Join("/public", prefix, m.Themes.Cur().Name))
			return nil
		},
		Recover(),
		Logger(),
	)
	m.Group.Use(mw...)

	// 模块登记
	Modules[m.id] = m

	return m
}

// 获取Id
func (this *Module) GetId() string {
	return this.id
}

// 设置Id
func (this *Module) SetId(id string) *Module {
	this.id = id
	return this
}

// 获取Name
func (this *Module) GetName() string {
	return this.Name
}

// 获取Description
func (this *Module) GetDescription() string {
	return this.Description
}

// 设置主题，并默认设置传入的第1个主题为当前主题
func (this *Module) SetThemes(themes ...*Theme) *Module {
	this.Themes.Set(themes...)
	return this
}

// 设置当前主题
func (this *Module) UseTheme(name string) *Module {
	this.Themes.Use(name)
	return this
}

// 定义中间件
func (this *Module) Use(m ...Middleware) *Module {
	this.Group.Use(m...)
	return this
}

// 注册智能路由
func (this *Module) Router(c Controller, m ...Middleware) *Module {
	t := reflect.TypeOf(c)
	e := t.Elem()
	cname := SnakeString(strings.TrimSuffix(e.Name(), "Controller"))
	group := this.Group.Group(cname, m...)
	for i := t.NumMethod() - 1; i >= 0; i-- {
		fname := t.Method(i).Name
		idx := strings.LastIndex(fname, "_")
		if idx == -1 {
			continue
		}
		pattern := SnakeString(fname[:idx])
		method := strings.ToUpper(fname[idx+1:])
		switch method {
		case "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "TRACE", "SOCKET":
			group.Match([]string{method}, pattern, func(ctx *Context) error {
				var v = reflect.New(e)
				v.Interface().(Controller).AutoInit(ctx)
				rets := v.MethodByName(fname).Call([]reflect.Value{})
				if len(rets) > 0 {
					if err, ok := rets[0].Interface().(error); ok {
						return err
					}
				}
				return nil
			})
		}
	}
	return this
}

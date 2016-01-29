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
	"sync"
)

type (
	// 应用模块
	Module struct {
		Name        string
		Class       string
		Description string
		id          string
		status      int
		sync.Mutex
		*Themes
		*Group
	}
	// 登记模块列表
	Modules struct {
		// 快速调用列表
		Map map[string]*Module
		// 有序列表 [分组][Id]*Module
		Slice [][]*Module
	}
)

var (
	re = regexp.MustCompile("^[/]?([a-zA-Z0-9_]+)([\\./\\?])?")
)

func newModules() *Modules {
	return &Modules{
		Map:   map[string]*Module{},
		Slice: [][]*Module{},
	}
}

// 初始化模块，文件名作为id，且文件名应与模块目录名、包名保存一致
func ModulePrepare(m *Module) *Module {
	_, file, _, _ := runtime.Caller(1)
	m.id = strings.TrimSuffix(filepath.Base(file), ".go")
	prefix := "/" + m.id
	// 创建分组并修改请求路径c.path "/[模块]/[控制器]/[操作]"为"/[模块]/[主题]/[控制器]/[操作]"
	m.Group = ThinkGo.echo.Group(prefix, func(c *Context) error {
		// 补全主题字段
		p := strings.Split(c.Path(), "/:")[0]
		p = path.Join(prefix, m.Themes.Cur, strings.TrimPrefix(p, prefix))
		c.SetPath(p)
		// 静态文件前缀
		c.Set("__PUBLIC__", path.Join("/public", prefix, m.Themes.Cur))
		return nil
	})
	m.Group.Use(Recover(), Logger())

	// 登记并排序
	insertModule(m)
	return m
}

// 获取Id
func (this *Module) GetId() string {
	return this.id
}

// 设置主题，自动设置传入的第1个主题为当前主题
func (this *Module) SetThemes(themes ...*Theme) *Module {
	if len(themes) == 0 {
		return this
	}
	this.Themes = &Themes{
		Cur:  themes[0].Name,
		List: make(map[string]*Theme),
	}
	this.Themes.AddThemes(themes...)
	return this
}

// 设置当前主题
func (this *Module) UseTheme(name string) *Module {
	this.Themes.UseTheme(name)
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

// 顺序插入插件
func insertModule(m *Module) {
	// 添加至插件索引列表
	ThinkGo.Modules.Map[m.id] = m

	// 添加至插件有序列表
	var (
		add   bool
		class []string
	)

	for _, ms := range ThinkGo.Modules.Slice {
		c := ms[0].Class
		class = append(class, c)
		if c != m.Class {
			continue
		}
		for k, v := range ms {
			if v.Name > m.Name {
				x := append([]*Module{m}, ms[k:]...)
				ms = append(ms[:k], x...)
				break
			}
		}
		add = true
		break
	}
	if add {
		return
	}

	if len(class) == 0 {
		ThinkGo.Modules.Slice = append(ThinkGo.Modules.Slice, []*Module{m})
		return
	}

	for k, v := range class {
		if v > m.Class {
			x := append([][]*Module{{m}}, ThinkGo.Modules.Slice[k:]...)
			ThinkGo.Modules.Slice = append(ThinkGo.Modules.Slice[:k], x...)
			break
		}
	}
}

func GetMoudleSlice() [][]*Module {
	return ThinkGo.Modules.Slice
}

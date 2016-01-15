// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"log"
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
		*Themes
		id     string
		status int
		*RouterGroup
		sync.Mutex
	}
	// 登记模块列表
	Modules struct {
		// 快速调用列表
		Map map[string]*Module
		// 有序列表 [分组][Id]*Module
		Slice [][]*Module
	}
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

	// 初始化
	m.RouterGroup = App.Engine.Group(m.id)

	// 登记并排序
	insertModule(m)
	return m
}

// 获取Id
func (this *Module) GetId() string {
	return this.id
}

// 获取url
func (this *Module) GetUrl(pattern string) string {
	return path.Join(this.id, pattern)
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
func (this *Module) Use(middleware ...HandlerFunc) *Module {
	this.RouterGroup.Use(middleware...)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) GET(pattern string, controllerOrhandler ...interface{}) *Module {
	this.router("GET", pattern, controllerOrhandler)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) POST(pattern string, controllerOrhandler ...interface{}) *Module {
	this.router("POST", pattern, controllerOrhandler)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) HEAD(pattern string, controllerOrhandler ...interface{}) *Module {
	this.router("HEAD", pattern, controllerOrhandler)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) PUT(pattern string, controllerOrhandler ...interface{}) *Module {
	this.router("PUT", pattern, controllerOrhandler)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) DELETE(pattern string, controllerOrhandler ...interface{}) *Module {
	this.router("DELETE", pattern, controllerOrhandler)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) PATCH(pattern string, controllerOrhandler ...interface{}) *Module {
	this.router("PATCH", pattern, controllerOrhandler)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) OPTIONS(pattern string, controllerOrhandler ...interface{}) *Module {
	this.router("OPTIONS", pattern, controllerOrhandler)
	return this
}

// 公用路由注册方法
// 注：路由规则"/ReadMail/:id" 将被自动转为 "/read_mail/:id"
// 注：路由规则"read" 将被自动转为 "/read"
var re = regexp.MustCompile("^([/]?[a-zA-Z0-9_]+)([\\./\\?])?")

func (this *Module) router(method, pattern string, controllerOrhandler []interface{}) {
	this.Mutex.Lock()
	defer func() {
		recover()
		this.Mutex.Unlock()
	}()
	if pattern[0] != '/' {
		pattern = "/" + SnakeString(pattern)
	} else {
		pattern = "/" + SnakeString(pattern[1:])
	}
	pattern = strings.Replace(pattern, "/?", "?", -1)
	pattern = strings.Trim(pattern, "?")
	pattern = strings.TrimSuffix(pattern, "/")
	a := re.FindStringSubmatch(pattern)
	if len(a) < 3 {
		log.Panicln(`[ERROR]  配置路由规则: 匹配规则 "` + pattern + `" 不正确`)
	}
	var (
		hfs             = make([]HandlerFunc, len(controllerOrhandler))
		countController int
		cName           string
		callfunc        = CamelString(strings.TrimPrefix(a[1], "/"))
	)
	for i, v := range controllerOrhandler {
		c, ok := v.(Controller)
		if ok {
			cName, callfunc, hfs[i] = this.newHandler(method, callfunc, c)
			countController++
			continue
		}
		h, ok := v.(HandlerFunc)
		if ok {
			hfs[i] = h
			continue
		}
		log.Panicln(`[ERROR] 配置路由规则: "` + this.RouterGroup.BasePath() + method + `" 指定了类型错误的操作`)
	}
	if countController != 1 {
		log.Panicln(`[ERROR] 配置路由规则: "` + this.RouterGroup.BasePath() + method + `" 须且仅须设置1个控制器`)
	}

	callMethod := reflect.ValueOf(this.RouterGroup).MethodByName(method)
	hfsv := reflect.ValueOf(hfs)
	if callfunc == "index" && a[2] != "/" {
		// 允许省略index
		p := path.Join(cName, pattern[len(a[1]):])
		p = strings.Replace(p, "/?", "?", -1)
		callMethod.CallSlice([]reflect.Value{reflect.ValueOf(p), hfsv})
		if cName == "index" {
			callMethod.CallSlice([]reflect.Value{reflect.ValueOf(pattern[len(a[1]):]), hfsv})
		}
	}
	callMethod.CallSlice([]reflect.Value{reflect.ValueOf(path.Join(cName, pattern)), hfsv})
}

// 返回闭包操作
func (this *Module) newHandler(method, callfunc string, c Controller) (_name, _callfunc string, hf HandlerFunc) {
	var (
		t   = reflect.TypeOf(c)
		has bool
	)
	for i := t.NumMethod() - 1; i >= 0; i-- {
		mName := t.Method(i).Name
		if strings.EqualFold(mName, callfunc) {
			callfunc = mName
			has = true
			break
		}
	}
	t = t.Elem()
	if !has {
		log.Panicln(`[ERROR]  配置路由规则: 指定方法名 "` + callfunc + `" 在控制器 [` + t.Name() + `] 不存在(忽略大小写)`)
	}

	_name = SnakeString(strings.TrimSuffix(t.Name(), "Controller"))
	_callfunc = SnakeString(callfunc)
	hf = func(ctx *Context) {
		var newCiValue = reflect.New(t)
		// 预处理
		newCiValue.Interface().(Controller).AutoInit(method, ctx, _name, _callfunc, this)
		// 开始执行
		newCiValue.MethodByName(callfunc).Call([]reflect.Value{})
	}
	return
}

// 顺序插入插件
func insertModule(m *Module) {
	// 添加至插件索引列表
	App.Modules.Map[m.id] = m

	// 添加至插件有序列表
	var (
		add   bool
		class []string
	)

	for _, ms := range App.Modules.Slice {
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
		App.Modules.Slice = append(App.Modules.Slice, []*Module{m})
		return
	}

	for k, v := range class {
		if v > m.Class {
			x := append([][]*Module{{m}}, App.Modules.Slice[k:]...)
			App.Modules.Slice = append(App.Modules.Slice[:k], x...)
			break
		}
	}
}

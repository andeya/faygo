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
		Id          string
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
	m.Id = strings.TrimSuffix(filepath.Base(file), ".go")
	prefix := "/" + m.Id
	if m.Id == "home" {
		prefix = ""
	}
	// 创建分组并修改请求路径c.path "/[模块]/[控制器]/[操作]"为"/[模块]/[主题]/[控制器]/[操作]"
	m.Group = ThinkGo.Echo.Group(prefix, func(c *Context) error {
		p := strings.Split(c.Path(), "/:")[0]
		if p == "/" || p == "" {
			c.SetPath("/index/index")
		} else if strings.HasSuffix(p, "/index") {
			l := 3
			if c.Echo().Prefix() == "/" {
				l = 2
			}
			num := l - strings.Count(p, "/")
			if num > 0 {
				c.SetPath(p + strings.Repeat("/index", num))
			}
		}
		p = path.Join(prefix, m.Themes.Cur, strings.TrimPrefix(p, prefix))
		// 插入主题字段
		c.SetPath(p)
		// 静态文件前缀
		c.Set("__PUBLIC__", path.Join(PUBLIC_PREFIX, prefix, m.Themes.Cur))
		return nil
	})

	// 登记并排序
	insertModule(m)
	return m
}

// 获取Id
func (this *Module) GetId() string {
	return this.Id
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

// @pattern: 方法名[/参数规则]
func (this *Module) CONNECT(pattern string, c Controller) *Module {
	this.router(CONNECT, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) DELETE(pattern string, c Controller) *Module {
	this.router(DELETE, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) GET(pattern string, c Controller) *Module {
	this.router(GET, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) HEAD(pattern string, c Controller) *Module {
	this.router(HEAD, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) OPTIONS(pattern string, c Controller) *Module {
	this.router(OPTIONS, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) PATCH(pattern string, c Controller) *Module {
	this.router(PATCH, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) POST(pattern string, c Controller) *Module {
	this.router(POST, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) PUT(pattern string, c Controller) *Module {
	this.router(PUT, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) TRACE(pattern string, c Controller) *Module {
	this.router(TRACE, pattern, c)
	return this
}

// @pattern: 方法名[/参数规则]
func (this *Module) SOCKET(pattern string, c Controller) *Module {
	this.router(SOCKET, pattern, c)
	return this
}

// 公用路由注册方法
// 注：路由规则"/ReadMail/:id" 将被自动转为 "/read_mail/:id"
// 注：路由规则"read" 将被自动转为 "/read"
// 注：当为"home"模块时，同时在根目录注册路由
func (this *Module) router(method, pattern string, c Controller) {
	this.Mutex.Lock()
	defer func() {
		recover()
		this.Mutex.Unlock()
	}()
	pattern, a := dealPattern(pattern)

	cname, fname, h := this.newHandler(CamelString(a[1]), c)

	echo := this.Group.Echo()
	prefix := echo.Prefix()
	cname = "/" + cname
	if method != SOCKET {
		add := echo.Match
		if prefix == "/home" {
			// 当为"home"模块时，添加注册根路由
			add = ThinkGo.Echo.Match
		}

		// 允许省略index
		if fname == "index" && a[2] != "/" {
			p := path.Join(cname, pattern[len(a[1])+1:])
			// p = strings.Replace(p, "/?", "?", -1)
			add([]string{method}, p, h)
			if cname == "/index" {
				add([]string{method}, "/"+pattern[len(a[1])+1:], h)
			}
		}
		add([]string{method}, path.Join(cname, pattern), h)

	} else {
		add := echo.WebSocket
		if prefix == "/home" {
			// 当为"home"模块时，添加注册根路由
			add = ThinkGo.Echo.WebSocket
		}

		// 允许省略index
		if fname == "index" && a[2] != "/" {
			p := path.Join(cname, pattern[len(a[1]):])
			// p = strings.Replace(p, "/?", "?", -1)
			add(p, h)
			if cname == "/index" {
				add("/"+pattern[len(a[1]):], h)
			}
		}
		add(path.Join(cname, pattern), h)
	}
}

// 返回闭包操作
func (this *Module) newHandler(fname string, c Controller) (_cname, _fname string, h HandlerFunc) {
	t := reflect.TypeOf(c)
	has := false
	for i := t.NumMethod() - 1; i >= 0; i-- {
		mname := t.Method(i).Name
		if strings.EqualFold(mname, fname) {
			fname = mname
			has = true
			break
		}
	}
	t = t.Elem()
	if !has {
		ThinkGo.Logger().Fatal(`[ERROR]  配置路由规则: 指定方法名 "` + fname + `" 在控制器 [` + t.Name() + `] 不存在(忽略大小写)`)
	}

	_cname = SnakeString(strings.TrimSuffix(t.Name(), "Controller"))
	_fname = SnakeString(fname)
	h = func(ctx *Context) error {
		var newv = reflect.New(t)
		newv.Interface().(Controller).AutoInit(ctx, this)
		vs := newv.MethodByName(fname).Call([]reflect.Value{})
		if len(vs) > 0 {
			if err, ok := vs[0].Interface().(error); ok {
				return err
			}
		}
		return nil
	}
	return
}

func dealPattern(s string) (string, []string) {
	s = strings.Trim(s, "/")
	s = strings.Trim(s, "?")
	// s = strings.Replace(s, "/?", "?", -1)
	s = strings.Split(s, "?")[0]
	a := re.FindStringSubmatch(s)
	s = "/" + SnakeString(s)
	if len(a) < 3 {
		ThinkGo.Logger().Fatal(`[ERROR]  配置路由规则: 匹配规则 "` + s + `" 不正确`)
	}
	return s, a
}

// 顺序插入插件
func insertModule(m *Module) {
	// 添加至插件索引列表
	ThinkGo.Modules.Map[m.Id] = m

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

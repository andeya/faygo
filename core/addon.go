// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"path/filepath"
	"runtime"
	"strings"
)

type (
	// 应用插件
	Addon struct {
		*Module
	}
	// 登记插件列表
	Addons struct {
		// 快速调用列表
		Map map[string]*Addon
		// 有序列表 [分组][Id]*Addon
		Slice [][]*Addon
	}
)

func newAddons() *Addons {
	return &Addons{
		Map:   map[string]*Addon{},
		Slice: [][]*Addon{},
	}
}

// 初始化插件，文件名作为id，且文件名应与插件目录名、包名保存一致
func AddonPrepare(m *Module) *Addon {
	_, file, _, _ := runtime.Caller(1)
	m.id = strings.TrimSuffix(filepath.Base(file), ".go")

	// 初始化
	m.RouterGroup = App.Engine.Group(m.id)

	a := &Addon{
		Module: m,
	}

	// 登记并排序
	insertAddon(a)
	return a
}

// 顺序插入插件
func insertAddon(a *Addon) {
	// 添加至插件索引列表
	App.Addons.Map[a.id] = a

	// 添加至插件有序列表
	var (
		add   bool
		class []string
	)

	for _, as := range App.Addons.Slice {
		c := as[0].Class
		class = append(class, c)
		if c != a.Class {
			continue
		}
		for k, v := range as {
			if v.Name > a.Name {
				x := append([]*Addon{a}, as[k:]...)
				as = append(as[:k], x...)
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
		App.Addons.Slice = append(App.Addons.Slice, []*Addon{a})
		return
	}

	for k, v := range class {
		if v > a.Class {
			x := append([][]*Addon{{a}}, App.Addons.Slice[k:]...)
			App.Addons.Slice = append(App.Addons.Slice[:k], x...)
			break
		}
	}
}

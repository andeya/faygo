// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"sort"
)

type (
	Themes struct {
		cur string
		Map map[string]*Theme
	}
	Theme struct {
		Name        string
		Description string
		Src         map[string]string // 预览图片地址等
	}
)

func NewThemes() *Themes {
	return &Themes{
		Map: make(map[string]*Theme),
	}
}

func (this *Themes) List() []*Theme {
	i := len(this.Map)
	a := make([]string, i)
	for k, _ := range this.Map {
		i--
		a[i] = k
	}
	sort.Strings(a)
	ts := make([]*Theme, len(a))
	for i, k := range a {
		ts[i] = this.Map[k]
	}
	return ts
}

func (this *Themes) Cur() *Theme {
	return this.Map[this.cur]
}

func (this *Themes) Use(name string) {
	this.cur = name
}

func (this *Themes) Add(themes ...*Theme) {
	for _, theme := range themes {
		this.Map[theme.Name] = theme
	}
}

// 设置主题，并默认设置传入的第1个主题为当前主题
func (this *Themes) Set(themes ...*Theme) *Themes {
	if len(themes) == 0 {
		return this
	}

	this.Map = make(map[string]*Theme)
	for _, theme := range themes {
		this.Map[theme.Name] = theme
	}

	this.cur = themes[0].Name

	return this
}

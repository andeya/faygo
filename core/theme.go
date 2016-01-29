// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

type (
	Themes struct {
		Cur  string
		List map[string]*Theme
	}
	Theme struct {
		Name        string
		Description string
		Src         map[string]string // 预览图片地址等
	}
)

func NewThemes() *Themes {
	return &Themes{
		List: make(map[string]*Theme),
	}
}

func (this *Themes) CurTheme() *Theme {
	return this.List[this.Cur]
}

func (this *Themes) UseTheme(name string) {
	this.Cur = name
}

func (this *Themes) AddThemes(themes ...*Theme) {
	for _, theme := range themes {
		this.List[theme.Name] = theme
	}
}

// 设置主题，并默认设置传入的第1个主题为当前主题
func (this *Themes) SetThemes(themes ...*Theme) *Themes {
	if len(themes) == 0 {
		return this
	}

	this.List = make(map[string]*Theme)
	for _, theme := range themes {
		this.List[theme.Name] = theme
	}

	this.Cur = themes[0].Name

	return this
}

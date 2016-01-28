// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io"

	"github.com/henrylee2cn/thinkgo/core/template"
)

type Template struct {
	debug    bool
	suffix   string
	basepath string
	delims   [2]string
	*template.Template
	// [[模块]/[主题]/[控制器]/[操作]]:path
	pathmap map[string]string
	// [[模块]/[主题]/[控制器]/[操作]]:html
	text map[string]string
}

func NewRender() *Template {
	return &Template{
		Template: template.New("thinkgo").Funcs(template.FuncMap{}),
		pathmap:  make(map[string]string),
	}
}

func (t *Template) Delims(left, right string) {
	t.delims = [2]string{left, right}
}

func (t *Template) SetSuffix(suffix string) {
	t.suffix = suffix
}

func (t *Template) SetBasepath(basepath string) {
	t.basepath = basepath
}

func (t *Template) SetDebug(debug bool) {
	t.debug = debug
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	f := t.pathmap[name]
	if f == "" {
		return fmt.Errorf("索引模板不存在: %s", name)
	}
	if t.debug {
		return template.Must(template.Must(t.Template.Clone()).ParseFiles(f)).ExecuteTemplate(w, f, data)
	}
	return t.Template.ExecuteTemplate(w, f, data)
}

func (t *Template) Map() map[string]string {
	return t.pathmap
}

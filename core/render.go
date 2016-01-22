// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io"
	"regexp"

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

func (t *Template) Prepare() {
	var re = regexp.MustCompile(t.basepath + "(/[^/]+)/" + VIEW_PACKAGE + "(/[^/]+)(/[^/]+)(/[^/]+)" + t.suffix)
	var re2 = regexp.MustCompile(t.basepath + "/" + COMMON_PACKAGE + "/" + VIEW_PACKAGE + "(/[^/]+)" + t.suffix)
	var paths []string
	for _, f := range WalkRelFiles(t.basepath, t.suffix) {
		a := re.FindStringSubmatch(f)
		if len(a) < 5 {
			b := re2.FindStringSubmatch(f)
			if len(b) == 2 {
				t.pathmap["/"+COMMON_PACKAGE+b[1]] = f
				paths = append(paths, f)
			}
			continue
		}
		if a[1] == "/home" {
			a[1] = ""
		}
		r := a[1] + a[2] + a[3] + a[4]
		t.pathmap[r] = f
		paths = append(paths, f)
	}
	if !t.debug {
		t.Template.ParseFiles(paths...)
	}

	t.Template.Delims(t.delims[0], t.delims[1])
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

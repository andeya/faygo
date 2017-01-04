// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package thinkgo

import (
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/henrylee2cn/thinkgo/ini"
	"github.com/henrylee2cn/thinkgo/utils"
)

// JoinStatic adds the static directory prefix to the file name.
func JoinStatic(shortFilename string) string {
	return path.Join(StaticDir(), shortFilename)
}

// SyncINI quickly create your own configuration files.
// Struct tags reference `https://github.com/go-ini/ini`
func SyncINI(structPointer interface{}, callback func() error, filename ...string) error {
	t := reflect.TypeOf(structPointer)
	if t.Kind() != reflect.Ptr {
		return errors.New("SyncINI's param must be struct pointer type.")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("SyncINI's param must be struct pointer type.")
	}

	var fname string
	if len(filename) > 0 {
		fname = filename[0]
	} else {
		fname = strings.TrimSuffix(t.Name(), "Config")
		fname = strings.TrimSuffix(fname, "INI")
		fname = utils.SnakeString(fname) + ".ini"
		fname = filepath.Join(CONFIG_DIR, fname)
	}
	var cfg *ini.File
	var err error
	var exist bool
	cfg, err = ini.Load(fname)
	if err != nil {
		os.MkdirAll(filepath.Dir(fname), 0777)
		cfg, err = ini.LooseLoad(fname)
		if err != nil {
			return err
		}
	} else {
		exist = true
	}

	err = cfg.MapTo(structPointer)
	if err != nil {
		return err
	}

	if callback != nil {
		if err = callback(); err != nil {
			return err
		}
	}

	if !exist {
		err = cfg.ReflectFrom(structPointer)
		if err != nil {
			return err
		}
		return cfg.SaveTo(fname)
	}
	return nil
}

/**
 * WrapDoc add a document notes to handler
 */
type docWrap struct {
	Handler
	doc Doc
}

var _ APIDoc = new(docWrap)

func (w *docWrap) Doc() Doc {
	return w.doc
}

// WrapDoc adds a note to the handler
func WrapDoc(handler Handler, note string, ret interface{}, params ...ParamInfo) Handler {
	return &docWrap{
		Handler: handler,
		doc: Doc{
			Note:   note,
			Return: ret,
			Params: params,
		},
	}
}

/**
 * define common middlewares.
 */

// NewIPFilter creates middleware that intercepts the specified IP prefix.
func NewIPFilter(whitelist []string, realIP bool) HandlerFunc {
	var noAccess bool
	var match []string
	var prefix []string

	if len(whitelist) == 0 {
		noAccess = true
	} else {
		for _, s := range whitelist {
			if strings.HasSuffix(s, "*") {
				prefix = append(prefix, s[:len(s)-1])
			} else {
				match = append(match, s)
			}
		}
	}

	return func(ctx *Context) error {
		if noAccess {
			ctx.Error(http.StatusForbidden, "no access")
			return nil
		}

		var ip string
		if realIP {
			ip = ctx.RealIP()
		} else {
			ip = ctx.IP()
		}
		for _, ipMatch := range match {
			if ipMatch == ip {
				ctx.Next()
				return nil
			}
		}
		for _, ipPrefix := range prefix {
			if strings.HasPrefix(ip, ipPrefix) {
				ctx.Next()
				return nil
			}
		}
		ctx.Error(http.StatusForbidden, "not allow to access: "+ip)
		return nil
	}
}

// CrossOrigin creates Cross-Domain middleware
var CrossOrigin = HandlerFunc(func(ctx *Context) error {
	ctx.SetHeader(HeaderAccessControlAllowOrigin, ctx.HeaderParam(HeaderOrigin))
	// ctx.SetHeader(HeaderAccessControlAllowOrigin, "*")
	ctx.SetHeader(HeaderAccessControlAllowCredentials, "true")
	return nil
})

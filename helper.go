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
func SyncINI(structPointer interface{}, callback func() error, filename ...string) (err error) {
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

	os.MkdirAll(filepath.Dir(fname), 0777)

	cfg, err := ini.LooseLoad(fname)
	if err != nil {
		return err
	}

	err = cfg.MapTo(structPointer)
	if err != nil {
		return err
	}

	if callback != nil {
		if err := callback(); err != nil {
			return err
		}
	}

	err = cfg.ReflectFrom(structPointer)
	if err != nil {
		return err
	}

	return cfg.SaveTo(fname)
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

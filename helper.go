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
	"strings"
)

// The static directory prefix is automatically added for the file name
func JoinStatic(shortFilename string) string {
	return path.Join(Global.staticDir, shortFilename)
}

/*
 * RenderFS
 */

// New a file system with auto-rendering.
func RenderFS(frame *Framework, pattern, root string, tplVar Map) *MuxAPI {
	return NamedRenderFS(frame, "renderserver", pattern, root, tplVar)
}

// New a file system with auto-rendering.
func NamedRenderFS(frame *Framework, name, pattern, root string, tplVar Map) *MuxAPI {
	return frame.NamedStaticFS(name, pattern, &renderFS{
		dir:    root,
		tplVar: tplVar,
		render: Global.Render(),
	}, true, false)
}

type renderFS struct {
	dir    string
	tplVar Map
	render *Render
}

func (fs *renderFS) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("RenderFS: invalid character in file path")
	}
	dir := fs.dir
	if dir == "" {
		dir = "."
	}
	fname := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	b, fileInfo, err := fs.render.renderForFS(fname, fs.tplVar)
	if err != nil {
		if strings.Contains(err.Error(), "not find") {
			return nil, os.ErrNotExist
		}
		// Error("RenderFS:", err)
		return NewFile(b, fileInfo), err
	}
	return NewFile(b, fileInfo), nil
}

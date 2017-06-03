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

package faygo

import (
	"errors"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	fayerrors "github.com/henrylee2cn/faygo/errors"
	"github.com/henrylee2cn/faygo/ini"
	"github.com/henrylee2cn/faygo/utils"
)

// JoinStatic adds the static directory prefix to the file name.
func JoinStatic(shortFilename string) string {
	return path.Join(StaticDir(), shortFilename)
}

// SyncINI quickly create your own configuration files.
// Struct tags reference `https://github.com/go-ini/ini`
func SyncINI(structPointer interface{}, callback func(existed bool, saveOnce func() error) error, filename ...string) error {
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
	var existed bool
	cfg, err = ini.Load(fname)
	if err != nil {
		os.MkdirAll(filepath.Dir(fname), 0777)
		cfg, err = ini.LooseLoad(fname)
		if err != nil {
			return err
		}
	} else {
		existed = true
	}

	err = cfg.MapTo(structPointer)
	if err != nil {
		return err
	}

	var once sync.Once
	var saveOnce = func() error {
		var err error
		once.Do(func() {
			err = cfg.ReflectFrom(structPointer)
			if err != nil {
				return
			}
			err = cfg.SaveTo(fname)
			if err != nil {
				return
			}
		})
		return err
	}

	if callback != nil {
		if err = callback(existed, saveOnce); err != nil {
			return err
		}
	}

	if !existed {
		return saveOnce()
	}
	return nil
}

// RemoveUseless when there's not frame instance, remove files: config, log, static and upload .
func RemoveUseless() {
	if len(AllFrames()) > 0 {
		return
	}
	var files []string
	filepath.Walk(CONFIG_DIR, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, retpath)
		return err
	})
	confile := filepath.Join(CONFIG_DIR, GLOBAL_CONFIG_FILE)
	if len(files) == 1 || len(files) == 2 && files[1] == confile {
		os.Remove(confile)
		os.Remove(CONFIG_DIR)
		os.Remove(LogDir())
		os.Remove(StaticDir())
		os.Remove(UploadDir())
	}
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

// WrapDoc adds a note to the handler func.
func WrapDoc(fn HandlerFunc, note string, ret interface{}, params ...ParamInfo) Handler {
	return &docWrap{
		Handler: fn,
		doc: Doc{
			Note:   note,
			Return: ret,
			Params: params,
		},
	}
}

/**
 * common utils
 */

// ContentTypeByExtension gets the content type from ext string.
// MIME type is given in mime package.
// It returns `application/octet-stream` incase MIME type is not
// found.
func ContentTypeByExtension(ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		return ctype
	}
	return MIMEOctetStream
}

// SelfPath gets compiled executable file absolute path.
//  func SelfPath() string
var SelfPath = utils.SelfPath

// SelfDir gets compiled executable file directory
//  func SelfDir() string
var SelfDir = utils.SelfDir

// RelPath gets relative path.
//  func RelPath() string
var RelPath = utils.RelPath

// SelfChdir switch the working path to my own path.
//  func SelfChdir()
var SelfChdir = utils.SelfChdir

// FileExists reports whether the named file or directory exists.
//  func FileExists(name string) bool
var FileExists = utils.FileExists

// SearchFile Search a file in paths.
// this is often used in search config file in /etc ~/
//  func SearchFile(filename string, paths ...string) (fullpath string, err error)
var SearchFile = utils.SearchFile

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
//  func GrepFile(patten string, filename string) (lines []string, err error)
var GrepFile = utils.GrepFile

// WalkDirs traverses the directory, return to the relative path.
// You can specify the suffix.
//  func WalkDirs(targpath string, suffixes ...string) (dirlist []string)
var WalkDirs = utils.WalkDirs

// SnakeString converts the accepted string to a snake string (XxYy to xx_yy)
//  func SnakeString(s string) string
var SnakeString = utils.SnakeString

// CamelString converts the accepted string to a camel string (xx_yy to XxYy)
//  func CamelString(s string) string
var CamelString = utils.CamelString

// ObjectName gets the type name of the object
//  func ObjectName(i interface{}) string
var ObjectName = utils.ObjectName

// CleanPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
// 1. Replace multiple slashes with a single slash.
// 2. Eliminate each . path name element (the current directory).
// 3. Eliminate each inner .. path name element (the parent directory) along with the non-.. element that precedes it.
// 4. Eliminate .. elements that begin a rooted path: that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned.
//  func CleanPath(p string) string
var CleanPath = utils.CleanPath

// RandomString returns a URL-safe, base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
//  func RandomString(n int) string
var RandomString = utils.RandomString

// Errors merge multiple errors.
//  func Errors(errs []error) error
var Errors = fayerrors.Errors

// Byte2String convert []byte type to string type.
//  func Bytes2String(b []byte) string
var Bytes2String = utils.Bytes2String

// String2Bytes convert *string type to []byte type.
// NOTE: panic if modify the member value of the []byte.
//  func String2Bytes(s *string) []byte
var String2Bytes = utils.String2Bytes

/**
 * define internal middlewares.
 */

// newIPFilter creates middleware that intercepts the specified IP prefix.
func newIPFilter(whitelist []string, realIP bool) HandlerFunc {
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

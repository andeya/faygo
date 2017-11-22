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
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/henrylee2cn/ini"

	"github.com/henrylee2cn/goutil"
	"github.com/henrylee2cn/goutil/errors"
)

// JoinStatic adds the static directory prefix to the file name.
func JoinStatic(shortFilename string) string {
	return path.Join(StaticDir(), shortFilename)
}

// SyncINI quickly create your own configuration files.
// Struct tags reference `https://github.com/go-ini/ini`
func SyncINI(structPtr interface{}, f func(onecUpdateFunc func() error) error, filename ...string) error {
	t := reflect.TypeOf(structPtr)
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
		fname = goutil.SnakeString(fname) + ".ini"
		fname = filepath.Join(CONFIG_DIR, fname)
	}
	return ini.SyncINI(structPtr, f, fname)
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

// WritePid write pid to the specified file.
func WritePid(pidFilename string) error {
	abs, err := filepath.Abs(pidFilename)
	if err != nil {
		return err
	}
	dir := filepath.Dir(abs)
	os.MkdirAll(dir, 0777)
	pid := os.Getpid()
	return ioutil.WriteFile(abs, []byte(fmt.Sprintf("%d\n", pid)), 0666)
}

// CleanToURL is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned
func CleanToURL(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 2 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r++

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 2

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return BytesToString(buf[:w])
}

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}

// SelfPath gets compiled executable file absolute path.
//  func SelfPath() string
var SelfPath = goutil.SelfPath

// SelfDir gets compiled executable file directory
//  func SelfDir() string
var SelfDir = goutil.SelfDir

// RelPath gets relative path.
//  func RelPath() string
var RelPath = goutil.RelPath

// SelfChdir switch the working path to my own path.
//  func SelfChdir()
var SelfChdir = goutil.SelfChdir

// FileExists reports whether the named file or directory exists.
//  func FileExists(name string) bool
var FileExists = goutil.FileExists

// SearchFile Search a file in paths.
// this is often used in search config file in /etc ~/
//  func SearchFile(filename string, paths ...string) (fullpath string, err error)
var SearchFile = goutil.SearchFile

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
//  func GrepFile(patten string, filename string) (lines []string, err error)
var GrepFile = goutil.GrepFile

// WalkDirs traverses the directory, return to the relative path.
// You can specify the suffix.
//  func WalkDirs(targpath string, suffixes ...string) (dirlist []string)
var WalkDirs = goutil.WalkDirs

// SnakeString converts the accepted string to a snake string (XxYy to xx_yy)
//  func SnakeString(s string) string
var SnakeString = goutil.SnakeString

// CamelString converts the accepted string to a camel string (xx_yy to XxYy)
//  func CamelString(s string) string
var CamelString = goutil.CamelString

// ObjectName gets the type name of the object
//  func ObjectName(i interface{}) string
var ObjectName = goutil.ObjectName

// RandomString returns a URL-safe, base64 encoded securely generated
// random string. It will panic if the system's secure random number generator
// fails to function correctly.
// The length n must be an integer multiple of 4, otherwise the last character will be padded with `=`.
//  func RandomString(n int) string
var RandomString = goutil.URLRandomString

// BytesToString convert []byte type to string type.
//  func BytesToString(b []byte) string
var BytesToString = goutil.BytesToString

// StringToBytes convert string type to []byte type.
// NOTE: panic if modify the member value of the []byte.
//  func StringToBytes(s string) []byte
var StringToBytes = goutil.StringToBytes

// JsQueryEscape escapes the string in javascript standard so it can be safely placed
// inside a URL query.
//  func JsQueryEscape(s string) string
var JsQueryEscape = goutil.JsQueryEscape

// JsQueryUnescape does the inverse transformation of JsQueryEscape, converting
// %AB into the byte 0xAB and '+' into ' ' (space). It returns an error if
// any % is not followed by two hexadecimal digits.
//  func JsQueryUnescape(s string) (string, error)
var JsQueryUnescape = goutil.JsQueryUnescape

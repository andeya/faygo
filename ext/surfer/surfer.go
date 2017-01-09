// Copyright 2015 henrylee2cn Author. All Rights Reserved.
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
//
// Package surfer is a high level concurrency http client.
// It has `surf` and` phantom` download engines, highly simulated browser behavior, the function of analog login and so on.
// Features:
// - Both surf and phantomjs engines are supported
// - Support random User-Agent
// - Support cache cookie
// - Support http/https
//
// Usage:
// package main

// import (
//     "github.com/henrylee2cn/surfer"
//     "io/ioutil"
//     "log"
// )
//
// func main() {
//     // Use surf engine
//     resp, err := surfer.Download(&surfer.Request{
//         Url: "http://github.com/henrylee2cn/surfer",
//     })
//     if err != nil {
//         log.Fatal(err)
//     }
//     b, err := ioutil.ReadAll(resp.Body)
//     log.Println(string(b), err)
//
//     // Use phantomjs engine
//     resp, err = surfer.Download(&surfer.Request{
//         Url:          "http://github.com/henrylee2cn",
//         DownloaderID: 1,
//     })
//     if err != nil {
//         log.Fatal(err)
//     }
//     b, err = ioutil.ReadAll(resp.Body)
//     log.Println(string(b), err)
//     resp.Body.Close()
//     surfer.DestroyJsFiles()
// }
package surfer

import (
	"net/http"
	"sync"
	// "os"
	// "path"
	// "path/filepath"
)

var (
	surf         Surfer
	phantom      Surfer
	once_surf    sync.Once
	once_phantom sync.Once
	tempJsDir    = "./tmp"
	// phantomjsFile = filepath.Clean(path.Join(os.Getenv("GOPATH"), `/src/github.com/henrylee2cn/surfer/phantomjs/phantomjs`))
	phantomjsFile = `./phantomjs`
)

// Download 实现surfer下载器接口
func Download(req *Request) (resp *http.Response, err error) {
	switch req.DownloaderID {
	case SurfID:
		once_surf.Do(func() { surf = New() })
		resp, err = surf.Download(req)
	case PhomtomJsID:
		once_phantom.Do(func() { phantom = NewPhantom(phantomjsFile, tempJsDir) })
		resp, err = phantom.Download(req)
	}
	return
}

// DestroyJsFiles 销毁Phantomjs的js临时文件
func DestroyJsFiles() {
	if pt, ok := phantom.(*Phantom); ok {
		pt.DestroyJsFiles()
	}
}

// Surfer represents an core of HTTP web browser for crawler.
type Surfer interface {
	// GET @param url string, header http.Header, cookies []*http.Cookie
	// HEAD @param url string, header http.Header, cookies []*http.Cookie
	// POST PostForm @param url, referer string, values url.Values, header http.Header, cookies []*http.Cookie
	// POST-M PostMultipart @param url, referer string, values url.Values, header http.Header, cookies []*http.Cookie
	Download(*Request) (resp *http.Response, err error)
}

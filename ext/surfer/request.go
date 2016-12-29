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

package surfer

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	SurfID             = 0               // Surf下载器标识符
	PhomtomJsID        = 1               // PhomtomJs下载器标识符
	DefaultMethod      = "GET"           // 默认请求方法
	DefaultDialTimeout = 2 * time.Minute // 默认请求服务器超时
	DefaultConnTimeout = 2 * time.Minute // 默认下载超时
	DefaultTryTimes    = 3               // 默认最大下载次数
	DefaultRetryPause  = 2 * time.Second // 默认重新下载前停顿时长
)

// Request contains the necessary prerequisite information.
type Request struct {
	// url (必须填写)
	Url string
	url *url.URL
	// GET POST POST-M HEAD (默认为GET)
	Method string
	// http header
	Header http.Header
	// 是否使用cookies，在Spider的EnableCookie设置
	EnableCookie bool
	// POST values
	Values url.Values
	// POST files
	Files []PostFile
	body  io.Reader
	// dial tcp: i/o timeout
	DialTimeout time.Duration
	// WSARecv tcp: i/o timeout
	ConnTimeout time.Duration
	// the max times of download
	TryTimes int
	// how long pause when retry
	RetryPause time.Duration
	// max redirect times
	// when RedirectTimes equal 0, redirect times is ∞
	// when RedirectTimes less than 0, redirect times is 0
	RedirectTimes int
	// the download ProxyHost
	Proxy string
	proxy *url.URL
	// 指定下载器ID
	// 0为Surf高并发下载器，各种控制功能齐全
	// 1为PhantomJS下载器，特点破防力强，速度慢，低并发
	DownloaderID int
	client       *http.Client
}

// PostFile post form file
type PostFile struct {
	Fieldname string
	Filename  string
	Bytes     []byte
}

func (self *Request) prepare() error {
	var err error
	self.url, err = UrlEncode(self.Url)
	if err != nil {
		return err
	}
	self.Url = self.url.String()
	if self.Proxy != "" {
		if self.proxy, err = url.Parse(self.Proxy); err != nil {
			return err
		}
	}
	if self.DialTimeout < 0 {
		self.DialTimeout = 0
	} else if self.DialTimeout == 0 {
		self.DialTimeout = DefaultDialTimeout
	}

	if self.ConnTimeout < 0 {
		self.ConnTimeout = 0
	} else if self.ConnTimeout == 0 {
		self.ConnTimeout = DefaultConnTimeout
	}

	if self.TryTimes == 0 {
		self.TryTimes = DefaultTryTimes
	}

	if self.RetryPause <= 0 {
		self.RetryPause = DefaultRetryPause
	}

	if self.DownloaderID != PhomtomJsID {
		self.DownloaderID = SurfID
	}

	if self.Header == nil {
		self.Header = make(http.Header)
	}
	var commonUserAgentIndex int
	if !self.EnableCookie {
		commonUserAgentIndex = rand.Intn(len(UserAgents["common"]))
		self.Header.Set("User-Agent", UserAgents["common"][commonUserAgentIndex])
	} else if len(self.Header["User-Agent"]) == 0 {
		self.Header.Set("User-Agent", UserAgents["common"][commonUserAgentIndex])
	}

	self.Method = strings.ToUpper(self.Method)
	switch self.Method {
	case "POST", "PUT", "PATCH", "DELETE":
		if len(self.Files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for _, postfile := range self.Files {
					fileWriter, err := bodyWriter.CreateFormFile(postfile.Fieldname, postfile.Filename)
					if err != nil {
						log.Println("[E] Surfer: Multipart:", err)
					}
					_, err = fileWriter.Write(postfile.Bytes)
					if err != nil {
						log.Println("[E] Surfer: Multipart:", err)
					}
				}
				for k, v := range self.Values {
					for _, vv := range v {
						bodyWriter.WriteField(k, vv)
					}
				}
				bodyWriter.Close()
				pw.Close()
			}()
			self.Header.Set("Content-Type", bodyWriter.FormDataContentType())
			self.body = ioutil.NopCloser(pr)
		} else {
			self.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			self.body = strings.NewReader(self.Values.Encode())
		}

	default:
		if len(self.Method) == 0 {
			self.Method = DefaultMethod
		}
	}

	return nil
}

// 回写Request内容
func (self *Request) writeback(resp *http.Response) *http.Response {
	if resp == nil {
		resp = new(http.Response)
		resp.Request = new(http.Request)
	} else if resp.Request == nil {
		resp.Request = new(http.Request)
	}

	resp.Request.Method = self.Method
	resp.Request.Header = self.Header
	resp.Request.Host = self.url.Host

	return resp
}

// checkRedirect is used as the value to http.Client.CheckRedirect
// when redirectTimes equal 0, redirect times is ∞
// when redirectTimes less than 0, not allow redirects
func (self *Request) checkRedirect(req *http.Request, via []*http.Request) error {
	if self.RedirectTimes == 0 {
		return nil
	}
	if len(via) >= self.RedirectTimes {
		if self.RedirectTimes < 0 {
			return fmt.Errorf("not allow redirects.")
		}
		return fmt.Errorf("stopped after %v redirects.", self.RedirectTimes)
	}
	return nil
}

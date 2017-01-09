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
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// constant
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
	// GET POST HEAD (默认为GET)
	Method string
	// http header
	Header http.Header
	// 是否使用cookies，在Spider的EnableCookie设置
	EnableCookie bool
	// request body interface
	Body body
	body io.Reader
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

func (r *Request) prepare() error {
	var err error
	r.url, err = UrlEncode(r.Url)
	if err != nil {
		return err
	}
	r.Url = r.url.String()
	if r.Proxy != "" {
		if r.proxy, err = url.Parse(r.Proxy); err != nil {
			return err
		}
	}
	if r.DialTimeout < 0 {
		r.DialTimeout = 0
	} else if r.DialTimeout == 0 {
		r.DialTimeout = DefaultDialTimeout
	}

	if r.ConnTimeout < 0 {
		r.ConnTimeout = 0
	} else if r.ConnTimeout == 0 {
		r.ConnTimeout = DefaultConnTimeout
	}

	if r.TryTimes == 0 {
		r.TryTimes = DefaultTryTimes
	}

	if r.RetryPause <= 0 {
		r.RetryPause = DefaultRetryPause
	}

	if r.DownloaderID != PhomtomJsID {
		r.DownloaderID = SurfID
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	var commonUserAgentIndex int
	if !r.EnableCookie {
		commonUserAgentIndex = rand.Intn(len(UserAgents["common"]))
		r.Header.Set("User-Agent", UserAgents["common"][commonUserAgentIndex])
	} else if len(r.Header["User-Agent"]) == 0 {
		r.Header.Set("User-Agent", UserAgents["common"][commonUserAgentIndex])
	}
	if len(r.Method) == 0 {
		r.Method = DefaultMethod
	} else {
		r.Method = strings.ToUpper(r.Method)
	}
	r.body = nil
	if r.Body != nil {
		return r.Body.SetBody(r)
	}
	return nil
}

// ReadBody returns body bytes
func (r *Request) ReadBody() (b []byte, err error) {
	if r.url == nil {
		r.prepare()
	}
	if r.body != nil {
		b, err = ioutil.ReadAll(r.body)
	}
	return b, err
}

// 回写Request内容
func (r *Request) writeback(resp *http.Response) *http.Response {
	if resp == nil {
		resp = new(http.Response)
		resp.Request = new(http.Request)
	} else if resp.Request == nil {
		resp.Request = new(http.Request)
	}

	resp.Request.Method = r.Method
	resp.Request.Header = r.Header
	resp.Request.Host = r.url.Host

	// reset url
	r.url = nil

	return resp
}

// checkRedirect is used as the value to http.Client.CheckRedirect
// when redirectTimes equal 0, redirect times is ∞
// when redirectTimes less than 0, not allow redirects
func (r *Request) checkRedirect(req *http.Request, via []*http.Request) error {
	if r.RedirectTimes == 0 {
		return nil
	}
	if len(via) >= r.RedirectTimes {
		if r.RedirectTimes < 0 {
			return fmt.Errorf("not allow redirects.")
		}
		return fmt.Errorf("stopped after %v redirects.", r.RedirectTimes)
	}
	return nil
}

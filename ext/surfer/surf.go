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
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

// Surf is the default Download implementation.
type Surf struct {
	cookieJar *cookiejar.Jar
}

// New 创建一个Surf下载器
func New() Surfer {
	s := new(Surf)
	s.cookieJar, _ = cookiejar.New(nil)
	return s
}

// Download 实现surfer下载器接口
func (surf *Surf) Download(param *Request) (*http.Response, error) {
	err := param.prepare()
	if err != nil {
		return nil, err
	}
	param.client = surf.buildClient(param)
	resp, err := surf.httpRequest(param)

	if err == nil {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			var gzipReader *gzip.Reader
			gzipReader, err = gzip.NewReader(resp.Body)
			if err == nil {
				resp.Body = gzipReader
			}

		case "deflate":
			resp.Body = flate.NewReader(resp.Body)

		case "zlib":
			var readCloser io.ReadCloser
			readCloser, err = zlib.NewReader(resp.Body)
			if err == nil {
				resp.Body = readCloser
			}
		}
	}

	return param.writeback(resp), err
}

// buildClient creates, configures, and returns a *http.Client type.
func (surf *Surf) buildClient(req *Request) *http.Client {
	client := &http.Client{
		CheckRedirect: req.checkRedirect,
	}

	if req.EnableCookie {
		client.Jar = surf.cookieJar
	}

	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, req.DialTimeout)
			if err != nil {
				return nil, err
			}
			if req.ConnTimeout > 0 {
				c.SetDeadline(time.Now().Add(req.ConnTimeout))
			}
			return c, nil
		},
	}

	if req.proxy != nil {
		transport.Proxy = http.ProxyURL(req.proxy)
	}

	if strings.ToLower(req.url.Scheme) == "https" {
		transport.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
		transport.DisableCompression = true
	}
	client.Transport = transport
	return client
}

// send uses the given *http.Request to make an HTTP request.
func (surf *Surf) httpRequest(param *Request) (resp *http.Response, err error) {
	req, err := http.NewRequest(param.Method, param.Url, param.body)
	if err != nil {
		return nil, err
	}

	req.Header = param.Header

	if param.TryTimes <= 0 {
		for {
			resp, err = param.client.Do(req)
			if err != nil {
				if !param.EnableCookie {
					l := len(UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", UserAgents["common"][r.Intn(l)])
				}
				time.Sleep(param.RetryPause)
				continue
			}
			break
		}
	} else {
		for i := 0; i < param.TryTimes; i++ {
			resp, err = param.client.Do(req)
			if err != nil {
				if !param.EnableCookie {
					l := len(UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", UserAgents["common"][r.Intn(l)])
				}
				time.Sleep(param.RetryPause)
				continue
			}
			break
		}
	}

	return resp, err
}

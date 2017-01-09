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

// surfer是一款Go语言编写的高并发爬虫下载器，支持 GET/POST/HEAD 方法及 http/https 协议，同时支持固定UserAgent自动保存cookie与随机大量UserAgent禁用cookie两种模式，高度模拟浏览器行为，可实现模拟登录等功能。

package surfer

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSurf(t *testing.T) {
	req := &Request{
		Method:       "GET",
		Url:          "https://www.bing.com/search?q=golang",
		EnableCookie: true,
		Header: http.Header{
			"Origin": []string{"https://cn.bing.com"},
		},
	}
	resp, _ := Download(req)
	b, _ := ioutil.ReadAll(resp.Body)
	t.Logf("request:\n%#v", req)
	t.Logf("response:\n%#v\nresponse_body:\n%s", resp, b[:200])
	resp, _ = Download(req)
	b, _ = ioutil.ReadAll(resp.Body)
	t.Logf("request:\n%#v", req)
	t.Logf("response:\n%#v\nresponse_body:\n%s", resp, b[:200])
}

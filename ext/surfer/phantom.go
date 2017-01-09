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
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type (
	// Phantom 基于Phantomjs的下载器实现，作为surfer的补充
	// 效率较surfer会慢很多，但是因为模拟浏览器，破防性更好
	// 支持UserAgent/TryTimes/RetryPause/自定义js
	Phantom struct {
		PhantomjsFile string            //Phantomjs完整文件名
		TempJsDir     string            //临时js存放目录
		jsFileMap     map[string]string //已存在的js文件
	}
	// Response 用于解析Phantomjs的响应内容
	Response struct {
		Cookies []string
		Body    string
	}
)

// NewPhantom 创建一个Phantomjs下载器
func NewPhantom(phantomjsFile, tempJsDir string) Surfer {
	phantom := &Phantom{
		PhantomjsFile: phantomjsFile,
		TempJsDir:     tempJsDir,
		jsFileMap:     make(map[string]string),
	}
	if !filepath.IsAbs(phantom.PhantomjsFile) {
		phantom.PhantomjsFile, _ = filepath.Abs(phantom.PhantomjsFile)
	}
	if !filepath.IsAbs(phantom.TempJsDir) {
		phantom.TempJsDir, _ = filepath.Abs(phantom.TempJsDir)
	}
	// 创建/打开目录
	err := os.MkdirAll(phantom.TempJsDir, 0777)
	if err != nil {
		log.Printf("[E] Surfer: %v\n", err)
		return phantom
	}
	phantom.createJsFile("js", js)
	return phantom
}

// Download 实现surfer下载器接口
func (phantom *Phantom) Download(req *Request) (resp *http.Response, err error) {
	err = req.prepare()
	if err != nil {
		return resp, err
	}
	var encoding = "utf-8"
	if _, params, err := mime.ParseMediaType(req.Header.Get("Content-Type")); err == nil {
		if cs, ok := params["charset"]; ok {
			encoding = strings.ToLower(strings.TrimSpace(cs))
		}
	}

	req.Header.Del("Content-Type")

	var b, _ = req.ReadBody()

	resp = req.writeback(resp)

	var args = []string{
		phantom.jsFileMap["js"],
		req.Url,
		req.Header.Get("Cookie"),
		encoding,
		req.Header.Get("User-Agent"),
		string(b),
		strings.ToLower(req.Method),
	}

	for i := 0; i < req.TryTimes; i++ {
		cmd := exec.Command(phantom.PhantomjsFile, args...)
		if resp.Body, err = cmd.StdoutPipe(); err != nil {
			time.Sleep(req.RetryPause)
			continue
		}
		err = cmd.Start()
		if err != nil || resp.Body == nil {
			time.Sleep(req.RetryPause)
			continue
		}
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			time.Sleep(req.RetryPause)
			continue
		}
		retResp := Response{}
		err = json.Unmarshal(b, &retResp)
		if err != nil {
			time.Sleep(req.RetryPause)
			continue
		}
		resp.Header = req.Header
		delete(resp.Header, "Set-Cookie")
		for _, c := range retResp.Cookies {
			resp.Header.Add("Set-Cookie", c)
		}
		resp.Body = ioutil.NopCloser(strings.NewReader(retResp.Body))
		break
	}

	if err == nil {
		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusOK)
	} else {
		resp.StatusCode = http.StatusBadGateway
		resp.Status = err.Error()
	}
	return resp, err
}

// DestroyJsFiles 销毁js临时文件
func (phantom *Phantom) DestroyJsFiles() {
	p, _ := filepath.Split(phantom.TempJsDir)
	if p == "" {
		return
	}
	for _, filename := range phantom.jsFileMap {
		os.Remove(filename)
	}
	if len(WalkDir(p)) == 1 {
		os.Remove(p)
	}
}

func (phantom *Phantom) createJsFile(fileName, jsCode string) {
	fullFileName := filepath.Join(phantom.TempJsDir, fileName)
	// 创建并写入文件
	f, _ := os.Create(fullFileName)
	f.Write([]byte(jsCode))
	f.Close()
	phantom.jsFileMap[fileName] = fullFileName
}

/*
* system.args[0] == js
* system.args[1] == url
* system.args[2] == cookie
* system.args[3] == pageEncode
* system.args[4] == userAgent
* system.args[5] == postdata
* system.args[6] == method
 */
const js string = `
var system = require('system');
var page = require('webpage').create();
var url = system.args[1];
var cookie = system.args[2];
var pageEncode = system.args[3];
var userAgent = system.args[4];
var postdata = system.args[5];
var method = system.args[6];
page.onResourceRequested = function(requestData, request) {
    request.setHeader('Cookie', cookie)
};
phantom.outputEncoding = pageEncode;
page.settings.userAgent = userAgent;
page.open(url, method, postdata, function(status) {
   if (status !== 'success') {
        console.log('Unable to access network');
    } else {
        var cookies = new Array();
        for(var i in page.cookies) {
        	var cookie = page.cookies[i];
        	var c = cookie["name"] + "=" + cookie["value"];
        	for (var obj in cookie){
        		if(obj == 'name' || obj == 'value'){
        			continue;
        		}
				c +=  "; " +　obj + "=" +  cookie[obj];
    		}
			cookies[i] = c;
		}
        var resp = {
            "Cookies": cookies,
            "Body": page.content
        };
        console.log(JSON.stringify(resp));
    }
    phantom.exit();
});
`

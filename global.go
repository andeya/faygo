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
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/henrylee2cn/apiware"
	"github.com/henrylee2cn/thinkgo/acceptencoder"
	"github.com/henrylee2cn/thinkgo/logging"
	"github.com/henrylee2cn/thinkgo/utils"
)

// GlobalSetting defines the global configuration and functions...
type GlobalSetting struct {
	config GlobalConfig
	// Error replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further
	// writes are done to response.
	// The error message should be plain text.
	errorFunc ErrorFunc
	// Decode params from request body.
	bodyDecodeFunc apiware.BodyDecodeFunc
	// The following is only for the APIHandler
	bindErrorFunc BindErrorFunc
	// When the APIHander's parameter name (struct tag) is unsetted,
	// it is mapped from the structure field name by default.
	// If `paramMapping` is nil, use snake style.
	// If the APIHander's parameter binding fails, the default handler is invoked
	paramMapping apiware.ParamNameFunc
	// global file cache system manager
	fsManager    *FileServerManager
	pongo2Render *Pongo2Render
	// The path for the upload files
	uploadDir string
	// The path for the static files
	staticDir string
	// The path for the log files
	logDir string

	syslog *logging.Logger
	bizlog *logging.Logger
}

// global configuration and functions...
var Global = func() *GlobalSetting {
	global := &GlobalSetting{
		config:         globalConfig,
		errorFunc:      defaultErrorFunc,
		bodyDecodeFunc: defaultBodyJSONFunc,
		bindErrorFunc:  defaultBindErrorFunc,
		paramMapping:   defaultParamMapping,
		fsManager: newFileServerManager(
			globalConfig.Cache.SizeMB*1024*1024,
			globalConfig.Cache.Expire,
			globalConfig.Cache.Enable,
			globalConfig.Gzip.Enable,
		),
		uploadDir: defaultUploadDir,
		staticDir: defaultStaticDir,
		logDir:    defaultLogDir,
	}
	if globalConfig.Cache.Enable {
		global.pongo2Render = newPongo2Render(global.fsManager.OpenFile)
	} else {
		global.pongo2Render = newPongo2Render(nil)
	}
	global.initLogger()
	return global
}()

var (
	defaultErrorFunc = func(ctx *Context, errStr string, status int) {
		ctx.Log().Error(errStr)
		statusText := http.StatusText(status)
		if len(errStr) > 0 {
			errStr = `<br><p><b style="color:red;">[ERROR]</b> <pre>` + errStr + `</pre></p>`
		}
		ctx.W.Header().Set(HeaderXContentTypeOptions, nosniff)
		ctx.HTML(status, fmt.Sprintf("<html>\n"+
			"<head><title>%d %s</title></head>\n"+
			"<body bgcolor=\"white\">\n"+
			"<center><h1>%d %s</h1></center>\n"+
			"<hr>\n<center>thinkgo/%s</center>\n%s\n</body>\n</html>\n",
			status, statusText, status, statusText, VERSION, errStr),
		)
	}
	// The default body decoder is json format decoding
	defaultBodyJSONFunc = func(dest reflect.Value, body []byte) error {
		var err error
		if dest.Kind() == reflect.Ptr {
			err = json.Unmarshal(body, dest.Interface())
		} else {
			err = json.Unmarshal(body, dest.Addr().Interface())
		}
		return err
	}
	defaultBindErrorFunc = func(ctx *Context, err error) {
		ctx.String(http.StatusBadRequest, "%v", err)
	}
	defaultParamMapping = utils.SnakeString
	// The default path for the upload files
	defaultUploadDir = "./upload/"
	// The default path for the static files
	defaultStaticDir = "./static/"
	// The default path for the log files
	defaultLogDir = "./log/"
)

func init() {
	fmt.Println(banner[1:])
	Global.syslog.Criticalf("The PID of the current process is %d", os.Getpid())
	if Global.config.warnMsg != "" {
		Warning(Global.config.warnMsg)
		Global.config.warnMsg = ""
	}
	// init file cache
	acceptencoder.InitGzip(Global.config.Gzip.MinLength, Global.config.Gzip.CompressLevel, Global.config.Gzip.Methods)
}

// When an error occurs, the default handler is invoked.
func (global *GlobalSetting) Error(ctx *Context, errStr string, status int) {
	global.errorFunc(ctx, errStr, status)
}

// Set the global default `ErrorFunc` function.
func (global *GlobalSetting) SetErrorFunc(errorFunc ErrorFunc) {
	global.errorFunc = errorFunc
}

// Decode params from request body.
func (global *GlobalSetting) BodyDecode(dest reflect.Value, body []byte) error {
	return global.bodyDecodeFunc(dest, body)
}

// Set the global default `BodyDecodeFunc` function.
func (global *GlobalSetting) SetBodyDecodeFunc(bodyDecodeFunc apiware.BodyDecodeFunc) {
	global.bodyDecodeFunc = bodyDecodeFunc
}

// If the APIHander's parameter binding fails, the default handler is invoked.
func (global *GlobalSetting) BindError(ctx *Context, err error) {
	global.bindErrorFunc(ctx, err)
}

// Set the global default `BindErrorFunc` function.
func (global *GlobalSetting) SetBindErrorFunc(bindErrorFunc BindErrorFunc) {
	global.bindErrorFunc = bindErrorFunc
}

// When the APIHander's parameter name (struct tag) is unsetted,
// it is mapped from the structure field name by default.
// If `paramMapping` is nil, use snake style.
func (global *GlobalSetting) ParamMapping(fieldName string) (paramName string) {
	return global.paramMapping(fieldName)
}

// Set the global default `ParamNameFunc` function.
func (global *GlobalSetting) SetParamMapping(paramMapping apiware.ParamNameFunc) {
	global.paramMapping = paramMapping
}

// Sets the global template variable or function for pongo2 render.
func (global *GlobalSetting) TemplateVariable(name string, v interface{}) {
	global.pongo2Render.TemplateVariable(name, v)
}

// UploadDir returns logs folder path with a slash at the end
func (global *GlobalSetting) LogDir() string {
	return global.logDir
}

// SetUpload sets upload folder path such as `./upload/`
// note: it should be called before Run()
// with a slash at the end
func (global *GlobalSetting) SetUpload(dir string) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	global.uploadDir = dir
}

// UploadDir returns upload folder path with a slash at the end
func (global *GlobalSetting) UploadDir() string {
	return global.uploadDir
}

// SetStatic sets static folder path, such as `./staic/`
// note: it should be called before Run()
// with a slash `/` at the end
func (global *GlobalSetting) SetStatic(dir string) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	global.staticDir = dir
}

// StaticDir returns static folder path with a slash at the end
func (global *GlobalSetting) StaticDir() string {
	return global.staticDir
}

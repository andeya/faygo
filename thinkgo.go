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
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/henrylee2cn/thinkgo/logging"
	"github.com/henrylee2cn/thinkgo/logging/color"
	"github.com/henrylee2cn/thinkgo/session"
	"github.com/henrylee2cn/thinkgo/swagger"
)

const (
	// VERSION is thinkgo web framework's version
	VERSION = "0.5"
	banner  = `
   _______  _                _                    
  |__   __|| |    [ ]       | |                   
     | |   | |_    _    _   | |  _   ___    ___   
     | |   |  _ \ | | / _ \ | |/ /  / _ \  / _ \  
     | |   | | | || || | | ||   -  | |_| || |_| | 
     |_|   |_| |_||_||_| |_||_| \_\ \_  /  \___/  
                                    _ \ \         
                                    \_\_/         ` + VERSION + "\n"
)

// Framework is the thinkgo web framework.
type Framework struct {
	name           string // name of the application
	version        string // version of the application
	config         Config
	*MuxAPI        // root muxAPI node
	muxesForRouter MuxAPIs
	filter         HandlerChain // called before the route is matched
	servers        []*Server
	once           sync.Once
	sessionManager *session.Manager
	syslog         *logging.Logger // for framework
	bizlog         *logging.Logger // for user bissness
	apidoc         *swagger.Swagger
}

// New uses the thinkgo web framework to create a new application.
func New(name string, version ...string) *Framework {
	mutexNewApp.Lock()
	defer mutexNewApp.Unlock()
	configFileName, ver := createConfigFilenameAndVersion(name, version...)
	frame := &Framework{
		name:           name,
		version:        ver,
		muxesForRouter: nil,
		config:         newConfig(configFileName),
	}
	frame.initSysLogger()
	frame.initBizLogger()
	frame.MuxAPI = newMuxAPI(frame, "root", "", "/")

	id := frame.NameWithVersion()
	if _, ok := Apps[id]; ok {
		Fatalf("There are two applications with exactly the same name and version: %s", id)
	}

	Apps[frame.NameWithVersion()] = frame

	return frame
}

var (
	// Apps is the list of applications that have been created.
	Apps          = make(map[string]*Framework)
	mutexNewApp   sync.Mutex
	mutexForBuild sync.Mutex
)

// Name returns the name of the application
func (frame *Framework) Name() string {
	return frame.name
}

// Version returns the version of the application
func (frame *Framework) Version() string {
	return frame.version
}

// NameWithVersion returns the name with version
func (frame *Framework) NameWithVersion() string {
	if len(frame.version) == 0 {
		return frame.name
	}
	return frame.name + "_" + frame.version
}

// Run starts web services.
func (frame *Framework) Run() {
	frame.once.Do(func() {
		frame.build()
		last := len(frame.servers) - 1
		for i := 0; i < last; i++ {
			go frame.servers[i].run()
		}
		frame.servers[last].run()
	})
}

func (frame *Framework) build() {
	// Make sure that the initialization logs for multiple applications are printed in sequence
	mutexForBuild.Lock()
	defer mutexForBuild.Unlock()

	// register the default MuxAPIs
	{
		// apidoc
		if frame.config.APIdoc.Enable {
			frame.regAPIdoc()
		}
		// static
		frame.presetSystemMuxes()
	}

	// build router
	var router = &Router{
		RedirectTrailingSlash:  frame.config.Router.RedirectTrailingSlash,
		RedirectFixedPath:      frame.config.Router.RedirectFixedPath,
		HandleMethodNotAllowed: frame.config.Router.HandleMethodNotAllowed,
		HandleOPTIONS:          frame.config.Router.HandleOPTIONS,
		filter:                 frame.makeFilterHandle(),
		NotFound:               frame.makeErrorHandler(http.StatusNotFound),
		MethodNotAllowed:       frame.makeErrorHandler(http.StatusMethodNotAllowed),
		PanicHandler:           frame.makePanicHandler(),
	}

	// register router
	for _, node := range frame.MuxAPIsForRouter() {
		handle := frame.makeHandle(node.handlers)
		for _, method := range node.methods {
			frame.syslog.Criticalf("%7s | %-30s", method, node.path)
			router.Handle(method, node.path, handle)
		}
	}

	// new server
	nameWithVersion := frame.NameWithVersion()
	for i, netType := range frame.config.NetTypes {
		frame.servers = append(frame.servers, &Server{
			nameWithVersion: nameWithVersion,
			netType:         netType,
			tlsCertFile:     frame.config.TLSCertFile,
			tlsKeyFile:      frame.config.TLSKeyFile,
			letsencryptFile: frame.config.LetsencryptFile,
			unixFileMode:    frame.config.UNIXFileMode,
			Server: &http.Server{
				Addr:         frame.config.Addrs[i],
				Handler:      router,
				ReadTimeout:  frame.config.ReadTimeout,
				WriteTimeout: frame.config.WriteTimeout,
			},
			log: frame.syslog,
		})
	}

	// register session
	frame.registerSession()
}

// Log returns the log used by the user bissness
func (frame *Framework) Log() *logging.Logger {
	return frame.bizlog
}

// MuxAPIsForRouter get an ordered list of nodes used to register router.
func (frame *Framework) MuxAPIsForRouter() []*MuxAPI {
	if frame.muxesForRouter == nil {
		// comb mux.handlers, mux.paramInfos, mux.returns and mux.path,.
		frame.MuxAPI.comb()

		frame.muxesForRouter = frame.MuxAPI.HandlerProgeny()
	}
	return frame.muxesForRouter
}

// Filter operations that are called before the route is matched.
func (frame *Framework) Filter(fn ...HandlerFunc) *Framework {
	handlers := make([]Handler, len(fn))
	for i, h := range fn {
		handlers[i] = h
	}
	frame.filter = append(handlers, frame.filter...)
	return frame
}

// Route append middlewares of function type to root muxAPI.
// Used to register router in tree style.
func (frame *Framework) Route(children ...*MuxAPI) *MuxAPI {
	frame.MuxAPI.children = append(frame.MuxAPI.children, children...)
	for _, child := range children {
		child.parent = frame.MuxAPI
	}
	return frame.MuxAPI
}

// NewGroup create an isolated grouping muxAPI node.
func (frame *Framework) NewGroup(pattern string, children ...*MuxAPI) *MuxAPI {
	return frame.NewNamedGroup("", pattern, children...)
}

// NewAPI creates an isolated muxAPI node.
func (frame *Framework) NewAPI(methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI("", methodset, pattern, handlers...)
}

// NewNamedGroup creates an isolated grouping muxAPI node with the name.
func (frame *Framework) NewNamedGroup(name string, pattern string, children ...*MuxAPI) *MuxAPI {
	group := frame.NewNamedAPI(name, "", pattern)
	group.children = append(group.children, children...)
	for _, child := range children {
		child.parent = group
	}
	return group
}

// NewNamedAPI creates an isolated muxAPI node with the name.
func (frame *Framework) NewNamedAPI(name string, methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return newMuxAPI(frame, name, methodset, pattern, handlers...)
}

// NewGET is a shortcut for frame.NewAPI("GET", pattern, handlers...)
func (frame *Framework) NewGET(pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewAPI("GET", pattern, handlers...)
}

// NewHEAD is a shortcut for frame.NewAPI("HEAD", pattern, handlers...)
func (frame *Framework) NewHEAD(pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewAPI("HEAD", pattern, handlers...)
}

// NewOPTIONS is a shortcut for frame.NewAPI("OPTIONS", pattern, handlers...)
func (frame *Framework) NewOPTIONS(pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewAPI("OPTIONS", pattern, handlers...)
}

// NewPOST is a shortcut for frame.NewAPI("POST", pattern, handlers...)
func (frame *Framework) NewPOST(pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewAPI("POST", pattern, handlers...)
}

// NewPUT is a shortcut for frame.NewAPI("PUT", pattern, handlers...)
func (frame *Framework) NewPUT(pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewAPI("PUT", pattern, handlers...)
}

// NewPATCH is a shortcut for frame.NewAPI("PATCH", pattern, handlers...)
func (frame *Framework) NewPATCH(pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewAPI("PATCH", pattern, handlers...)
}

// NewDELETE is a shortcut for frame.NewAPI("DELETE", pattern, handlers...)
func (frame *Framework) NewDELETE(pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewAPI("DELETE", pattern, handlers...)
}

// NewNamedGET is a shortcut for frame.NewNamedAPI(name, "GET", pattern, handlers...)
func (frame *Framework) NewNamedGET(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI(name, "GET", pattern, handlers...)
}

// NewNamedHEAD is a shortcut for frame.NewNamedAPI(name, "HEAD", pattern, handlers...)
func (frame *Framework) NewNamedHEAD(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI(name, "HEAD", pattern, handlers...)
}

// NewNamedOPTIONS is a shortcut for frame.NewNamedAPI(name, "OPTIONS", pattern, handlers...)
func (frame *Framework) NewNamedOPTIONS(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI(name, "OPTIONS", pattern, handlers...)
}

// NewNamedPOST is a shortcut for frame.NewNamedAPI(name, "POST", pattern, handlers...)
func (frame *Framework) NewNamedPOST(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI(name, "POST", pattern, handlers...)
}

// NewNamedPUT is a shortcut for frame.NewNamedAPI(name, "PUT", pattern, handlers...)
func (frame *Framework) NewNamedPUT(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI(name, "PUT", pattern, handlers...)
}

// NewNamedPATCH is a shortcut for frame.NewNamedAPI(name, "PATCH", pattern, handlers...)
func (frame *Framework) NewNamedPATCH(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI(name, "PATCH", pattern, handlers...)
}

// NewNamedDELETE is a shortcut for frame.NewNamedAPI(name, "DELETE", pattern, handlers...)
func (frame *Framework) NewNamedDELETE(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NewNamedAPI(name, "DELETE", pattern, handlers...)
}

// NewStatic creates an isolated static muxAPI node.
func (frame *Framework) NewStatic(pattern string, root string, nocompressAndNocache ...bool) *MuxAPI {
	return frame.NewNamedStatic("", pattern, root, nocompressAndNocache...)
}

// NewNamedStatic creates an isolated static muxAPI node with the name.
func (frame *Framework) NewNamedStatic(name, pattern string, root string, nocompressAndNocache ...bool) *MuxAPI {
	return (&MuxAPI{frame: frame}).NamedStatic(name, pattern, root, nocompressAndNocache...)
}

// NewStaticFS creates an isolated static muxAPI node.
func (frame *Framework) NewStaticFS(pattern string, fs FileSystem) *MuxAPI {
	return frame.NewNamedStaticFS("", pattern, fs)
}

// NewNamedStaticFS creates an isolated static muxAPI node with the name.
func (frame *Framework) NewNamedStaticFS(name, pattern string, fs FileSystem) *MuxAPI {
	return (&MuxAPI{frame: frame}).NamedStaticFS(name, pattern, fs)
}

// makeFilterHandle makes an FilterFunc.
func (frame *Framework) makeFilterHandle() FilterFunc {
	if len(frame.filter) == 0 {
		return nil
	}
	ctxPool := sync.Pool{
		New: func() interface{} {
			return newFilterContext(frame)
		},
	}
	return func(w http.ResponseWriter, r *http.Request) (map[interface{}]interface{}, bool) {
		ctx := ctxPool.Get().(*Context)
		ctx.reset(w, r, nil, nil)
		defer func() {
			ctxPool.Put(ctx)
		}()
		ctx.posReset()

		var u = ctx.URI()
		start := time.Now()

		ctx.Next()
		if ctx.IsBreak() {
			if !ctx.W.Committed() {
				ctx.Error(http.StatusForbidden, http.StatusText(http.StatusForbidden))
			}
			stop := time.Now()
			method := ctx.Method()
			if u == "" {
				u = "/"
			}
			n := ctx.W.Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Magenta(n)
			case n >= 300:
				code = color.Grey(n)
			}
			ctx.Log().Infof("%15s %7s  %3s %10d %12s %-30s | ", ctx.RealIP(), method, code, ctx.W.Size(), stop.Sub(start), u)
			return nil, false
		}
		return ctx.data, true
	}
}

// makeHandle makes an *apiware.ParamsAPI implements the Handle interface.
func (frame *Framework) makeHandle(handlerChain HandlerChain) Handle {
	ctxPool := sync.Pool{
		New: func() interface{} {
			return newContext(frame, handlerChain)
		},
	}
	return func(w http.ResponseWriter, r *http.Request, pathParams Params, data map[interface{}]interface{}) {
		ctx := ctxPool.Get().(*Context)
		ctx.reset(w, r, pathParams, data)
		defer func() {
			ctxPool.Put(ctx)
		}()
		ctx.do()
	}
}

// Create the handle to be called by the router
func (frame *Framework) makeErrorHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Global.errorFunc(newEmptyContext(frame, w, r), http.StatusText(status), status)
	})
}

// Create the handle to be called by the router
func (frame *Framework) makePanicHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	return func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		stack := make([]byte, 4<<10) //4KB
		length := runtime.Stack(stack, true)
		start := bytes.Index(stack, s)
		stack = stack[start:length]
		start = bytes.Index(stack, line) + 1
		stack = stack[start:]
		end := bytes.LastIndex(stack, line)
		if end != -1 {
			stack = stack[:end]
		}
		end = bytes.Index(stack, e)
		if end != -1 {
			stack = stack[:end]
		}
		stack = bytes.TrimRight(stack, "\n")
		errStr := fmt.Sprintf("%v\n[TRACE]\n%s\n", rcv, stack)
		Global.errorFunc(newEmptyContext(frame, w, r), errStr, http.StatusInternalServerError)
	}
}

func (frame *Framework) presetSystemMuxes() {
	frame.Use(accessLogWare())
	var hadUpload, hadStatic bool
	for _, child := range frame.MuxAPI.children {
		if strings.Contains(child.pattern, "/upload/") {
			hadUpload = true
		}
		if strings.Contains(child.pattern, "/static/") {
			hadUpload = true
		}
	}
	// When does not have a custom route, the route is automatically created.
	if !hadUpload {
		frame.MuxAPI.NamedStatic(
			"Directory for uploading files",
			"/upload/",
			Global.upload.root,
			Global.upload.nocompress,
			Global.upload.nocache,
		).Use(Global.upload.handlers...)
	}
	if !hadStatic {
		frame.MuxAPI.NamedStatic(
			"Directory for public static files",
			"/static/",
			Global.static.root,
			Global.static.nocompress,
			Global.static.nocache,
		).Use(Global.static.handlers...)
	}
}

func (frame *Framework) registerSession() {
	if !frame.config.Session.Enable {
		return
	}
	conf := &session.ManagerConfig{
		CookieName:              frame.config.Session.Name,
		EnableSetCookie:         frame.config.Session.AutoSetCookie,
		Gclifetime:              frame.config.Session.GCMaxLifetime,
		Secure:                  true,
		CookieLifeTime:          frame.config.Session.CookieLifetime,
		ProviderConfig:          frame.config.Session.ProviderConfig,
		Domain:                  frame.config.Session.Domain,
		EnableSidInHttpHeader:   frame.config.Session.EnableSidInHttpHeader,
		SessionNameInHttpHeader: frame.config.Session.NameInHttpHeader,
		EnableSidInUrlQuery:     frame.config.Session.EnableSidInUrlQuery,
	}
	var err error
	frame.sessionManager, err = session.NewManager(frame.config.Session.Provider, conf)
	if err != nil {
		panic(err)
	}
	go frame.sessionManager.GC()
}

func createConfigFilenameAndVersion(name string, version ...string) (fileName string, ver string) {
	if len(version) > 0 && len(version[0]) > 0 {
		ver = version[0]
		fileName = CONFIG_DIR + "/" + name + "_" + ver + ".ini"
	} else {
		fileName = CONFIG_DIR + "/" + name + ".ini"
	}
	return
}

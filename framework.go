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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/henrylee2cn/faygo/logging"
	"github.com/henrylee2cn/faygo/logging/color"
	"github.com/henrylee2cn/faygo/session"
	"github.com/henrylee2cn/faygo/swagger"
	"github.com/henrylee2cn/faygo/utils"
)

// Framework is the faygo web framework.
type Framework struct {
	// name of the application
	name string
	// version of the application
	version string
	config  Config
	// root muxAPI node
	*MuxAPI
	muxesForRouter MuxAPIs
	// called before the route is matched
	filter         HandlerChain
	servers        []*Server
	running        bool
	buildOnce      sync.Once
	lock           sync.RWMutex
	sessionManager *session.Manager
	// for framework
	syslog *logging.Logger
	// for user bissness
	bizlog *logging.Logger
	apidoc *swagger.Swagger
	trees  map[string]*node
	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	redirectTrailingSlash bool
	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// redirectTrailingSlash is independent of this option.
	redirectFixedPath bool
	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	handleMethodNotAllowed bool
	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	handleOPTIONS bool
	contextPool   sync.Pool
}

// Make sure the Framework conforms with the http.Handler interface
var _ http.Handler = new(Framework)

// newFramework uses the faygo web framework to create a new application.
func newFramework(name string, version ...string) *Framework {
	mutexNewApp.Lock()
	defer mutexNewApp.Unlock()
	var frame = new(Framework)

	frame.name = strings.TrimSpace(name)
	if len(version) > 0 && len(version[0]) > 0 {
		frame.version = strings.TrimSpace(version[0])
	}

	id := frame.NameWithVersion()
	if _, ok := GetFrame(id); ok {
		Fatalf("There are two applications with exactly the same name and version: %s", id)
	}

	configFilename := frame.ConfigFilename()
	frame.setConfig(newConfig(configFilename))

	frame.redirectTrailingSlash = frame.config.Router.RedirectTrailingSlash
	frame.redirectFixedPath = frame.config.Router.RedirectFixedPath
	frame.handleMethodNotAllowed = frame.config.Router.HandleMethodNotAllowed
	frame.handleOPTIONS = frame.config.Router.HandleOPTIONS
	frame.contextPool = sync.Pool{
		New: func() interface{} {
			ctx := &Context{
				frame:         frame,
				enableGzip:    global.config.Gzip.Enable,
				enableSession: frame.config.Session.Enable,
				enableXSRF:    frame.config.XSRF.Enable,
			}
			ctx.W = &Response{context: ctx}
			return ctx
		},
	}
	frame.initSysLogger()
	frame.initBizLogger()
	frame.MuxAPI = newMuxAPI(frame, "root", "", "/")

	addFrame(frame)
	return frame
}

var (
	mutexNewApp   sync.Mutex
	mutexForBuild sync.Mutex
)

func (frame *Framework) setConfig(config Config) {
	frame.config = config
}

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

// Config returns the framework's config copy.
func (frame *Framework) Config() Config {
	return frame.config
}

// ConfigFilename returns the framework's config file name.
func (frame *Framework) ConfigFilename() string {
	return CONFIG_DIR + "/" + frame.NameWithVersion() + ".ini"
}

// Run starts the web service.
func (frame *Framework) Run() {
	if frame.Running() {
		return
	}
	go frame.run()
	global.graceOnce.Do(func() {
		graceSignal()
	})
	select {}
}

// Running returns whether the frame service is running.
func (frame *Framework) Running() bool {
	frame.lock.RLock()
	defer frame.lock.RUnlock()
	return frame.running
}

func (frame *Framework) run() {
	frame.lock.Lock()
	frame.build()
	frame.running = true
	count := len(frame.servers)
	for i := 0; i < count; i++ {
		go frame.servers[i].run()
	}
	frame.lock.Unlock()
}

func (frame *Framework) build() {
	frame.buildOnce.Do(func() {
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

		// register router
		for _, api := range frame.MuxAPIsForRouter() {
			handle := frame.makeHandle(api.handlers)
			for _, method := range api.methods {
				if api.path[0] != '/' {
					Panic("path must begin with '/' in path '" + api.path + "'")
				}
				if frame.trees == nil {
					frame.trees = make(map[string]*node)
				}
				root := frame.trees[method]
				if root == nil {
					root = new(node)
					frame.trees[method] = root
				}
				root.addRoute(api.path, handle)
				frame.syslog.Criticalf("\x1b[46m[SYS]\x1b[0m %7s | %-30s", method, api.path)
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
				letsencryptDir:  frame.config.LetsencryptDir,
				unixFileMode:    frame.config.UNIXFileMode,
				Server: &http.Server{
					Addr:         frame.config.Addrs[i],
					Handler:      frame,
					ReadTimeout:  frame.config.ReadTimeout,
					WriteTimeout: frame.config.WriteTimeout,
				},
				log: frame.syslog,
			})
		}

		// register session
		frame.registerSession()
	})
}

// shutdown closes the frame service gracefully.
func (frame *Framework) shutdown(ctxTimeout context.Context) (graceful bool) {
	frame.lock.Lock()
	defer frame.lock.Unlock()
	if !frame.running {
		return true
	}
	var flag int32 = 1
	count := new(sync.WaitGroup)
	for _, server := range frame.servers {
		count.Add(1)
		go func(srv *Server) {
			if err := srv.Shutdown(ctxTimeout); err != nil {
				atomic.StoreInt32(&flag, 0)
				frame.Log().Errorf("[shutdown-%s] %s", frame.NameWithVersion(), err.Error())
			}
			count.Done()
		}(server)
	}
	count.Wait()
	frame.running = false
	frame.CloseLog()
	return flag == 1
}

// Log returns the logger used by the user bissness.
func (frame *Framework) Log() *logging.Logger {
	return frame.bizlog
}

// CloseLog closes loggers.
func (frame *Framework) CloseLog() {
	frame.bizlog.Close()
	frame.syslog.Close()
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

func (frame *Framework) presetSystemMuxes() {
	var hadUpload, hadStatic bool
	for _, child := range frame.MuxAPI.children {
		if strings.Contains(child.pattern, "/upload/") {
			hadUpload = true
		}
		if strings.Contains(child.pattern, "/static/") {
			hadStatic = true
		}
	}
	// When does not have a custom route, the route is automatically created.
	if !hadUpload {
		frame.MuxAPI.NamedStatic(
			"Directory for uploading files",
			"/upload/",
			global.upload.root,
			global.upload.nocompress,
			global.upload.nocache,
		).Use(global.upload.handlers...)
	}
	if !hadStatic {
		frame.MuxAPI.NamedStatic(
			"Directory for public static files",
			"/static/",
			global.static.root,
			global.static.nocompress,
			global.static.nocache,
		).Use(global.static.handlers...)
	}
}

func (frame *Framework) registerSession() {
	if !frame.config.Session.Enable {
		return
	}
	conf := &session.ManagerConfig{
		CookieName:              frame.config.Session.Name,
		EnableSetCookie:         frame.config.Session.AutoSetCookie,
		CookieLifeTime:          frame.config.Session.CookieLifetime,
		Gclifetime:              frame.config.Session.GCLifetime,
		Maxlifetime:             frame.config.Session.MaxLifetime,
		Secure:                  true,
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

// ServeHTTP makes the router implement the http.Handler interface.
func (frame *Framework) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var start = time.Now()
	var ctx = frame.getContext(w, req)
	defer func() {
		if rcv := recover(); rcv != nil {
			panicHandler(ctx, rcv)
		}
		frame.putContext(ctx)
	}()
	var method = ctx.Method()
	var u = ctx.URI()
	if u == "" {
		u = "/"
	}
	frame.serveHTTP(ctx)
	var n = ctx.Status()
	var code string
	switch {
	case n >= 500:
		code = color.Red(n)
	case n >= 400:
		code = color.Magenta(n)
	case n >= 300:
		code = color.Grey(n)
	default:
		code = color.Green(n)
	}
	frame.syslog.Infof("[I] %15s %7s  %3s %10d %12s %-30s | ", ctx.RealIP(), method, code, ctx.Size(), time.Since(start), u)
}

func (frame *Framework) serveHTTP(ctx *Context) {
	if !ctx.doFilter() {
		return
	}
	var path = ctx.Path()
	var method = ctx.Method()
	if root := frame.trees[method]; root != nil {
		if handle, ps, tsr := root.getValue(path); handle != nil {
			handle(ctx, ps)
			return
		} else if method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && frame.redirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					ctx.ModifyPath(path[:len(path)-1])
				} else {
					ctx.ModifyPath(path + "/")
				}
				http.Redirect(ctx.W, ctx.R, ctx.URL().String(), code)
				return
			}

			// Try to fix the request path
			if frame.redirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					utils.CleanPath(path),
					frame.redirectTrailingSlash,
				)
				if found {
					ctx.ModifyPath(string(fixedPath))
					http.Redirect(ctx.W, ctx.R, ctx.URL().String(), code)
					return
				}
			}
		}
	}

	if method == "OPTIONS" {
		// Handle OPTIONS requests
		if frame.handleOPTIONS {
			if allow := frame.allowed(path, method); len(allow) > 0 {
				ctx.SetHeader("Allow", allow)
				return
			}
		}
	} else {
		// Handle 405
		if frame.handleMethodNotAllowed {
			if allow := frame.allowed(path, method); len(allow) > 0 {
				ctx.SetHeader("Allow", allow)
				global.errorFunc(ctx, "Method Not Allowed", 405)
				return
			}
		}
	}

	// Handle 404
	global.errorFunc(ctx, "Not Found", 404)
}

func (frame *Framework) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range frame.trees {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range frame.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _, _ := frame.trees[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

// makeHandle makes an *apiware.ParamsAPI implements the Handle interface.
func (frame *Framework) makeHandle(handlerChain HandlerChain) Handle {
	return func(ctx *Context, pathParams PathParams) {
		ctx.doHandler(handlerChain, pathParams)
	}
}

func panicHandler(ctx *Context, rcv interface{}) {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
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
	global.errorFunc(ctx, fmt.Sprintf("%v\n[TRACE]\n%s\n", rcv, stack), http.StatusInternalServerError)
}

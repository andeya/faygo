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
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/henrylee2cn/apiware"
	"github.com/henrylee2cn/thinkgo/acceptencoder"
	"github.com/henrylee2cn/thinkgo/session"
	"github.com/henrylee2cn/thinkgo/swagger"
	"github.com/henrylee2cn/thinkgo/utils"
	"github.com/henrylee2cn/thinkgo/utils/errors"
	"github.com/rsc/letsencrypt"
	// "github.com/facebookgo/grace/gracehttp"
)

const (
	VERSION = "0.1"
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

// Thinkgo web framework.
type Framework struct {
	name           string // name of the application
	version        string // version of the application
	config         Config
	mux            *MuxAPI
	muxesForRouter MuxAPIs
	server         *http.Server
	// Error replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further
	// writes are done to response.
	// The error message should be plain text.
	errorFunc         ErrorFunc
	fileServerManager *FileServerManager
	once              sync.Once
	// The following is only for the APIHandler
	bindErrorFunc BindErrorFunc
	// When the APIHander's parameter name (struct tag) is unsetted,
	// it is mapped from the structure field name by default.
	// If `paramMapping` is nil, use snake style.
	// If the APIHander's parameter binding fails, the default handler is invoked
	paramMapping   apiware.ParamNameFunc
	sessionManager *session.Manager
	apidoc         *swagger.Swagger
}

// Use the thinkgo web framework to create a new application.
func New(name string, version ...string) *Framework {
	configFileName, ver := createConfigFilenameAndVersion(name, version...)
	frame := &Framework{
		name:              name,
		version:           ver,
		muxesForRouter:    nil,
		config:            newConfig(configFileName),
		fileServerManager: new(FileServerManager),
	}
	frame.mux = newMuxAPI(frame, "root", "", "/", nil)
	return frame
}

// Quick to use.
var (
	defaultName      = "myapp"
	defaultFramework = New(defaultName)
	defaultErrorFunc = func(ctx *Context, errStr string, status int) {
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

	defaultBindErrorFunc = func(ctx *Context, err error) {
		ctx.String(http.StatusBadRequest, "%v", err)
	}

	initOnce   = new(sync.Once)
	bannerOnce = new(sync.Once)
)

// Initializes the name and version of the default application,
// returns the default application that was created.
func Init(nameAndVersion ...string) *Framework {
	initOnce.Do(func() {
		count := len(nameAndVersion)
		if count == 0 {
			return
		}
		var name = nameAndVersion[0]
		var version string
		if count > 1 {
			version = nameAndVersion[1]
		}
		configFileName, _ := createConfigFilenameAndVersion(name, version)
		if defaultFramework.name != defaultName || len(name) == 0 ||
			(defaultFramework.name == name && defaultFramework.version == version) {
			return
		}
		defaultFramework.config = newConfig(configFileName, defaultFramework.config.Addr)

		configFileName, _ = createConfigFilenameAndVersion(defaultFramework.name, defaultFramework.version)
		os.Remove(configFileName)

		defaultFramework.version = version
		defaultFramework.name = name
	})
	return defaultFramework
}

// name of the application
func (frame *Framework) Name() string {
	return frame.name
}

// version of the application
func (frame *Framework) Version() string {
	return frame.version
}

// Start web service.
func Run() {
	defaultFramework.Run()
}

// Start web service.
func (frame *Framework) Run() {
	bannerOnce.Do(func() { fmt.Println(banner[1:]) })
	frame.build()
	var err error
	switch frame.config.NetType {
	case NETTYPE_NORMAL:
		err = frame.listenAndServe()
	case NETTYPE_TLS:
		err = frame.listenAndServeTLS(frame.config.TLSCertFile, frame.config.TLSKeyFile)
	case NETTYPE_LETSENCRYPT:
		err = frame.listenAndServeLETSENCRYPT(frame.config.LetsencryptFile)
	case NETTYPE_UNIX:
		err = frame.listenAndServeUNIX(frame.config.UNIXFileMode)
	}
	if err != nil {
		panic(err)
	}
}

// listenAndServe listens on the TCP network address and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// If srv.Addr is blank, ":http" is used, listenAndServe always returns a non-nil error.
func (frame *Framework) listenAndServe() error {
	return frame.server.ListenAndServe()
}

// listenAndServeTLS listens on the TCP network address and
// then calls Serve to handle requests on incoming TLS connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Filenames containing a certificate and matching private key for the
// server must be provided if neither the Server's TLSConfig.Certificates
// nor TLSConfig.GetCertificate are populated. If the certificate is
// signed by a certificate authority, the certFile should be the
// concatenation of the server's certificate, any intermediates, and
// the CA's certificate.
//
// If frame.config.Addr is blank, ":https" is used, listenAndServeTLS always returns a non-nil error.
func (frame *Framework) listenAndServeTLS(certFile, keyFile string) error {
	return frame.server.ListenAndServeTLS(certFile, keyFile)
}

// listenAndServeLETSENCRYPT listens on a new Automatic TLS using letsencrypt.org service.
// if you want to disable cache file then simple give cacheFileOptional a value of empty string ""
func (frame *Framework) listenAndServeLETSENCRYPT(cacheFileOptional string) error {
	if frame.server.Addr == "" {
		frame.server.Addr = ":https"
	}

	ln, err := net.Listen("tcp4", frame.server.Addr)
	if err != nil {
		return err
	}

	var m letsencrypt.Manager
	if cacheFileOptional != "" {
		if err = m.CacheFile(cacheFileOptional); err != nil {
			return err
		}
	}

	tlsConfig := &tls.Config{GetCertificate: m.GetCertificate}
	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, tlsConfig)

	return frame.server.Serve(tlsListener)
}

var (
	errPortAlreadyUsed = errors.New("Port is already used")
	errRemoveUnix      = errors.New("Unexpected error when trying to remove unix socket file. Addr: %s | Trace: %s")
	errChmod           = errors.New("Cannot chmod %#o for %q: %s")
	errCertKeyMissing  = errors.New("You should provide certFile and keyFile for TLS/SSL")
	errParseTLS        = errors.New("Couldn't load TLS, certFile=%q, keyFile=%q. Trace: %s")
)

// listenAndServeUNIX announces on the Unix domain socket laddr and listens a Unix service.
func (frame *Framework) listenAndServeUNIX(fileMode os.FileMode) error {
	if errOs := os.Remove(frame.server.Addr); errOs != nil && !os.IsNotExist(errOs) {
		return errRemoveUnix.Format(frame.server.Addr, errOs.Error())
	}

	ln, err := net.Listen("unix", frame.server.Addr)
	if err != nil {
		return errPortAlreadyUsed.AppendErr(err)
	}

	if err = os.Chmod(frame.server.Addr, fileMode); err != nil {
		return errChmod.Format(fileMode, frame.server.Addr, err.Error())
	}
	return frame.server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (frame *Framework) build() {
	frame.once.Do(func() {
		// register the default MuxAPIs
		{
			// apidoc
			if frame.config.APIdoc.Enable {
				frame.regAPIdoc()
			}
			// static
			frame.registerSystemStatic()
		}

		if frame.errorFunc == nil {
			frame.errorFunc = defaultErrorFunc
		}
		if frame.bindErrorFunc == nil {
			frame.bindErrorFunc = defaultBindErrorFunc
		}
		if frame.paramMapping == nil {
			frame.paramMapping = utils.SnakeString
		}

		// build router
		var router = &Router{
			RedirectTrailingSlash:  frame.config.Router.RedirectTrailingSlash,
			RedirectFixedPath:      frame.config.Router.RedirectFixedPath,
			HandleMethodNotAllowed: frame.config.Router.HandleMethodNotAllowed,
			HandleOPTIONS:          frame.config.Router.HandleOPTIONS,
			NotFound:               frame.makeErrorHandler(http.StatusNotFound),
			MethodNotAllowed:       frame.makeErrorHandler(http.StatusMethodNotAllowed),
			PanicHandler:           frame.makePanicHandler(),
		}

		// register router
		for _, node := range frame.MuxAPIsForRouter() {
			handle := frame.makeHandle(node.handlers)
			for _, method := range node.methods {

				println("URL:" + node.path)

				router.Handle(method, node.path, handle)
			}
		}

		// new server
		frame.server = &http.Server{
			Addr:         frame.config.Addr,
			Handler:      router,
			ReadTimeout:  frame.config.ReadTimeout,
			WriteTimeout: frame.config.WriteTimeout,
		}

		// init file cache
		acceptencoder.InitGzip(frame.config.Gzip.MinLength, frame.config.Gzip.CompressLevel, frame.config.Gzip.Methods)
		*frame.fileServerManager = *newFileServerManager(
			frame.config.Cache.SizeMB*1024*1024,
			frame.config.Cache.Expire,
			frame.config.Cache.Enable,
			frame.config.Gzip.Enable,
			frame.errorFunc,
		)

		// register session
		frame.registerSession()
	})
}

// When an error occurs, the default handler is invoked.
func SetErrorFunc(errorFunc ErrorFunc) {
	defaultFramework.SetErrorFunc(errorFunc)
}

// When an error occurs, the default handler is invoked.
func (frame *Framework) SetErrorFunc(errorFunc ErrorFunc) {
	frame.errorFunc = errorFunc
}

// If the APIHander's parameter binding fails, the default handler is invoked.
func SetBindErrorFunc(bindErrorFunc BindErrorFunc) {
	defaultFramework.SetBindErrorFunc(bindErrorFunc)
}

// If the APIHander's parameter binding fails, the default handler is invoked.
func (frame *Framework) SetBindErrorFunc(bindErrorFunc BindErrorFunc) {
	frame.bindErrorFunc = bindErrorFunc
}

// When the APIHander's parameter name (struct tag) is unsetted,
// it is mapped from the structure field name by default.
// If `paramMapping` is nil, use snake style.
func SetParamMapping(paramMapping apiware.ParamNameFunc) {
	defaultFramework.SetParamMapping(paramMapping)
}

// When the APIHander's parameter name (struct tag) is unsetted,
// it is mapped from the structure field name by default.
// If `paramMapping` is nil, use snake style.
func (frame *Framework) SetParamMapping(paramMapping apiware.ParamNameFunc) {
	frame.paramMapping = paramMapping
}

// Root returns the root node of the muxAPI.
func Root() *MuxAPI {
	return defaultFramework.Root()
}

// Root returns the root node of the muxAPI.
func (frame *Framework) Root() *MuxAPI {
	return frame.mux
}

// Append middlewares of function type to root muxAPI.
func Route(children ...*MuxAPI) *MuxAPI {
	return defaultFramework.Route(children...)
}

// Append middlewares of function type to root muxAPI.
func (frame *Framework) Route(children ...*MuxAPI) *MuxAPI {
	frame.mux.children = append(frame.mux.children, children...)
	for _, child := range children {
		child.parent = frame.mux
	}
	return frame.mux
}

// Insert the middlewares at the left end of the node's handler chain.
func Use(handlers ...HandlerWithoutPath) *MuxAPI {
	return defaultFramework.Use(handlers...)
}

// Insert the middlewares at the left end of the node's handler chain.
func (frame *Framework) Use(handlers ...HandlerWithoutPath) *MuxAPI {
	return frame.mux.Use(handlers...)
}

// Create a new muxAPI node that contains the Handlers, but it is only group.
func Group(pattern string, children ...*MuxAPI) *MuxAPI {
	return defaultFramework.Group(pattern, children...)
}

// Create a new MuxAPI node that contains the Handlers, but it is only group.
func (frame *Framework) Group(pattern string, children ...*MuxAPI) *MuxAPI {
	return frame.NamedGroup("", pattern, children...)
}

// Create a new MuxAPI group node with a name.
func NamedGroup(name string, pattern string, children ...*MuxAPI) *MuxAPI {
	return defaultFramework.NamedGroup(name, pattern, children...)
}

// Create a new MuxAPI group node with a name.
func (frame *Framework) NamedGroup(name string, pattern string, children ...*MuxAPI) *MuxAPI {
	group := frame.NamedAPI(name, "", pattern)
	group.children = append(group.children, children...)
	for _, child := range children {
		child.parent = group
	}
	return group
}

// Create a new MuxAPI node that contains the Handlers.
func API(methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.API(methodset, pattern, handlers...)
}

// Create a new MuxAPI node that contains the Handlers.
func (frame *Framework) API(methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI("", methodset, pattern, handlers...)
}

// Create a new MuxAPI node that contains the named Handlers.
func NamedAPI(name string, methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedAPI(name, methodset, pattern, handlers...)
}

// Create a new MuxAPI node that contains the named Handlers.
func (frame *Framework) NamedAPI(name string, methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return newMuxAPI(frame, name, methodset, pattern, handlers...)
}

// GET is a shortcut for defaultFramework.GET(pattern, handlers...)
func GET(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.GET(pattern, handlers...)
}

// GET is a shortcut for frame.API("GET", pattern, handlers...)
func (frame *Framework) GET(pattern string, handlers ...Handler) *MuxAPI {
	return frame.API("GET", pattern, handlers...)
}

// HEAD is a shortcut for defaultFramework.HEAD(pattern, handlers...)
func HEAD(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.HEAD(pattern, handlers...)
}

// HEAD is a shortcut for frame.API("HEAD", pattern, handlers...)
func (frame *Framework) HEAD(pattern string, handlers ...Handler) *MuxAPI {
	return frame.API("HEAD", pattern, handlers...)
}

// OPTIONS is a shortcut for defaultFramework.OPTIONS(pattern, handlers...)
func OPTIONS(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.OPTIONS(pattern, handlers...)
}

// OPTIONS is a shortcut for frame.API("OPTIONS", pattern, handlers...)
func (frame *Framework) OPTIONS(pattern string, handlers ...Handler) *MuxAPI {
	return frame.API("OPTIONS", pattern, handlers...)
}

// POST is a shortcut for defaultFramework.POST(pattern, handlers...)
func POST(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.POST(pattern, handlers...)
}

// POST is a shortcut for frame.API("POST", pattern, handlers...)
func (frame *Framework) POST(pattern string, handlers ...Handler) *MuxAPI {
	return frame.API("POST", pattern, handlers...)
}

// PUT is a shortcut for defaultFramework.PUT(pattern, handlers...)
func PUT(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.PUT(pattern, handlers...)
}

// PUT is a shortcut for frame.API("PUT", pattern, handlers...)
func (frame *Framework) PUT(pattern string, handlers ...Handler) *MuxAPI {
	return frame.API("PUT", pattern, handlers...)
}

// PATCH is a shortcut for defaultFramework.PATCH(pattern, handlers...)
func PATCH(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.PATCH(pattern, handlers...)
}

// PATCH is a shortcut for frame.API("PATCH", pattern, handlers...)
func (frame *Framework) PATCH(pattern string, handlers ...Handler) *MuxAPI {
	return frame.API("PATCH", pattern, handlers...)
}

// DELETE is a shortcut for defaultFramework.DELETE(pattern, handlers...)
func DELETE(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.DELETE(pattern, handlers...)
}

// DELETE is a shortcut for frame.API("DELETE", pattern, handlers...)
func (frame *Framework) DELETE(pattern string, handlers ...Handler) *MuxAPI {
	return frame.API("DELETE", pattern, handlers...)
}

// NamedGET is a shortcut for defaultFramework.NamedGET(name, pattern, handlers...)
func NamedGET(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedGET(name, pattern, handlers...)
}

// NamedGET is a shortcut for frame.NamedAPI(name, "GET", pattern, handlers...)
func (frame *Framework) NamedGET(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI(name, "GET", pattern, handlers...)
}

// NamedHEAD is a shortcut for defaultFramework.NamedHEAD(name, pattern, handlers...)
func NamedHEAD(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedHEAD(name, pattern, handlers...)
}

// NamedHEAD is a shortcut for frame.NamedAPI(name, "HEAD", pattern, handlers...)
func (frame *Framework) NamedHEAD(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI(name, "HEAD", pattern, handlers...)
}

// NamedOPTIONS is a shortcut for defaultFramework.NamedOPTIONS(name, pattern, handlers...)
func NamedOPTIONS(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedOPTIONS(name, pattern, handlers...)
}

// NamedOPTIONS is a shortcut for frame.NamedAPI(name, "OPTIONS", pattern, handlers...)
func (frame *Framework) NamedOPTIONS(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI(name, "OPTIONS", pattern, handlers...)
}

// NamedPOST is a shortcut for defaultFramework.NamedPOST(name, pattern, handlers...)
func NamedPOST(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedPOST(name, pattern, handlers...)
}

// NamedPOST is a shortcut for frame.NamedAPI(name, "POST", pattern, handlers...)
func (frame *Framework) NamedPOST(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI(name, "POST", pattern, handlers...)
}

// NamedPUT is a shortcut for defaultFramework.NamedPUT(name, pattern, handlers...)
func NamedPUT(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedPUT(name, pattern, handlers...)
}

// NamedPUT is a shortcut for frame.NamedAPI(name, "PUT", pattern, handlers...)
func (frame *Framework) NamedPUT(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI(name, "PUT", pattern, handlers...)
}

// NamedPATCH is a shortcut for defaultFramework.NamedPATCH(name, pattern, handlers...)
func NamedPATCH(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedPATCH(name, pattern, handlers...)
}

// NamedPATCH is a shortcut for frame.NamedAPI(name, "PATCH", pattern, handlers...)
func (frame *Framework) NamedPATCH(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI(name, "PATCH", pattern, handlers...)
}

// NamedDELETE is a shortcut for defaultFramework.NamedDELETE(name, pattern, handlers...)
func NamedDELETE(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedDELETE(name, pattern, handlers...)
}

// NamedDELETE is a shortcut for frame.NamedAPI(name, "DELETE", pattern, handlers...)
func (frame *Framework) NamedDELETE(name string, pattern string, handlers ...Handler) *MuxAPI {
	return frame.NamedAPI(name, "DELETE", pattern, handlers...)
}

// StaticFS serves files from the given file system fs.
// The pattern must end with "/*filepath", files are then served from the local
// pattern /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     frame.StaticFS("/src/*filepath", http.Dir("/var/www"))
func StaticFS(pattern string, fs http.FileSystem) *MuxAPI {
	return defaultFramework.StaticFS(pattern, fs)
}

// StaticFS serves files from the given file system fs.
// The pattern must end with "/*filepath", files are then served from the local
// pattern /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     frame.StaticFS("/src/*filepath", http.Dir("/var/www"))
func (frame *Framework) StaticFS(pattern string, fs http.FileSystem) *MuxAPI {
	return (&MuxAPI{frame: frame}).StaticFS(pattern, fs)
}

func NamedStaticFS(name, pattern string, fs http.FileSystem) *MuxAPI {
	return defaultFramework.NamedStaticFS(name, pattern, fs)
}

func (frame *Framework) NamedStaticFS(name, pattern string, fs http.FileSystem) *MuxAPI {
	return (&MuxAPI{frame: frame}).NamedStaticFS(name, pattern, fs)
}

// An easy way to register a static route
func Static(pattern string, root string) *MuxAPI {
	return defaultFramework.Static(pattern, root)
}

// An easy way to register a static route
func (frame *Framework) Static(pattern string, root string) *MuxAPI {
	return (&MuxAPI{frame: frame}).Static(pattern, root)
}

// An easy way to register a static named route
func NamedStatic(name, pattern string, root string) *MuxAPI {
	return defaultFramework.NamedStatic(name, pattern, root)
}

// An easy way to register a static named route
func (frame *Framework) NamedStatic(name, pattern string, root string) *MuxAPI {
	return (&MuxAPI{frame: frame}).NamedStatic(name, pattern, root)
}

// Get an ordered list of nodes used to register router.
func MuxAPIsForRouter() []*MuxAPI {
	return defaultFramework.MuxAPIsForRouter()
}

// Get an ordered list of nodes used to register router.
func (frame *Framework) MuxAPIsForRouter() []*MuxAPI {
	if frame.muxesForRouter == nil {
		// comb mux.handlers, mux.paramInfos, mux.returns and mux.path,.
		frame.mux.comb()

		frame.muxesForRouter = frame.mux.HandlerProgeny()
	}
	return frame.muxesForRouter
}

// makeHandle makes an *apiware.ParamsAPI implements the Handle interface.
func (frame *Framework) makeHandle(handlerChain HandlerChain) Handle {
	ctxPool := sync.Pool{
		New: func() interface{} {
			return newContext(frame, handlerChain)
		},
	}
	return func(w http.ResponseWriter, r *http.Request, pathParams Params) {
		ctx := ctxPool.Get().(*Context)
		ctx.reset(w, r, pathParams)
		defer func() {
			ctxPool.Put(ctx)
		}()
		ctx.start()
	}
}

// Create the handle to be called by the router
func (frame *Framework) makeErrorHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		frame.errorFunc(newEmptyContext(frame, w, r), http.StatusText(status), status)
	})
}

// Create the handle to be called by the router
func (frame *Framework) makePanicHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	s := []byte("/src/runtime/panic.go")
	line := []byte("\n")
	return func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		stack := make([]byte, 4<<10) //4KB
		length := runtime.Stack(stack, true)
		start := bytes.Index(stack, s)
		stack = stack[start:length]
		start = bytes.Index(stack, line) + 1
		errStr := fmt.Sprintf("%v\n\n[STACK]\n%s", rcv, stack[start:])
		frame.errorFunc(newEmptyContext(frame, w, r), errStr, http.StatusInternalServerError)
	}
}

func (frame *Framework) registerSystemStatic() {
	frame.mux.NamedStatic("Directory for uploading files", "/upload/", UPLOAD_DIR)
	frame.mux.NamedStatic("Directory for public static files", "/static/", STATIC_DIR)
}

func (frame *Framework) registerSession() {
	if !frame.config.Session.Enable {
		return
	}
	conf := &session.ManagerConfig{
		CookieName:              frame.config.Session.Name,
		EnableSetCookie:         frame.config.Session.AutoSetCookie,
		Gclifetime:              frame.config.Session.GCMaxLifetime,
		Secure:                  frame.config.NetType == "tls" || frame.config.NetType == "letsencrypt",
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

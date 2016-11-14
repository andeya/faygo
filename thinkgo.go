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
	"github.com/henrylee2cn/thinkgo/logging"
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
	*MuxAPI        // root muxAPI node
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
	syslog         *logging.Logger // for framework
	bizlog         *logging.Logger // for user bissness
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
	frame.initSysLogger()
	frame.initBizLogger()
	frame.MuxAPI = newMuxAPI(frame, "root", "", "/")
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
		defaultFramework.initSysLogger()
		defaultFramework.initBizLogger()
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
			frame.presetSystemMuxes()
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

				frame.syslog.Criticalf("%7s | %-30s", method, node.path)

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

// Default returns the default framework.
func Default() *Framework {
	return defaultFramework
}

// The log used by the user bissness
func Log() *logging.Logger {
	return defaultFramework.Log()
}

// The log used by the user bissness
func (frame *Framework) Log() *logging.Logger {
	return frame.bizlog
}

// Get an ordered list of nodes used to register router.
func MuxAPIsForRouter() []*MuxAPI {
	return defaultFramework.MuxAPIsForRouter()
}

// Get an ordered list of nodes used to register router.
func (frame *Framework) MuxAPIsForRouter() []*MuxAPI {
	if frame.muxesForRouter == nil {
		// comb mux.handlers, mux.paramInfos, mux.returns and mux.path,.
		frame.MuxAPI.comb()

		frame.muxesForRouter = frame.MuxAPI.HandlerProgeny()
	}
	return frame.muxesForRouter
}

/**
 * -----------------------Register the middleware for the root node---------------------------
 */

// Insert the middlewares at the left end of the node's handler chain.
func Use(handlers ...HandlerWithoutPath) *MuxAPI {
	return defaultFramework.Use(handlers...)
}

/**
 * -----------------------------Add subordinate muxAPI nodes----------------------------------
 * ------------------------Used to register router in chain style-----------------------------
 */

// Group adds a subordinate subgroup node to the current muxAPI grouping node.
func Group(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.Group(pattern, handlers...)
}

// NamedGroup adds a subordinate subgroup node with the name to the current muxAPI grouping node.
func NamedGroup(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedGroup(name, pattern, handlers...)
}

// API adds a subordinate node to the current muxAPI grouping node.
func API(methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.API(methodset, pattern, handlers...)
}

// NamedAPI adds a subordinate node with the name to the current muxAPI grouping node.
func NamedAPI(name string, methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedAPI(name, methodset, pattern, handlers...)
}

// GET is a shortcut for defaultFramework.GET(pattern, handlers...)
func GET(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.GET(pattern, handlers...)
}

// HEAD is a shortcut for defaultFramework.HEAD(pattern, handlers...)
func HEAD(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.HEAD(pattern, handlers...)
}

// OPTIONS is a shortcut for defaultFramework.OPTIONS(pattern, handlers...)
func OPTIONS(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.OPTIONS(pattern, handlers...)
}

// POST is a shortcut for defaultFramework.POST(pattern, handlers...)
func POST(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.POST(pattern, handlers...)
}

// PUT is a shortcut for defaultFramework.PUT(pattern, handlers...)
func PUT(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.PUT(pattern, handlers...)
}

// PATCH is a shortcut for defaultFramework.PATCH(pattern, handlers...)
func PATCH(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.PATCH(pattern, handlers...)
}

// DELETE is a shortcut for defaultFramework.DELETE(pattern, handlers...)
func DELETE(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.DELETE(pattern, handlers...)
}

// NamedGET is a shortcut for defaultFramework.NamedGET(name, pattern, handlers...)
func NamedGET(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedGET(name, pattern, handlers...)
}

// NamedHEAD is a shortcut for defaultFramework.NamedHEAD(name, pattern, handlers...)
func NamedHEAD(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedHEAD(name, pattern, handlers...)
}

// NamedOPTIONS is a shortcut for defaultFramework.NamedOPTIONS(name, pattern, handlers...)
func NamedOPTIONS(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedOPTIONS(name, pattern, handlers...)
}

// NamedPOST is a shortcut for defaultFramework.NamedPOST(name, pattern, handlers...)
func NamedPOST(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedPOST(name, pattern, handlers...)
}

// NamedPUT is a shortcut for defaultFramework.NamedPUT(name, pattern, handlers...)
func NamedPUT(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedPUT(name, pattern, handlers...)
}

// NamedPATCH is a shortcut for defaultFramework.NamedPATCH(name, pattern, handlers...)
func NamedPATCH(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedPATCH(name, pattern, handlers...)
}

// NamedDELETE is a shortcut for defaultFramework.NamedDELETE(name, pattern, handlers...)
func NamedDELETE(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NamedDELETE(name, pattern, handlers...)
}

// Static is a shortcut for defaultFramework.Static(pattern, root)
func Static(pattern string, root string) *MuxAPI {
	return defaultFramework.Static(pattern, root)
}

// NamedStatic is a shortcut for defaultFramework.NamedStatic(name, pattern, root)
func NamedStatic(name, pattern string, root string) *MuxAPI {
	return defaultFramework.NamedStatic(name, pattern, root)
}

// StaticFS is a shortcut for defaultFramework.StaticFS(pattern, fs)
func StaticFS(pattern string, fs http.FileSystem) *MuxAPI {
	return defaultFramework.StaticFS(pattern, fs)
}

// NamedStaticFS is a shortcut for defaultFramework.NamedStaticFS(name, pattern, fs)
func NamedStaticFS(name, pattern string, fs http.FileSystem) *MuxAPI {
	return defaultFramework.NamedStaticFS(name, pattern, fs)
}

/**
 * -----------------------------Create isolated muxAPI nodes----------------------------------
 * -------------------------Used to register router in tree style-----------------------------
 */

// Append middlewares of function type to root muxAPI.
// Used to register router in tree style.
func Route(children ...*MuxAPI) *MuxAPI {
	return defaultFramework.Route(children...)
}

// Append middlewares of function type to root muxAPI.
// Used to register router in tree style.
func (frame *Framework) Route(children ...*MuxAPI) *MuxAPI {
	frame.MuxAPI.children = append(frame.MuxAPI.children, children...)
	for _, child := range children {
		child.parent = frame.MuxAPI
	}
	return frame.MuxAPI
}

// NewGroup create an isolated grouping muxAPI node.
func NewGroup(pattern string, children ...*MuxAPI) *MuxAPI {
	return defaultFramework.NewGroup(pattern, children...)
}

// NewAPI creates an isolated muxAPI node.
func NewAPI(methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewAPI(methodset, pattern, handlers...)
}

// NewNamedGroup creates an isolated grouping muxAPI node with the name.
func NewNamedGroup(name string, pattern string, children ...*MuxAPI) *MuxAPI {
	return defaultFramework.NewNamedGroup(name, pattern, children...)
}

// NewNamedAPI creates an isolated muxAPI node with the name.
func NewNamedAPI(name string, methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedAPI(name, methodset, pattern, handlers...)
}

// NewGET is a shortcut for defaultFramework.NewGET(name,pattern, handlers ...)
func NewGET(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewGET(pattern, handlers...)
}

// NewHEAD is a shortcut for defaultFramework.NewHEAD(name,pattern, handlers ...)
func NewHEAD(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewHEAD(pattern, handlers...)
}

// NewOPTIONS is a shortcut for defaultFramework.NewOPTIONS(name,pattern, handlers ...)
func NewOPTIONS(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewOPTIONS(pattern, handlers...)
}

// NewPOST is a shortcut for defaultFramework.NewPOST(name,pattern, handlers ...)
func NewPOST(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewPOST(pattern, handlers...)
}

// NewPUT is a shortcut for defaultFramework.NewPUT(name,pattern, handlers ...)
func NewPUT(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewPUT(pattern, handlers...)
}

// NewPATCH is a shortcut for defaultFramework.NewPATCH(name,pattern, handlers ...)
func NewPATCH(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewPATCH(pattern, handlers...)
}

// NewDELETE is a shortcut for defaultFramework.NewDELETE(name,pattern, handlers ...)
func NewDELETE(pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewDELETE(pattern, handlers...)
}

// NewNamedGET is a shortcut for defaultFramework.NewNamedGET(name,pattern, handlers ...)
func NewNamedGET(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedGET(name, pattern, handlers...)
}

// NewNamedHEAD is a shortcut for defaultFramework.NewNamedHEAD(name,pattern, handlers ...)
func NewNamedHEAD(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedHEAD(name, pattern, handlers...)
}

// NewNamedOPTIONS is a shortcut for defaultFramework.NewNamedOPTIONS(name,pattern, handlers ...)
func NewNamedOPTIONS(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedOPTIONS(name, pattern, handlers...)
}

// NewNamedPOST is a shortcut for defaultFramework.NewNamedPOST(name,pattern, handlers ...)
func NewNamedPOST(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedPOST(name, pattern, handlers...)
}

// NewNamedPUT is a shortcut for defaultFramework.NewNamedPUT(name,pattern, handlers ...)
func NewNamedPUT(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedPUT(name, pattern, handlers...)
}

// NewNamedPATCH is a shortcut for defaultFramework.NewNamedPATCH(name,pattern, handlers ...)
func NewNamedPATCH(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedPATCH(name, pattern, handlers...)
}

// NewNamedDELETE is a shortcut for defaultFramework.NewNamedDELETE(name,pattern, handlers ...)
func NewNamedDELETE(name string, pattern string, handlers ...Handler) *MuxAPI {
	return defaultFramework.NewNamedDELETE(name, pattern, handlers...)
}

// NewNamedStatic creates an isolated static muxAPI node.
func NewStatic(pattern string, root string) *MuxAPI {
	return defaultFramework.NewStatic(pattern, root)
}

// NewNamedStatic creates an isolated static muxAPI node with the name.
func NewNamedStatic(name, pattern string, root string) *MuxAPI {
	return defaultFramework.NewNamedStatic(name, pattern, root)
}

// NewNamedStatic creates an isolated static muxAPI node.
func NewStaticFS(pattern string, fs http.FileSystem) *MuxAPI {
	return defaultFramework.NewStaticFS(pattern, fs)
}

// NewNamedStatic creates an isolated static muxAPI node with the name.
func NewNamedStaticFS(name, pattern string, fs http.FileSystem) *MuxAPI {
	return defaultFramework.NewNamedStaticFS(name, pattern, fs)
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

// NewNamedStatic creates an isolated static muxAPI node.
func (frame *Framework) NewStatic(pattern string, root string) *MuxAPI {
	return frame.NewNamedStatic("", pattern, root)
}

// NewNamedStatic creates an isolated static muxAPI node with the name.
func (frame *Framework) NewNamedStatic(name, pattern string, root string) *MuxAPI {
	return (&MuxAPI{frame: frame}).NamedStatic(name, pattern, root)
}

// NewNamedStatic creates an isolated static muxAPI node.
func (frame *Framework) NewStaticFS(pattern string, fs http.FileSystem) *MuxAPI {
	return frame.NewNamedStaticFS("", pattern, fs)
}

// NewNamedStatic creates an isolated static muxAPI node with the name.
func (frame *Framework) NewNamedStaticFS(name, pattern string, fs http.FileSystem) *MuxAPI {
	return (&MuxAPI{frame: frame}).NamedStaticFS(name, pattern, fs)
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

func (frame *Framework) presetSystemMuxes() {
	frame.Use(AccessLogWare())
	frame.MuxAPI.NamedStatic("Directory for uploading files", "/upload/", UPLOAD_DIR)
	frame.MuxAPI.NamedStatic("Directory for public static files", "/static/", STATIC_DIR)
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

// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless ruired by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package thinkgo

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"sync"

	"github.com/henrylee2cn/thinkgo/logging"
	"github.com/henrylee2cn/thinkgo/session"
	"github.com/henrylee2cn/thinkgo/utils"
)

// Headers
const (
	HeaderAccept                        = "Accept"
	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAuthorization                 = "Authorization"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLength                 = "Content-Length"
	HeaderContentType                   = "Content-Type"
	HeaderContentDescription            = "Content-Description"
	HeaderContentTransferEncoding       = "Content-Transfer-Encoding"
	HeaderCookie                        = "Cookie"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderLastModified                  = "Last-Modified"
	HeaderLocation                      = "Location"
	HeaderReferer                       = "Referer"
	HeaderUserAgent                     = "User-Agent"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderXForwardedProto               = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride           = "X-HTTP-Method-Override"
	HeaderXForwardedFor                 = "X-Forwarded-For"
	HeaderXRealIP                       = "X-Real-IP"
	HeaderXRequestedWith                = "X-Requested-With"
	HeaderServer                        = "Server"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderExpires                       = "Expires"
	HeaderCacheControl                  = "Cache-Control"
	HeaderPragma                        = "Pragma"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=utf-8"
	nosniff     = "nosniff" // Security
)

const (
	// stopExecutionposition used inside the Context,
	// is the number which shows us that the context's handlerChain manualy stop the execution
	stopExecutionposition = 125
)

type (
	// Map is just a conversion for a map[string]interface{}
	// should not be used inside Render when PongoEngine is used.
	Map map[string]interface{}

	// Context is resetting every time a ruest is coming to the server
	// it is not good practice to use this object in goroutines, for these cases use the .Clone()
	Context struct {
		R                  *http.Request // the *http.Request
		W                  *Response     // the *Response cooked by the http.ResponseWriter
		CruSession         session.Store
		limitedRequestBody []byte // the copy of requset body(Limited by maximum length)
		frame              *Framework
		handlerChain       HandlerChain                // keep track all registed handlers
		pathParams         Params                      // The parameter values on the URL path
		queryParams        url.Values                  // URL query string values
		data               map[interface{}]interface{} // Used to transfer variables between Handler-chains
		handlerChainLen    int8
		pos                int8 // pos is the position number of the Context, look .Next to understand
		enableGzip         bool // Note: Never reset!
		enableSession      bool // Note: Never reset!
		enableXSRF         bool // Note: Never reset!
		xsrfExpire         int
		_xsrfToken         string
		_xsrfTokenReset    bool
	}
)

// The log used by the user bissness
func (ctx *Context) Log() *logging.Logger {
	return ctx.frame.bizlog
}

// XSRFToken creates a xsrf token string and returns.
// If specifiedExpiration is empty, the value in the configuration is used.
func (ctx *Context) XSRFToken(specifiedExpiration ...int) string {
	if ctx._xsrfToken == "" {
		token, ok := ctx.SecureCookieParam(ctx.frame.config.XSRF.Key, "_xsrf")
		if !ok {
			ctx._xsrfTokenReset = true
			token = string(utils.RandomCreateBytes(32))
			if len(specifiedExpiration) > 0 && specifiedExpiration[0] > 0 {
				ctx.xsrfExpire = specifiedExpiration[0]
			} else if ctx.xsrfExpire == 0 {
				ctx.xsrfExpire = ctx.frame.config.XSRF.Expire
			}
		}
		ctx._xsrfToken = token
	}
	return ctx._xsrfToken
}

// XSRFFormHTML writes an input field contains xsrf token value.
func (ctx *Context) XSRFFormHTML() string {
	return `<input type="hidden" name="_xsrf" value="` +
		ctx.XSRFToken() + `" />`
}

// checkXSRFCookie checks xsrf token in this ruest is valid or not.
// the token can provided in ruest header "X-Xsrftoken" and "X-CsrfToken"
// or in form field value named as "_xsrf".
func (ctx *Context) checkXSRFCookie() bool {
	if !ctx.enableXSRF {
		return true
	}
	token := ctx.BizParam("_xsrf")
	if token == "" {
		token = ctx.R.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = ctx.R.Header.Get("X-Csrftoken")
	}
	if token == "" {
		ctx.Error(403, "'_xsrf' argument missing from POST")
		return false
	}
	if ctx._xsrfToken != token {
		ctx.Error(403, "XSRF cookie does not match POST argument")
		return false
	}
	return true
}

// StartSession starts session and load old session data info this controller.
func (ctx *Context) StartSession() (session.Store, error) {
	if !ctx.enableSession {
		return nil, errors.New("session function is disable.")
	}
	if ctx.CruSession != nil {
		return ctx.CruSession, nil
	}
	var err error
	ctx.CruSession, err = ctx.frame.sessionManager.SessionStart(ctx.W, ctx.R)
	return ctx.CruSession, err
}

// SetSession puts value into session.
func (ctx *Context) SetSession(key interface{}, value interface{}) {
	if ctx.CruSession == nil {
		if _, err := ctx.StartSession(); err != nil {
			return
		}
	}
	ctx.CruSession.Set(key, value)
}

// GetSession gets value from session.
func (ctx *Context) GetSession(key interface{}) interface{} {
	if ctx.CruSession == nil {
		if _, err := ctx.StartSession(); err != nil {
			return nil
		}
	}
	return ctx.CruSession.Get(key)
}

// DelSession removes value from session.
func (ctx *Context) DelSession(key interface{}) {
	if ctx.CruSession == nil {
		if _, err := ctx.StartSession(); err != nil {
			return
		}
	}
	ctx.CruSession.Delete(key)
}

// SessionRegenerateID regenerates session id for this session.
// the session data have no changes.
func (ctx *Context) SessionRegenerateID() {
	if ctx.CruSession == nil {
		if _, err := ctx.StartSession(); err != nil {
			return
		}
	}
	ctx.CruSession.SessionRelease(ctx.W)
	ctx.CruSession = ctx.frame.sessionManager.SessionRegenerateID(ctx.W, ctx.R)
}

// DestroySession cleans session data and session cookie.
func (ctx *Context) DestroySession() {
	if ctx.CruSession == nil {
		if _, err := ctx.StartSession(); err != nil {
			return
		}
	}
	ctx.CruSession.Flush()
	ctx.CruSession = nil
	ctx.frame.sessionManager.SessionDestroy(ctx.W, ctx.R)
}

// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
//
// The provided status code should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.
func (ctx *Context) Redirect(status int, urlStr string) error {
	if status < http.StatusMultipleChoices || status > http.StatusPermanentRedirect {
		return fmt.Errorf("The provided status code should be in the 3xx range and is usually 301, 302 or 303, yours: %d", status)
	}
	http.Redirect(ctx.W, ctx.R, urlStr, status)
	return nil
}

var proxyList = &struct {
	m map[string]*httputil.ReverseProxy
	sync.RWMutex
}{
	m: map[string]*httputil.ReverseProxy{},
}

// ReverseProxy routes URLs to the scheme, host, and base path provided in targetUrlBase.
// If pathAppend is "true" and the targetUrlBase's path is "/base" and the incoming ruest was for "/dir",
// the target ruest will be for /base/dir.
func (ctx *Context) ReverseProxy(targetUrlBase string, pathAppend bool) error {
	proxyList.RLock()
	var rp = proxyList.m[targetUrlBase]
	proxyList.RUnlock()
	if rp == nil {
		proxyList.Lock()
		defer proxyList.Unlock()
		rp = proxyList.m[targetUrlBase]
		if rp == nil {
			target, err := url.Parse(targetUrlBase)
			if err != nil {
				return err
			}
			targetQuery := target.RawQuery
			rp = &httputil.ReverseProxy{
				Director: func(r *http.Request) {
					r.Host = target.Host
					r.URL.Scheme = target.Scheme
					r.URL.Host = target.Host
					r.URL.Path = path.Join(target.Path, r.URL.Path)
					if targetQuery == "" || r.URL.RawQuery == "" {
						r.URL.RawQuery = targetQuery + r.URL.RawQuery
					} else {
						r.URL.RawQuery = targetQuery + "&" + r.URL.RawQuery
					}
				},
			}
			proxyList.m[targetUrlBase] = rp
		}
	}

	if !pathAppend {
		ctx.R.URL.Path = ""
	}
	rp.ServeHTTP(ctx.W, ctx.R)
	return nil
}

// Create the context for the router handle
func newEmptyContext(
	frame *Framework,
	w http.ResponseWriter,
	r *http.Request,
) *Context {
	ctx := &Context{
		frame:      frame,
		R:          r,
		enableGzip: Global.config.Gzip.Enable,
	}
	ctx.W = newResponse(ctx, w)
	return ctx
}

// Create the context for the filter
func newFilterContext(
	frame *Framework,
) *Context {
	ctx := &Context{
		frame:           frame,
		handlerChain:    frame.filter,
		handlerChainLen: int8(len(frame.filter)),
		pos:             0,
		data:            make(map[interface{}]interface{}),
	}
	ctx.W = newResponse(ctx, nil)
	return ctx
}

// Create the context for common handle
func newContext(
	frame *Framework,
	handlerChain HandlerChain,
) *Context {
	count := len(handlerChain)
	chain := make(HandlerChain, count)
	copy(chain, handlerChain)
	for i, h := range chain {
		if h2, ok := h.(*handlerStruct); ok {
			chain[i] = h2.new()
		}
	}
	ctx := &Context{
		frame:           frame,
		handlerChain:    chain,
		handlerChainLen: int8(count),
		pos:             0,
		enableGzip:      Global.config.Gzip.Enable,
		enableSession:   frame.config.Session.Enable,
		enableXSRF:      frame.config.XSRF.Enable,
		data:            make(map[interface{}]interface{}),
	}
	ctx.W = newResponse(ctx, nil)
	return ctx
}

// Called before the start
func (ctx *Context) prepare() bool {
	var pass = true
	//if XSRF is Enable then check cookie where there has any cookie in the  request's cookie _csrf
	if ctx.enableXSRF {
		ctx.XSRFToken()
		switch ctx.R.Method {
		case "POST", "DELETE", "PUT":
			pass = ctx.checkXSRFCookie()
		default:
			switch ctx.BizParam("_method") {
			case "POST", "DELETE", "PUT":
				pass = ctx.checkXSRFCookie()
			}
		}
	}
	return pass
}

// reset the cursor
func (ctx *Context) posReset() {
	ctx.pos = -1
}

// do calls the first handler only, it's like Next with negative pos, used only on Router&MemoryRouter
func (ctx *Context) do() {
	if !ctx.prepare() {
		return
	}
	ctx.posReset()
	ctx.Next()
}

// Next calls all the next handler from the middleware stack, it used inside a middleware
func (ctx *Context) Next() {
	//set position to the next
	ctx.pos++
	//run the next
	if ctx.pos < ctx.handlerChainLen {
		switch h := ctx.handlerChain[ctx.pos].(type) {
		case *handlerStruct:
			err := h.bind(ctx.R, ctx.pathParams)
			defer h.reset()
			if err != nil {
				Global.bindErrorFunc(ctx, err)
				ctx.Stop()
				return
			}
			err = h.Serve(ctx)
			if err != nil {
				Global.errorFunc(ctx, err.Error(), http.StatusInternalServerError)
				ctx.Stop()
				return
			}
		default:
			err := h.Serve(ctx)
			if err != nil {
				Global.errorFunc(ctx, err.Error(), http.StatusInternalServerError)
				ctx.Stop()
				return
			}
		}
		// If the next one exists, it is executed automatically.
		ctx.Next()
		return
	}
	ctx.pos--
}

func (ctx *Context) beforeWriteHeader() {
	if ctx._xsrfTokenReset {
		ctx.SetSecureCookie(ctx.frame.config.XSRF.Key, "_xsrf", ctx._xsrfToken, ctx.xsrfExpire)
	}
	if ctx.enableSession {
		if ctx.CruSession != nil {
			ctx.CruSession.SessionRelease(ctx.W)
			ctx.CruSession = nil
		}
	}
}

// Stop just sets the .pos to 125 in order to  not move to the next handlers(if any)
func (ctx *Context) Stop() {
	ctx.pos = stopExecutionposition
}

func (ctx *Context) isStop() bool {
	return ctx.pos >= ctx.handlerChainLen
}

func (ctx *Context) isActiveStop() bool {
	return ctx.pos == stopExecutionposition
}

// reset ctx.
// Note: Never reset `ctx.frame`, `ctx.W`, `ctx.enableGzip`, `ctx.enableSession` and `ctx.enableXSRF`!
func (ctx *Context) reset(w http.ResponseWriter, r *http.Request, pathParams Params, data map[interface{}]interface{}) {
	ctx.limitedRequestBody = nil
	ctx.data = data
	ctx.queryParams = nil
	ctx._xsrfToken = ""
	ctx._xsrfTokenReset = false
	ctx.pathParams = pathParams
	ctx.W.reset(w)
	ctx.R = r
}

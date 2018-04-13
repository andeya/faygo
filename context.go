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

package faygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"sync"

	"github.com/henrylee2cn/faygo/logging"
	"github.com/henrylee2cn/faygo/session"
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
	stopExecutionposition int16 = math.MaxInt16 - 1
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
		curSession         session.Store
		limitedRequestBody []byte // the copy of requset body(Limited by maximum length)
		frame              *Framework
		handlerChain       HandlerChain                // keep track all registed handlers
		pathParams         PathParams                  // The parameter values on the URL path
		queryParams        url.Values                  // URL query string values
		data               map[interface{}]interface{} // Used to transfer variables between Handler-chains
		handlerChainLen    int16
		pos                int16 // pos is the position number of the Context, look .Next to understand
		enableGzip         bool  // Note: Never reset!
		enableSession      bool  // Note: Never reset!
		enableXSRF         bool  // Note: Never reset!
		xsrfExpire         int
		_xsrfToken         string
		_xsrfTokenReset    bool
	}
)

// Log used by the user bissness
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
			token = RandomString(32)
			if len(specifiedExpiration) > 0 && specifiedExpiration[0] > 0 {
				ctx.xsrfExpire = specifiedExpiration[0]
			} else if ctx.xsrfExpire == 0 {
				ctx.xsrfExpire = ctx.frame.config.XSRF.ExpireSecond
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

// checkXSRFCookie checks xsrf token in this request is valid or not.
// the token can provided in request cookie "_xsrf",
// or in header "X-Xsrftoken" and "X-CsrfToken",
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
	// default cookie value
	if token == "" {
		token, _ = ctx.SecureCookieParam(ctx.frame.config.XSRF.Key, "_xsrf")
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

var errNotEnableSession = errors.New("before using the session, must set config `session::enable = true`...")

// startSession starts session and load old session data info this controller.
func (ctx *Context) startSession() (session.Store, error) {
	if ctx.curSession != nil {
		return ctx.curSession, nil
	}
	if !ctx.enableSession {
		return nil, errNotEnableSession
	}
	ctx.makeSureParseMultipartForm()
	var err error
	ctx.curSession, err = ctx.frame.sessionManager.SessionStart(ctx.W, ctx.R)
	return ctx.curSession, err
}

// getSessionStore return SessionStore.
func (ctx *Context) getSessionStore() (session.Store, error) {
	if ctx.curSession != nil {
		return ctx.curSession, nil
	}
	if !ctx.enableSession {
		return nil, errNotEnableSession
	}
	var err error
	ctx.curSession, err = ctx.frame.sessionManager.GetSessionStore(ctx.W, ctx.R)
	return ctx.curSession, err
}

// SetSession puts value into session.
func (ctx *Context) SetSession(key interface{}, value interface{}) {
	if _, err := ctx.startSession(); err != nil {
		ctx.Log().Warning(err.Error())
		return
	}
	ctx.curSession.Set(key, value)
}

// GetSession gets value from session.
func (ctx *Context) GetSession(key interface{}) interface{} {
	if _, err := ctx.getSessionStore(); err != nil {
		return nil
	}
	return ctx.curSession.Get(key)
}

// DelSession removes value from session.
func (ctx *Context) DelSession(key interface{}) {
	if _, err := ctx.getSessionStore(); err != nil {
		return
	}
	ctx.curSession.Delete(key)
}

// SessionRegenerateID regenerates session id for this session.
// the session data have no changes.
func (ctx *Context) SessionRegenerateID() {
	if _, err := ctx.getSessionStore(); err != nil {
		return
	}
	ctx.curSession.SessionRelease(ctx.W)
	ctx.curSession = ctx.frame.sessionManager.SessionRegenerateID(ctx.W, ctx.R)
}

// DestroySession cleans session data and session cookie.
func (ctx *Context) DestroySession() {
	if _, err := ctx.getSessionStore(); err != nil {
		return
	}
	ctx.curSession.Flush()
	ctx.curSession = nil
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

func (ctx *Context) doFilter() bool {
	if count := len(ctx.frame.filter); count > 0 {
		ctx.handlerChain = ctx.frame.filter
		ctx.handlerChainLen = int16(count)
		ctx.posReset()
		ctx.Next()
		if ctx.IsBreak() {
			if !ctx.W.Committed() {
				ctx.Error(http.StatusForbidden, http.StatusText(http.StatusForbidden))
			}
			return false
		}
	}
	return true
}

// doHandler calls the first handler only, it's like Next with negative pos, used only on Router&MemoryRouter
func (ctx *Context) doHandler(handlerChain HandlerChain, pathParams PathParams) {
	ctx.pathParams = pathParams
	ctx.handlerChain = handlerChain
	ctx.handlerChainLen = int16(len(handlerChain))
	ctx.posReset()
	if !ctx.prepare() {
		return
	}
	ctx.Next()
}

// Called before the start
func (ctx *Context) prepare() bool {
	var pass = true
	//if XSRF is Enable then check cookie where there has any cookie in the request's cookie _csrf
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

// Next calls all the next handler from the middleware stack, it used inside a middleware.
// Notes: Non-concurrent security.
func (ctx *Context) Next() {
	//set position to the next
	ctx.pos++
	//run the next
	if ctx.pos < ctx.handlerChainLen {
		if err := ctx.handlerChain[ctx.pos].Serve(ctx); err != nil {
			global.errorFunc(ctx, err.Error(), http.StatusInternalServerError)
			ctx.Stop()
			return
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
		if ctx.curSession != nil {
			ctx.curSession.SessionRelease(ctx.W)
			ctx.curSession = nil
		}
	}
}

// Stop just sets the .pos to 32766 in order to  not move to the next handlers(if any)
func (ctx *Context) Stop() {
	ctx.pos = stopExecutionposition
}

// Stopped returns whether the operation has stopped.
func (ctx *Context) Stopped() bool {
	return ctx.pos >= ctx.handlerChainLen
}

// IsBreak returns whether the operation is stopped halfway.
func (ctx *Context) IsBreak() bool {
	return ctx.pos == stopExecutionposition
}

func (ctx *Context) recordBody() []byte {
	if !ctx.frame.config.PrintBody {
		return nil
	}
	var b []byte
	formValues := ctx.FormParamAll()
	if len(formValues) > 0 || ctx.IsUpload() {
		v := multipart.Form{
			Value: formValues,
		}
		if ctx.R.MultipartForm != nil {
			v.File = ctx.R.MultipartForm.File
		}
		b, _ = json.Marshal(v)
	} else {
		b = ctx.LimitedBodyBytes()
	}
	if len(b) > 0 {
		bb := make([]byte, len(b)+2)
		bb[0] = '\n'
		copy(bb[1:], b)
		bb[len(bb)-1] = '\n'
		return bb
	}
	return b
}

func (frame *Framework) getContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := frame.contextPool.Get().(*Context)
	ctx.R = r
	ctx.W.reset(w)
	ctx.data = make(map[interface{}]interface{})
	if frame.config.PrintBody && !ctx.IsUpload() {
		ctx.LimitedBodyBytes()
	}
	return ctx
}

func (frame *Framework) putContext(ctx *Context) {
	if ctx.R.Body != nil {
		ctx.R.Body.Close()
	}
	ctx.R = nil
	ctx.W.writer = nil
	ctx.limitedRequestBody = nil
	ctx.data = nil
	ctx.queryParams = nil
	ctx._xsrfToken = ""
	ctx._xsrfTokenReset = false
	frame.contextPool.Put(ctx)
}

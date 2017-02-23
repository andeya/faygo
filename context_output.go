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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/henrylee2cn/faygo/acceptencoder"
)

// Size returns the current size, in bytes, of the response.
func (ctx *Context) Size() int64 {
	return ctx.W.Size()
}

// Committed returns whether the response has been submitted or not.
func (ctx *Context) Committed() bool {
	return ctx.W.committed
}

// Status returns the HTTP status code of the response.
func (ctx *Context) Status() int {
	return ctx.W.status
}

// IsCachable returns boolean of this request is cached.
// HTTP 304 means cached.
func (ctx *Context) IsCachable() bool {
	return ctx.W.status >= 200 && ctx.W.status < 300 || ctx.W.status == 304
}

// IsEmpty returns boolean of this request is empty.
// HTTP 201，204 and 304 means empty.
func (ctx *Context) IsEmpty() bool {
	return ctx.W.status == 201 || ctx.W.status == 204 || ctx.W.status == 304
}

// IsOk returns boolean of this request runs well.
// HTTP 200 means ok.
func (ctx *Context) IsOk() bool {
	return ctx.W.status == 200
}

// IsSuccessful returns boolean of this request runs successfully.
// HTTP 2xx means ok.
func (ctx *Context) IsSuccessful() bool {
	return ctx.W.status >= 200 && ctx.W.status < 300
}

// IsRedirect returns boolean of this request is redirection header.
// HTTP 301,302,307 means redirection.
func (ctx *Context) IsRedirect() bool {
	return ctx.W.status == 301 || ctx.W.status == 302 || ctx.W.status == 303 || ctx.W.status == 307
}

// IsForbidden returns boolean of this request is forbidden.
// HTTP 403 means forbidden.
func (ctx *Context) IsForbidden() bool {
	return ctx.W.status == 403
}

// IsNotFound returns boolean of this request is not found.
// HTTP 404 means forbidden.
func (ctx *Context) IsNotFound() bool {
	return ctx.W.status == 404
}

// IsClientError returns boolean of this request client sends error data.
// HTTP 4xx means forbidden.
func (ctx *Context) IsClientError() bool {
	return ctx.W.status >= 400 && ctx.W.status < 500
}

// IsServerError returns boolean of this server handler errors.
// HTTP 5xx means server internal error.
func (ctx *Context) IsServerError() bool {
	return ctx.W.status >= 500 && ctx.W.status < 600
}

// SetHeader sets response header item string via given key.
func (ctx *Context) SetHeader(key, val string) {
	ctx.W.Header().Set(key, val)
}

// SetCookie sets cookie value via given key.
// others are ordered as cookie's max age time, path, domain, secure and httponly.
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
	//fix cookie not work in IE
	if len(others) > 0 {
		var maxAge int64
		switch v := others[0].(type) {
		case int:
			maxAge = int64(v)
		case int32:
			maxAge = int64(v)
		case int64:
			maxAge = v
		}
		switch {
		case maxAge > 0:
			fmt.Fprintf(&b, "; Expires=%s; Max-Age=%d", time.Now().Add(time.Duration(maxAge)*time.Second).UTC().Format(time.RFC1123), maxAge)
		case maxAge < 0:
			fmt.Fprintf(&b, "; Max-Age=0")
		}
	}
	// the settings below
	// Path, Domain, Secure, HttpOnly
	// can use nil skip set

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", sanitizeValue(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}

	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(v))
		}
	}

	// default empty
	if len(others) > 3 {
		var secure bool
		switch v := others[3].(type) {
		case bool:
			secure = v
		default:
			if others[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&b, "; Secure")
		}
	}

	// default false. for session cookie default true
	httponly := false
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			// HttpOnly = true
			httponly = true
		}
	}

	if httponly {
		fmt.Fprintf(&b, "; HttpOnly")
	}

	ctx.W.Header().Add(HeaderSetCookie, b.String())
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}

// SetSecureCookie Set Secure cookie for response.
func (ctx *Context) SetSecureCookie(secret, name, value string, others ...interface{}) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.SetCookie(name, cookie, others...)
}

// NoContent sends a response with no body and a status code.
func (ctx *Context) NoContent(status int) {
	ctx.W.WriteHeader(status)
}

// Send error message and stop handler chain.
func (ctx *Context) Error(status int, errStr string) {
	global.errorFunc(ctx, errStr, status)
	ctx.Stop()
}

// Bytes writes the data bytes to the connection as part of an HTTP reply.
func (ctx *Context) Bytes(status int, contentType string, content []byte) error {
	if ctx.W.committed {
		ctx.W.multiCommitted()
		return nil
	}
	ctx.W.Header().Set(HeaderContentType, contentType)
	if ctx.enableGzip && len(ctx.W.Header()[HeaderContentEncoding]) == 0 {
		buf := &bytes.Buffer{}
		ok, encoding, _ := acceptencoder.WriteBody(acceptencoder.ParseEncoding(ctx.R), buf, content)
		if ok {
			ctx.W.Header().Set(HeaderContentEncoding, encoding)
			content = buf.Bytes()
		}
	}
	ctx.W.Header().Set(HeaderContentLength, strconv.Itoa(len(content)))
	ctx.W.WriteHeader(status)
	_, err := ctx.W.Write(content)
	return err
}

// String writes a string to the client, something like fmt.Fprintf
func (ctx *Context) String(status int, format string, s ...interface{}) error {
	if len(s) == 0 {
		return ctx.Bytes(status, MIMETextPlainCharsetUTF8, []byte(format))
	}
	return ctx.Bytes(status, MIMETextPlainCharsetUTF8, []byte(fmt.Sprintf(format, s...)))
}

// HTML sends an HTTP response with status code.
func (ctx *Context) HTML(status int, html string) error {
	x := (*[2]uintptr)(unsafe.Pointer(&html))
	h := [3]uintptr{x[0], x[1], x[1]}
	return ctx.Bytes(status, MIMETextHTMLCharsetUTF8, *(*[]byte)(unsafe.Pointer(&h)))
}

// JSON sends a JSON response with status code.
func (ctx *Context) JSON(status int, data interface{}, isIndent ...bool) error {
	var (
		b   []byte
		err error
	)
	if len(isIndent) > 0 && isIndent[0] {
		b, err = json.MarshalIndent(data, "", "  ")
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	return ctx.JSONBlob(status, b)
}

// JSONBlob sends a JSON blob response with status code.
func (ctx *Context) JSONBlob(status int, b []byte) error {
	return ctx.Bytes(status, MIMEApplicationJSONCharsetUTF8, b)
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (ctx *Context) JSONP(status int, callback string, data interface{}, isIndent ...bool) error {
	var (
		b   []byte
		err error
	)
	if len(isIndent) > 0 && isIndent[0] {
		b, err = json.MarshalIndent(data, "", "  ")
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	callback = template.JSEscapeString(callback)
	callbackContent := bytes.NewBufferString(" if(window." + callback + ")" + callback)
	callbackContent.WriteString("(")
	callbackContent.Write(b)
	callbackContent.WriteString(");\r\n")
	return ctx.Bytes(status, MIMEApplicationJavaScriptCharsetUTF8, callbackContent.Bytes())
}

// JSONMsg sends a JSON with JSONMsg format.
func (ctx *Context) JSONMsg(status int, msgcode int, info interface{}, isIndent ...bool) error {
	var (
		b    []byte
		err  error
		data = JSONMsg{
			Code: msgcode,
			Info: info,
		}
	)
	if len(isIndent) > 0 && isIndent[0] {
		b, err = json.MarshalIndent(data, "", "  ")
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}
	return ctx.JSONBlob(status, b)
}

// XML sends an XML response with status code.
func (ctx *Context) XML(status int, data interface{}, isIndent ...bool) error {
	var (
		b   []byte
		err error
	)
	if len(isIndent) > 0 && isIndent[0] {
		b, err = xml.MarshalIndent(data, "", "  ")
	} else {
		b, err = xml.Marshal(data)
	}
	if err != nil {
		return err
	}
	return ctx.XMLBlob(status, b)
}

// XMLBlob sends a XML blob response with status code.
func (ctx *Context) XMLBlob(status int, b []byte) error {
	content := bytes.NewBufferString(xml.Header)
	content.Write(b)
	return ctx.Bytes(status, MIMEApplicationXMLCharsetUTF8, content.Bytes())
}

// JSONOrXML serve Xml OR Json, depending on the value of the Accept header
func (ctx *Context) JSONOrXML(status int, data interface{}, isIndent ...bool) error {
	if ctx.AcceptJSON() || !ctx.AcceptXML() {
		return ctx.JSON(status, data, isIndent...)
	}
	return ctx.XML(status, data, isIndent...)
}

// File forces response for download file.
// it prepares the download response header automatically.
func (ctx *Context) File(file string, filename ...string) {
	ctx.W.Header().Set(HeaderContentDescription, "File Transfer")
	ctx.W.Header().Set(HeaderContentType, MIMEOctetStream)
	if len(filename) > 0 && filename[0] != "" {
		ctx.W.Header().Set(HeaderContentDisposition, "attachment; filename="+filename[0])
	} else {
		ctx.W.Header().Set(HeaderContentDisposition, "attachment; filename="+filepath.Base(file))
	}
	ctx.W.Header().Set(HeaderContentTransferEncoding, "binary")
	ctx.W.Header().Set(HeaderExpires, "0")
	ctx.W.Header().Set(HeaderCacheControl, "must-revalidate")
	ctx.W.Header().Set(HeaderPragma, "public")
	global.fsManager.ServeFile(ctx, file)
}

// Render renders a template with data and sends a text/html response with status code.
func (ctx *Context) Render(status int, name string, data Map) error {
	b, err := global.render.Render(name, data)
	if err != nil {
		return err
	}
	return ctx.Bytes(status, MIMETextHTMLCharsetUTF8, b)
}

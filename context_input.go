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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/henrylee2cn/thinkgo/utils"
)

// Regexes for checking the accept headers
// TODO make sure these are correct
var (
	acceptsHTMLRegex = regexp.MustCompile(`(text/html|application/xhtml\+xml)(?:,|$)`)
	acceptsXMLRegex  = regexp.MustCompile(`(application/xml|text/xml)(?:,|$)`)
	acceptsJSONRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
)

// Protocol returns request protocol name, such as HTTP/1.1 .
func (ctx *Context) Protocol() string {
	return ctx.R.Proto
}

// URI returns full request url with query string, fragment.
func (ctx *Context) URI() string {
	return ctx.R.RequestURI
}

// Path returns request url path (without query string, fragment).
func (ctx *Context) Path() string {
	return ctx.R.URL.Path
}

// Scheme returns request scheme as "http" or "https".
func (ctx *Context) Scheme() string {
	if scheme := ctx.HeaderParam(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if ctx.R.URL.Scheme != "" {
		return ctx.R.URL.Scheme
	}
	if ctx.R.TLS == nil {
		return "http"
	}
	return "https"
}

// Site returns base site url as scheme://host type.
func (ctx *Context) Site() string {
	return ctx.Scheme() + "://" + ctx.Host()
}

// HostWithPort returns a host:port string for this request,
// such as "example.com" or "example.com:8080".
func (ctx *Context) HostWithPort() string {
	return ctx.R.Host
}

// Host returns host name.
// `host` is `subDomain.domain`.
// if no host info in request, return localhost.
func (ctx *Context) Host() string {
	if ctx.R.Host != "" {
		hostParts := strings.Split(ctx.R.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return ctx.R.Host
	}
	return "localhost"
}

// Domain returns domain name.
// `host` is `subDomain.domain`.
// if aa.bb.domain.com, returns aa.bb .
// if no host info in request, return localhost.
func (ctx *Context) Domain() string {
	parts := strings.Split(ctx.Host(), ".")
	if len(parts) >= 3 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return "localhost"
}

// SubDomain returns sub domain string.
// `host` is `subDomain.domain`.
// if aa.bb.domain.com, returns aa.bb .
func (ctx *Context) SubDomain() string {
	parts := strings.Split(ctx.Host(), ".")
	if len(parts) >= 3 {
		return strings.Join(parts[:len(parts)-2], ".")
	}
	return ""
}

// Port returns host port for this request.
// when error or empty, return 80.
func (ctx *Context) Port() int {
	parts := strings.Split(ctx.R.Host, ":")
	if len(parts) == 2 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
}

// IP gets just the ip from the most direct one client.
func (ctx *Context) IP() string {
	var ip = strings.Split(ctx.R.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

// RealIP returns request client ip.
// if in proxy, return first proxy id.
// if error, return 127.0.0.1.
func (ctx *Context) RealIP() string {
	var ip = ctx.R.Header.Get(HeaderXRealIP)
	if len(ip) > 0 {
		return ip
	}
	ips := ctx.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		return rip[0]
	}
	return ctx.IP()
}

// Proxy returns proxy client ips slice.
func (ctx *Context) Proxy() []string {
	if ips := ctx.HeaderParam(HeaderXForwardedFor); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

// Referer returns http referer header.
func (ctx *Context) Referer() string {
	return ctx.HeaderParam(HeaderReferer)
}

// Method returns http request method.
func (ctx *Context) Method() string {
	return ctx.R.Method
}

// Is returns boolean of this request is on given method, such as Is("POST").
func (ctx *Context) Is(method string) bool {
	return ctx.Method() == method
}

// IsGet Is this a GET method request?
func (ctx *Context) IsGet() bool {
	return ctx.Is("GET")
}

// IsPost Is this a POST method request?
func (ctx *Context) IsPost() bool {
	return ctx.Is("POST")
}

// IsHead Is this a Head method request?
func (ctx *Context) IsHead() bool {
	return ctx.Is("HEAD")
}

// IsOptions Is this a OPTIONS method request?
func (ctx *Context) IsOptions() bool {
	return ctx.Is("OPTIONS")
}

// IsPut Is this a PUT method request?
func (ctx *Context) IsPut() bool {
	return ctx.Is("PUT")
}

// IsDelete Is this a DELETE method request?
func (ctx *Context) IsDelete() bool {
	return ctx.Is("DELETE")
}

// IsPatch Is this a PATCH method request?
func (ctx *Context) IsPatch() bool {
	return ctx.Is("PATCH")
}

// IsAjax returns boolean of this request is generated by ajax.
func (ctx *Context) IsAjax() bool {
	return ctx.HeaderParam(HeaderXRequestedWith) == "XMLHttpRequest"
}

// IsSecure returns boolean of this request is in https.
func (ctx *Context) IsSecure() bool {
	return ctx.Scheme() == "https"
}

// IsWebsocket returns boolean of this request is in webSocket.
func (ctx *Context) IsWebsocket() bool {
	return ctx.HeaderParam(HeaderUpgrade) == "websocket"
}

// IsUpload returns boolean of whether file uploads in this request or not..
func (ctx *Context) IsUpload() bool {
	return strings.Contains(ctx.HeaderParam(HeaderContentType), MIMEMultipartForm)
}

// AcceptHTML Checks if request accepts html response
func (ctx *Context) AcceptHTML() bool {
	return acceptsHTMLRegex.MatchString(ctx.HeaderParam(HeaderAccept))
}

// AcceptXML Checks if request accepts xml response
func (ctx *Context) AcceptXML() bool {
	return acceptsXMLRegex.MatchString(ctx.HeaderParam(HeaderAccept))
}

// AcceptJSON Checks if request accepts json response
func (ctx *Context) AcceptJSON() bool {
	return acceptsJSONRegex.MatchString(ctx.HeaderParam(HeaderAccept))
}

// UserAgent returns request client user agent string.
func (ctx *Context) UserAgent() string {
	return ctx.HeaderParam(HeaderUserAgent)
}

// Data return the implicit data in the context
func (ctx *Context) Data() map[interface{}]interface{} {
	if ctx.data == nil {
		ctx.data = make(map[interface{}]interface{})
	}
	return ctx.data
}

// GetData returns the stored data in this context.
func (ctx *Context) GetData(key interface{}) interface{} {
	if v, ok := ctx.data[key]; ok {
		return v
	}
	return nil
}

// SetData stores data with given key in this context.
// This data are only available in this context.
func (ctx *Context) SetData(key, val interface{}) {
	if ctx.data == nil {
		ctx.data = make(map[interface{}]interface{})
	}
	ctx.data[key] = val
}

// Contains checks if the key exists in the context.
func (ctx *Context) Contains(key interface{}) bool {
	_, ok := ctx.data[key]
	return ok
}

// Param returns the first value for the kinds of parameters.
// priority:
// path parameters > POST and PUT body parameters > URL query string values > header > cookie.Value.
//
// Param calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, Param returns the empty string.
// To access multiple values of the same key, call ParseForm and
// then inspect Request.Form directly.
func (ctx *Context) Param(key string) string {
	var value string
	value = ctx.BizParam(key)
	if len(value) > 0 {
		return value
	}
	value = ctx.R.Header.Get(key)
	if len(value) > 0 {
		return value
	}
	if cookie, err := ctx.R.Cookie(key); err != nil {
		return cookie.Value
	}
	return value
}

// BizParam returns the first value for the kinds of business parameters.
// priority:
// path parameters > POST and PUT body parameters > URL query string values.
//
// BizParam calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, BizParam returns the empty string.
// To access multiple values of the same key, call ParseForm and
// then inspect Request.Form directly.
func (ctx *Context) BizParam(key string) string {
	var value string
	value = ctx.pathParams.ByName(key)
	if len(value) > 0 {
		return value
	}
	if ctx.R.Form == nil {
		ctx.R.ParseMultipartForm(ctx.frame.config.multipartMaxMemory)
	}
	return ctx.R.FormValue(key)
}

// PathParam returns path param by key.
func (ctx *Context) PathParam(key string) string {
	return ctx.pathParams.ByName(key)
}

// PathParamAll returns whole path parameters.
func (ctx *Context) PathParamAll() Params {
	return ctx.pathParams
}

// ParseFormOrMulitForm parseForm or parseMultiForm based on Content-type
func (ctx *Context) ParseFormOrMulitForm(maxMemory int64) error {
	// Parse the body depending on the content type.
	if strings.Contains(ctx.HeaderParam(HeaderContentType), MIMEMultipartForm) {
		if err := ctx.R.ParseMultipartForm(maxMemory); err != nil {
			return errors.New("Error parsing request body:" + err.Error())
		}
	} else if err := ctx.R.ParseForm(); err != nil {
		return errors.New("Error parsing request body:" + err.Error())
	}
	return nil
}

// FormParam returns the first value for the named component of the POST or PUT ruest body.
// URL query parameters and path parameters are ignored.
// FormParam calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, FormParam returns the empty string.
func (ctx *Context) FormParam(key string) string {
	if ctx.R.PostForm == nil {
		ctx.R.ParseMultipartForm(ctx.frame.config.multipartMaxMemory)
	}
	return ctx.R.PostFormValue(key)
}

// FormParams returns the form field value with "[]string" for the provided key.
func (ctx *Context) FormParams(key string) []string {
	if ctx.R.PostForm == nil {
		ctx.R.ParseMultipartForm(ctx.frame.config.multipartMaxMemory)
	}
	return ctx.R.PostForm[key]
}

// FormParamAll returns the parsed form data from POST, PATCH,
// or PUT body parameters.
func (ctx *Context) FormParamAll() url.Values {
	if ctx.R.PostForm == nil {
		ctx.R.ParseMultipartForm(ctx.frame.config.multipartMaxMemory)
	}
	return ctx.R.PostForm
}

// QueryParam gets the first query value associated with the given key.
// If there are no values associated with the key, QueryParam returns
// the empty string.
func (ctx *Context) QueryParam(key string) string {
	if ctx.queryParams == nil {
		ctx.queryParams = ctx.R.URL.Query()
	}
	return ctx.queryParams.Get(key)
}

// QueryParams returns the query param with "[]string".
func (ctx *Context) QueryParams(key string) []string {
	if ctx.queryParams == nil {
		ctx.queryParams = ctx.R.URL.Query()
	}
	return ctx.queryParams[key]
}

// QueryParamAll returns all query params.
func (ctx *Context) QueryParamAll() url.Values {
	if ctx.queryParams == nil {
		ctx.queryParams = ctx.R.URL.Query()
	}
	return ctx.queryParams
}

// HeaderParam gets the first header value associated with the given key.
// If there are no values associated with the key, HeaderParam returns
// the empty string.
func (ctx *Context) HeaderParam(key string) string {
	return ctx.R.Header.Get(key)
}

// HeaderParamAll returns the whole ruest header.
func (ctx *Context) HeaderParamAll() http.Header {
	return ctx.R.Header
}

// CookieParam returns request cookie item string by a given key.
// if non-existed, return empty string.
func (ctx *Context) CookieParam(key string) string {
	cookie, err := ctx.R.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// SecureCookieParam Get secure cookie from request by a given key.
func (ctx *Context) SecureCookieParam(secret, key string) (string, bool) {
	val := ctx.CookieParam(key)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}

// FormFile returns the first file for the provided form key.
// FormFile calls ParseMultipartForm and ParseForm if necessary.
func (ctx *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.R.FormFile(key)
}

// SaveFile saves the file *Context.FormFile to UPLOAD_DIR,
// character "?" indicates that the original file name.
// for example newfname="a/?" -> UPLOAD_DIR/a/fname.
func (ctx *Context) SaveFile(key string, cover bool, newfname ...string) (fileUrl string, size int64, err error) {
	f, fh, err := ctx.R.FormFile(key)
	if err != nil {
		return
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil && err == nil {
			err = err2
		}
	}()

	// Sets the full file name
	var fullname string
	if len(newfname) == 0 {
		fullname = filepath.Join(UPLOAD_DIR, fh.Filename)
	} else {
		if strings.Contains(newfname[0], "?") {
			fullname = filepath.Join(UPLOAD_DIR, strings.Replace(newfname[0], "?", fh.Filename, -1))
		} else {
			fname := strings.TrimRight(newfname[0], ".")
			if filepath.Ext(fname) == "" {
				fullname = filepath.Join(UPLOAD_DIR, fname+filepath.Ext(fh.Filename))
			} else {
				fullname = filepath.Join(UPLOAD_DIR, fname)
			}
		}
	}

	// Create the completion file path
	p, _ := filepath.Split(fullname)
	err = os.MkdirAll(p, 0777)
	if err != nil {
		return
	}

	// If the file with the same name exists, add the suffix of the serial number
	idx := strings.LastIndex(fullname, filepath.Ext(fullname))
	_fullname := fullname
	for i := 2; utils.FileExists(_fullname) && !cover; i++ {
		_fullname = fmt.Sprintf("%s(%d)%s", fullname[:idx], i, fullname[idx:])
	}
	fullname = _fullname

	// Create the URL of the file
	fileUrl = "/" + strings.Replace(fullname, `\`, `/`, -1)

	// Save the file to local
	f2, err := os.OpenFile(fullname, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	size, err = io.Copy(f2, f)
	err3 := f2.Close()
	if err3 != nil && err == nil {
		err = err3
	}
	return
}

// Session returns current session item value by a given key.
// if non-existed, return nil.
func (ctx *Context) Session(key interface{}) interface{} {
	return ctx.CruSession.Get(key)
}

// CopyBody returns the raw request body data as bytes.
func (ctx *Context) CopyBody(MaxMemory int64) []byte {
	if ctx.R.Body == nil {
		return []byte{}
	}
	safe := &io.LimitedReader{R: ctx.R.Body, N: MaxMemory}
	requestbody, _ := ioutil.ReadAll(safe)
	ctx.R.Body.Close()
	bf := bytes.NewBuffer(requestbody)
	ctx.R.Body = ioutil.NopCloser(bf)
	ctx.RequestBody = requestbody
	return requestbody
}

// BizBind data from ctx.BizParam(key) to dest
// like /?id=123&isok=true&ft=1.2&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=astaxie
// var id int  ctx.BizBind(&id, "id")  id ==123
// var isok bool  ctx.BizBind(&isok, "isok")  isok ==true
// var ft float64  ctx.BizBind(&ft, "ft")  ft ==1.2
// ol := make([]int, 0, 2)  ctx.BizBind(&ol, "ol")  ol ==[1 2]
// ul := make([]string, 0, 2)  ctx.BizBind(&ul, "ul")  ul ==[str array]
// user struct{Name}  ctx.BizBind(&user, "user")  user == {Name:"astaxie"}
func (ctx *Context) BizBind(dest interface{}, key string) error {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return errors.New("thinkgo: non-pointer passed to Bind: " + key)
	}
	value = value.Elem()
	if !value.CanSet() {
		return errors.New("thinkgo: non-settable variable passed to Bind: " + key)
	}
	rv := ctx.bind(key, value.Type())
	if !rv.IsValid() {
		return errors.New("thinkgo: reflect value is empty")
	}
	value.Set(rv)
	return nil
}

func (ctx *Context) bind(key string, typ reflect.Type) reflect.Value {
	rv := reflect.Zero(typ)
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := ctx.BizParam(key)
		if len(val) == 0 {
			return rv
		}
		rv = ctx.bindInt(val, typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := ctx.BizParam(key)
		if len(val) == 0 {
			return rv
		}
		rv = ctx.bindUint(val, typ)
	case reflect.Float32, reflect.Float64:
		val := ctx.BizParam(key)
		if len(val) == 0 {
			return rv
		}
		rv = ctx.bindFloat(val, typ)
	case reflect.String:
		val := ctx.BizParam(key)
		if len(val) == 0 {
			return rv
		}
		rv = ctx.bindString(val, typ)
	case reflect.Bool:
		val := ctx.BizParam(key)
		if len(val) == 0 {
			return rv
		}
		rv = ctx.bindBool(val, typ)
	case reflect.Slice:
		rv = ctx.bindSlice(&ctx.R.Form, key, typ)
	case reflect.Struct:
		rv = ctx.bindStruct(&ctx.R.Form, key, typ)
	case reflect.Ptr:
		rv = ctx.bindPoint(key, typ)
	case reflect.Map:
		rv = ctx.bindMap(&ctx.R.Form, key, typ)
	}
	return rv
}

func (ctx *Context) bindValue(val string, typ reflect.Type) reflect.Value {
	rv := reflect.Zero(typ)
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv = ctx.bindInt(val, typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv = ctx.bindUint(val, typ)
	case reflect.Float32, reflect.Float64:
		rv = ctx.bindFloat(val, typ)
	case reflect.String:
		rv = ctx.bindString(val, typ)
	case reflect.Bool:
		rv = ctx.bindBool(val, typ)
	case reflect.Slice:
		rv = ctx.bindSlice(&url.Values{"": {val}}, "", typ)
	case reflect.Struct:
		rv = ctx.bindStruct(&url.Values{"": {val}}, "", typ)
	case reflect.Ptr:
		rv = ctx.bindPoint(val, typ)
	case reflect.Map:
		rv = ctx.bindMap(&url.Values{"": {val}}, "", typ)
	}
	return rv
}

func (ctx *Context) bindInt(val string, typ reflect.Type) reflect.Value {
	intValue, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetInt(intValue)
	return pValue.Elem()
}

func (ctx *Context) bindUint(val string, typ reflect.Type) reflect.Value {
	uintValue, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetUint(uintValue)
	return pValue.Elem()
}

func (ctx *Context) bindFloat(val string, typ reflect.Type) reflect.Value {
	floatValue, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return reflect.Zero(typ)
	}
	pValue := reflect.New(typ)
	pValue.Elem().SetFloat(floatValue)
	return pValue.Elem()
}

func (ctx *Context) bindString(val string, typ reflect.Type) reflect.Value {
	return reflect.ValueOf(val)
}

func (ctx *Context) bindBool(val string, typ reflect.Type) reflect.Value {
	val = strings.TrimSpace(strings.ToLower(val))
	switch val {
	case "true", "on", "1":
		return reflect.ValueOf(true)
	}
	return reflect.ValueOf(false)
}

type sliceValue struct {
	index int           // Index extracted from brackets.  If -1, no index was provided.
	value reflect.Value // the bound value for this slice element.
}

func (ctx *Context) bindSlice(params *url.Values, key string, typ reflect.Type) reflect.Value {
	maxIndex := -1
	numNoIndex := 0
	sliceValues := []sliceValue{}
	for reqKey, vals := range *params {
		if !strings.HasPrefix(reqKey, key+"[") {
			continue
		}
		// Extract the index, and the index where a sub-key starts. (e.g. field[0].subkey)
		index := -1
		leftBracket, rightBracket := len(key), strings.Index(reqKey[len(key):], "]")+len(key)
		if rightBracket > leftBracket+1 {
			index, _ = strconv.Atoi(reqKey[leftBracket+1 : rightBracket])
		}
		subKeyIndex := rightBracket + 1

		// Handle the indexed case.
		if index > -1 {
			if index > maxIndex {
				maxIndex = index
			}
			sliceValues = append(sliceValues, sliceValue{
				index: index,
				value: ctx.bind(reqKey[:subKeyIndex], typ.Elem()),
			})
			continue
		}

		// It's an un-indexed element.  (e.g. element[])
		numNoIndex += len(vals)
		for _, val := range vals {
			// Unindexed values can only be direct-bound.
			sliceValues = append(sliceValues, sliceValue{
				index: -1,
				value: ctx.bindValue(val, typ.Elem()),
			})
		}
	}
	resultArray := reflect.MakeSlice(typ, maxIndex+1, maxIndex+1+numNoIndex)
	for _, sv := range sliceValues {
		if sv.index != -1 {
			resultArray.Index(sv.index).Set(sv.value)
		} else {
			resultArray = reflect.Append(resultArray, sv.value)
		}
	}
	return resultArray
}

func (ctx *Context) bindStruct(params *url.Values, key string, typ reflect.Type) reflect.Value {
	result := reflect.New(typ).Elem()
	fieldValues := make(map[string]reflect.Value)
	for reqKey, val := range *params {
		var fieldName string
		if strings.HasPrefix(reqKey, key+".") {
			fieldName = reqKey[len(key)+1:]
		} else if strings.HasPrefix(reqKey, key+"[") && reqKey[len(reqKey)-1] == ']' {
			fieldName = reqKey[len(key)+1 : len(reqKey)-1]
		} else {
			continue
		}

		if _, ok := fieldValues[fieldName]; !ok {
			// Time to bind this field.  Get it and make sure we can set it.
			fieldValue := result.FieldByName(fieldName)
			if !fieldValue.IsValid() {
				continue
			}
			if !fieldValue.CanSet() {
				continue
			}
			boundVal := ctx.bindValue(val[0], fieldValue.Type())
			fieldValue.Set(boundVal)
			fieldValues[fieldName] = boundVal
		}
	}

	return result
}

func (ctx *Context) bindPoint(key string, typ reflect.Type) reflect.Value {
	return ctx.bind(key, typ.Elem()).Addr()
}

func (ctx *Context) bindMap(params *url.Values, key string, typ reflect.Type) reflect.Value {
	var (
		result    = reflect.MakeMap(typ)
		keyType   = typ.Key()
		valueType = typ.Elem()
	)
	for paramName, values := range *params {
		if !strings.HasPrefix(paramName, key+"[") || paramName[len(paramName)-1] != ']' {
			continue
		}

		key := paramName[len(key)+1 : len(paramName)-1]
		result.SetMapIndex(ctx.bindValue(key, keyType), ctx.bindValue(values[0], valueType))
	}
	return result
}

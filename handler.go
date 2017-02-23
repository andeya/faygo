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
	"errors"
	"net/http"
	"reflect"
	"sort"

	"github.com/henrylee2cn/faygo/apiware"
)

type (
	// Handler is the main Faygo Handler interface.
	Handler interface {
		Serve(ctx *Context) error
	}
	// APIHandler is the Faygo Handler interface,
	// which is implemented by a struct with API descriptor information.
	// It is an intelligent Handler of binding parameters.
	APIHandler interface {
		Handler
		APIDoc
	}
	// HandlerWithBody is the Faygo APIHandler interface but with DecodeBody method.
	HandlerWithBody interface {
		Handler
		Bodydecoder // Decode params from request body
	}
	// Bodydecoder is an interface to customize decoding operation
	Bodydecoder interface {
		Decode(dest reflect.Value, body []byte) error
	}
	// HandlerWithoutPath is handler without binding path parameter for middleware.
	HandlerWithoutPath interface {
		Handler
	}
	// APIDoc provides the API's note, result or parameters information.
	APIDoc interface {
		Doc() Doc
	}
	// ParamInfo is the request parameter information
	ParamInfo struct {
		Name     string      // Parameter name
		In       string      // The position of the parameter
		Required bool        // Is a required parameter
		Model    interface{} // A parameter value that is used to infer a value type and as a default value
		Desc     string      // Description
	}
	// Doc api information
	Doc struct {
		Note   string      `json:"note" xml:"note"`
		Return interface{} `json:"return,omitempty" xml:"return,omitempty"`
		Params []ParamInfo `json:"params,omitempty" xml:"params,omitempty"`
	}
	// Notes implementation notes of a response
	Notes struct {
		Note   string      `json:"note" xml:"note"`
		Return interface{} `json:"return,omitempty" xml:"return,omitempty"`
	}
	// JSONMsg is commonly used to return JSON format response.
	JSONMsg struct {
		Code int         `json:"code" xml:"code"`                     // the status code of the business process (required)
		Info interface{} `json:"info,omitempty" xml:"info,omitempty"` // response's schema and example value (optional)
	}
	// apiHandler is an intelligent Handler of binding parameters.
	apiHandler struct {
		paramsAPI   *apiware.ParamsAPI
		paramTypes  []reflect.Type
		paramValues []reflect.Value
		handler     Handler
	}
	// HandlerFunc type is an adapter to allow the use of
	// ordinary functions as HTTP handlers.  If f is a function
	// with the appropriate signature, HandlerFunc(f) is a
	// Handler that calls f.
	HandlerFunc func(ctx *Context) error
	// HandlerChain is the chain of handlers for a request.
	HandlerChain []Handler
	// ErrorFunc replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further
	// writes are done to ctx.
	// The error message should be plain text.
	ErrorFunc func(ctx *Context, errStr string, status int)
	// BinderrorFunc is called when binding or validation apiHandler parameters are wrong.
	BinderrorFunc func(ctx *Context, err error)
)

// Serve implements the Handler, is like ServeHTTP but for Faygo.
func (h HandlerFunc) Serve(ctx *Context) error {
	return h(ctx)
}

// common errors
var (
	ErrNotStructPtr   = errors.New("handler must be a structure type or a structure pointer type")
	ErrNoParamHandler = errors.New("handler does not define any parameter tags")
)
var _ APIDoc = new(apiHandler)

// ToAPIHandler tries converts it to an *apiHandler.
func ToAPIHandler(handler Handler) (*apiHandler, error) {
	v := reflect.Indirect(reflect.ValueOf(handler))
	if v.Kind() != reflect.Struct {
		return nil, ErrNotStructPtr
	}

	var structPointer = v.Addr().Interface()
	var bodydecoder = global.bodydecoder
	if h, ok := structPointer.(HandlerWithBody); ok {
		bodydecoder = h.Decode
	}

	paramsAPI, err := apiware.NewParamsAPI(structPointer, global.paramNameMapper, bodydecoder)
	if err != nil {
		return nil, err
	}
	if paramsAPI.Number() == 0 {
		return nil, ErrNoParamHandler
	}

	_, paramValues := paramsAPI.NewReceiver()
	var paramTypes = make([]reflect.Type, len(paramValues))
	for i, v := range paramValues {
		paramTypes[i] = v.Type()
	}
	// Reduce the creation of unnecessary field paramValues.
	return &apiHandler{
		paramsAPI:  paramsAPI,
		paramTypes: paramTypes,
		handler:    structPointer.(Handler),
	}, nil
}

// IsHandlerWithoutPath verifies that the Handler is an HandlerWithoutPath.
func IsHandlerWithoutPath(handler Handler) bool {
	v := reflect.Indirect(reflect.ValueOf(handler))
	if v.Kind() != reflect.Struct {
		return true
	}
	paramsAPI, err := apiware.NewParamsAPI(v.Addr().Interface(), nil, nil)
	if err != nil {
		return true
	}
	for _, param := range paramsAPI.Params() {
		if param.In() == "path" {
			return false
		}
	}
	return true
}

// Serve implements the APIHandler, is like ServeHTTP but for Faygo
func (h *apiHandler) Serve(ctx *Context) error {
	return h.handler.Serve(ctx)
}

// Doc returns the API's note, result or parameters information.
func (h *apiHandler) Doc() Doc {
	var doc Doc
	if d, ok := h.handler.(APIDoc); ok {
		doc = d.Doc()
	}
	for _, param := range h.paramsAPI.Params() {
		var had bool
		var info = ParamInfo{
			Name:     param.Name(),
			In:       param.In(),
			Required: param.IsRequired(),
			Desc:     param.Description(),
			Model:    param.Raw(),
		}
		for i, p := range doc.Params {
			if p.Name == info.Name {
				doc.Params[i] = info
				had = true
				break
			}
		}
		if !had {
			doc.Params = append(doc.Params, info)
		}
	}
	return doc
}

// Create a new `*apiHandler` by itself.
func (h *apiHandler) new() *apiHandler {
	h2 := &apiHandler{
		paramsAPI:  h.paramsAPI,
		paramTypes: h.paramTypes,
	}
	var object interface{}
	object, h2.paramValues = h.paramsAPI.NewReceiver()
	h2.handler = object.(Handler)
	return h2
}

// Bind the request path params to `apiHandler.handler`.
func (h *apiHandler) bind(req *http.Request, pathParams PathParams) error {
	return h.paramsAPI.BindFields(h.paramValues, req, pathParams)
}

// Reset all fields to a value of zero
func (h *apiHandler) reset() {
	for i, typ := range h.paramTypes {
		h.paramValues[i].Set(reflect.Zero(typ))
	}
}

// Get distinct and sorted parameters information.
func distinctAndSortedParamInfos(infos []ParamInfo) []ParamInfo {
	infoMap := make(map[string]ParamInfo, len(infos))
	ks := make([]string, 0, len(infos))
	for _, info := range infos {
		k := info.Name + "<\r-\n-\t>" + info.In
		ks = append(ks, k)
		// Filter duplicate parameters, and maximize access to information.
		if newinfo, ok := infoMap[k]; ok {
			if !newinfo.Required && info.Required {
				newinfo.Required = info.Required
			}
			if len(newinfo.Desc) == 0 && len(info.Desc) > 0 {
				newinfo.Desc = info.Desc
			}
			infoMap[k] = newinfo
			continue
		}
		infoMap[k] = info
	}
	sort.Strings(ks)
	newinfos := make([]ParamInfo, 0, len(ks))
	for _, k := range ks {
		newinfos = append(newinfos, infoMap[k])
	}
	return newinfos
}

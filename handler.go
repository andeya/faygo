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
	"net/http"
	"reflect"
	"sort"

	"github.com/henrylee2cn/apiware"
)

type (
	// Handler is the main Thinkgo Handler interface.
	Handler interface {
		Serve(ctx *Context) error
	}

	// APIHandler is the Thinkgo Handler interface,
	// which is implemented by a struct with API descriptor information.
	// It is an intelligent Handler of binding parameters.
	APIHandler interface {
		Handler
	}

	// APIHandlerWithBody is the Thinkgo APIHandler interface but with DecodeBody method.
	APIHandlerWithBody interface {
		APIHandler
		BodyDecoder // Decode params from request body
	}

	// BodyDecoder is an interface to customize decoding operation
	BodyDecoder interface {
		Decode(dest reflect.Value, body []byte) error
	}

	// Handler without binding path parameter for middleware, when Handler is implemented by a APIHandler.
	HandlerWithoutPath interface {
		Handler
	}

	// handlerStruct implemented `APIHandler` interface, and possibly implemented the `APIHandlerWithBody` interface.
	// It is an intelligent Handler of binding parameters.
	handlerStruct struct {
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

	// The chain of handlers for a request
	HandlerChain []Handler

	// ErrorFunc replies to the request with the specified error message and HTTP code.
	// It does not otherwise end the request; the caller should ensure no further
	// writes are done to ctx.
	// The error message should be plain text.
	ErrorFunc     func(ctx *Context, errStr string, status int)
	BindErrorFunc func(ctx *Context, err error)
)

// Serve implements the Handler, is like ServeHTTP but for Thinkgo.
func (h HandlerFunc) Serve(ctx *Context) error {
	return h(ctx)
}

// Verify that the Handler is an APIHandler and convert it.
// If fails, returns nil.
func ToAPIHandler(handler Handler) APIHandler {
	v := reflect.Indirect(reflect.ValueOf(handler))
	if v.Kind() != reflect.Struct {
		return nil
	}
	return v.Addr().Interface().(APIHandler)
}

// Verify that the Handler is an HandlerWithoutPath.
func IsHandlerWithoutPath(handler Handler) bool {
	apiHandler := ToAPIHandler(handler)
	if apiHandler == nil {
		return true
	}
	paramsAPI, err := apiware.NewParamsAPI(apiHandler, nil, nil)
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

// Create a `*handlerStruct`.
func newHandlerStruct(structPointer APIHandler, paramMapping apiware.ParamNameFunc) (*handlerStruct, error) {
	var bodyDecodeFunc = Global.bodyDecodeFunc
	if h, ok := structPointer.(APIHandlerWithBody); ok {
		bodyDecodeFunc = h.Decode
	}

	paramsAPI, err := apiware.NewParamsAPI(structPointer, paramMapping, bodyDecodeFunc)
	if err != nil {
		return nil, err
	}

	_, paramValues := paramsAPI.NewReceiver()
	var paramTypes = make([]reflect.Type, len(paramValues))
	for i, v := range paramValues {
		paramTypes[i] = v.Type()
	}
	// Reduce the creation of unnecessary field paramValues.
	return &handlerStruct{
		paramsAPI:  paramsAPI,
		paramTypes: paramTypes,
		handler:    structPointer,
	}, nil
}

// Create a new `*handlerStruct` by itself.
func (h *handlerStruct) new() *handlerStruct {
	h2 := &handlerStruct{
		paramsAPI:  h.paramsAPI,
		paramTypes: h.paramTypes,
	}
	var object interface{}
	object, h2.paramValues = h.paramsAPI.NewReceiver()
	h2.handler = object.(Handler)
	return h2
}

// Bind the request params to `handlerStruct.handler`.
func (h *handlerStruct) bind(req *http.Request, params Params) error {
	return h.paramsAPI.BindFields(h.paramValues, req, params)
}

// Reset all fields to a value of zero
func (h *handlerStruct) reset() {
	for i, typ := range h.paramTypes {
		h.paramValues[i].Set(reflect.Zero(typ))
	}
}

// Serve implements the APIHandler, is like ServeHTTP but for Thinkgo
func (h *handlerStruct) Serve(ctx *Context) error {
	return h.handler.Serve(ctx)
}

/*
Param tag value description:
    tag   |   key    | required |     value     |   desc
    ------|----------|----------|---------------|----------------------------------
    param |    in    | only one |     path      | (position of param) if `required` is unsetted, auto set it. e.g. url: "http://www.abc.com/a/{path}"
    param |    in    | only one |     query     | (position of param) e.g. url: "http://www.abc.com/a?b={query}"
    param |    in    | only one |     formData  | (position of param) e.g. "request body: a=123&b={formData}"
    param |    in    | only one |     body      | (position of param) request body can be any content
    param |    in    | only one |     header    | (position of param) request header info
    param |    in    | only one |     cookie    | (position of param) request cookie info, support: `http.Cookie`,`fasthttp.Cookie`,`string`,`[]byte`
    param |   name   |    no    |  (e.g. "id")  | specify request param`s name
    param | required |    no    |   required    | request param is required
    param |   desc   |    no    |  (e.g. "id")  | request param description
    param |   len    |    no    | (e.g. 3:6, 3) | length range of param's value
    param |   range  |    no    |  (e.g. 0:10)  | numerical range of param's value
    param |  nonzero |    no    |    nonzero    | param`s value can not be zero
    param |   maxmb  |    no    |   (e.g. 32)   | when request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
    regexp|          |    no    |(e.g. "^\\w+$")| param value can not be null
    err   |          |    no    |(e.g. "incorrect password format")| customize the prompt for validation error

    NOTES:
        1. the binding object must be a struct pointer
        2. the binding struct's field can not be a pointer
        3. `regexp` or `param` tag is only usable when `param:"type(xxx)"` is exist
        4. if the `param` tag is not exist, anonymous field will be parsed
        5. when the param's position(`in`) is `formData` and the field's type is `multipart.FileHeader`, the param receives file uploaded
        6. if param's position(`in`) is `cookie`, field's type must be `http.Cookie`
        7. param tags `in(formData)` and `in(body)` can not exist at the same time
        8. there should not be more than one `in(body)` param tag

List of supported param value types:
    base    |   slice    | special
    --------|------------|-------------------------------------------------------
    string  |  []string  | [][]byte
    byte    |  []byte    | [][]uint8
    uint8   |  []uint8   | multipart.FileHeader (only for `formData` param)
    bool    |  []bool    | http.Cookie (only for `net/http`'s `cookie` param)
    int     |  []int     | fasthttp.Cookie (only for `fasthttp`'s `cookie` param)
    int8    |  []int8    | struct (struct type only for `body` param or as an anonymous field to extend params)
    int16   |  []int16   |
    int32   |  []int32   |
    int64   |  []int64   |
    uint8   |  []uint8   |
    uint16  |  []uint16  |
    uint32  |  []uint32  |
    uint64  |  []uint64  |
    float32 |  []float32 |
    float64 |  []float64 |
*/
type (
	// request parameter information
	ParamInfo struct {
		Name string
		// the position of the parameter
		In       string
		Required bool
		Desc     string
		// a parameter value that is used to infer a value type and as a default value
		Model interface{}
	}
	Doc interface {
		Notes() Notes
	}
	// response description
	// Return struct {
	// 	Code         int         // HTTP status code (required)
	// 	Description  string      // response's reason (optional)
	// 	ExampleValue interface{} // response's schema and example value (optional)
	// 	Headers      interface{} // response's headers (optional)
	// }

	// implementation notes of a response
	Notes struct {
		Note   string
		Return interface{}
	}
	JSONMsg struct {
		Code int         `json:"code"`           // the status code of the business process (required)
		Info interface{} `json:"info,omitempty"` // response's schema and example value (optional)
	}
)

// Only the original instance is invoked.
func (h *handlerStruct) paramInfos() []ParamInfo {
	params := h.paramsAPI.Params()
	infos := make([]ParamInfo, 0, len(params))
	for _, param := range params {
		infos = append(infos, ParamInfo{
			Name:     param.Name(),
			In:       param.In(),
			Required: param.IsRequired(),
			Desc:     param.Description(),
			Model:    param.Raw(),
		})
	}
	return infos
	// return distinctAndSortedParamInfos(infos)
}

// Only the original instance is invoked.
func (h *handlerStruct) getNotes() *Notes {
	if doc, ok := h.handler.(Doc); ok {
		n := doc.Notes()
		return &n
	}
	return nil
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

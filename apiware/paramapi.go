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

package apiware

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	// "github.com/valyala/fasthttp"
)

type (
	// ParamsAPI defines a parameter model for an web api.
	ParamsAPI struct {
		name   string
		params []*Param
		//used to create a new struct (non-pointer)
		structType reflect.Type
		//the raw struct pointer
		rawStructPointer interface{}
		// rawStructPointer value bytes
		defaultValues []byte
		// create param name from struct field name
		paramNameMapper ParamNameMapper
		// decode params from request body
		bodydecoder Bodydecoder
		//when request Content-Type is multipart/form-data, the max memory for body.
		maxMemory int64
	}

	// Schema is a collection of ParamsAPI
	Schema struct {
		lib map[string]*ParamsAPI
		sync.RWMutex
	}

	// ParamNameMapper maps param name from struct param name
	ParamNameMapper func(fieldName string) (paramName string)

	// Bodydecoder decodes params from request body.
	Bodydecoder func(dest reflect.Value, body []byte) error
)

var (
	defaultSchema = &Schema{
		lib: map[string]*ParamsAPI{},
	}
)

// NewParamsAPI parses and store the struct object, requires a struct pointer,
// if `paramNameMapper` is nil, `paramNameMapper=toSnake`,
// if `bodydecoder` is nil, `bodydecoder=bodyJONS`,
func NewParamsAPI(
	structPointer interface{},
	paramNameMapper ParamNameMapper,
	bodydecoder Bodydecoder,
	useDefaultValues bool,
) (
	*ParamsAPI,
	error,
) {
	name := reflect.TypeOf(structPointer).String()
	v := reflect.ValueOf(structPointer)
	if v.Kind() != reflect.Ptr {
		return nil, NewError(name, "*", "the binding object must be a struct pointer")
	}
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return nil, NewError(name, "*", "the binding object must be a struct pointer")
	}
	var paramsAPI = &ParamsAPI{
		name:             name,
		params:           []*Param{},
		structType:       v.Type(),
		rawStructPointer: structPointer,
	}
	if paramNameMapper != nil {
		paramsAPI.paramNameMapper = paramNameMapper
	} else {
		paramsAPI.paramNameMapper = toSnake
	}
	if bodydecoder != nil {
		paramsAPI.bodydecoder = bodydecoder
	} else {
		paramsAPI.bodydecoder = bodyJONS
	}
	err := paramsAPI.addFields([]int{}, paramsAPI.structType, v)
	if err != nil {
		return nil, err
	}

	if useDefaultValues && !reflect.DeepEqual(reflect.New(paramsAPI.structType).Interface(), paramsAPI.rawStructPointer) {
		buf := bytes.NewBuffer(nil)
		err = gob.NewEncoder(buf).EncodeValue(v)
		if err == nil {
			paramsAPI.defaultValues = buf.Bytes()
		}
	}
	defaultSchema.set(paramsAPI)
	return paramsAPI, nil
}

// Register is similar to a `NewParamsAPI`, but only return error.
// Parse and store the struct object, requires a struct pointer,
// if `paramNameMapper` is nil, `paramNameMapper=toSnake`,
// if `bodydecoder` is nil, `bodydecoder=bodyJONS`,
func Register(
	structPointer interface{},
	paramNameMapper ParamNameMapper,
	bodydecoder Bodydecoder,
	useDefaultValues bool,
) error {
	_, err := NewParamsAPI(structPointer, paramNameMapper, bodydecoder, useDefaultValues)
	return err
}

func (paramsAPI *ParamsAPI) addFields(parentIndexPath []int, t reflect.Type, v reflect.Value) error {
	var err error
	var maxMemoryMB int64
	var hasFormData, hasBody bool
	var deep = len(parentIndexPath) + 1
	for i := 0; i < t.NumField(); i++ {
		indexPath := make([]int, deep)
		copy(indexPath, parentIndexPath)
		indexPath[deep-1] = i

		var field = t.Field(i)
		tag, ok := field.Tag.Lookup(TAG_PARAM)
		if !ok {
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				if err = paramsAPI.addFields(indexPath, field.Type, v.Field(i)); err != nil {
					return err
				}
			}
			continue
		}

		if tag == TAG_IGNORE_PARAM {
			continue
		}
		if field.Type.Kind() == reflect.Ptr && field.Type.String() != fileTypeString && field.Type.String() != cookieTypeString {
			return NewError(t.String(), field.Name, "field can not be a pointer")
		}

		var value = v.Field(i)
		if !value.CanSet() {
			return NewError(t.String(), field.Name, "field can not be a unexported field")
		}

		var parsedTags = ParseTags(tag)
		var paramPosition = parsedTags[KEY_IN]
		var paramTypeString = field.Type.String()

		switch paramTypeString {
		case fileTypeString, filesTypeString, fileTypeString2, filesTypeString2:
			if paramPosition != "formData" {
				return NewError(t.String(), field.Name, "when field type is `"+paramTypeString+"`, tag `in` value must be `formData`")
			}
		case cookieTypeString, cookieTypeString2 /*, fasthttpCookieTypeString*/ :
			if paramPosition != "cookie" {
				return NewError(t.String(), field.Name, "when field type is `"+paramTypeString+"`, tag `in` value must be `cookie`")
			}
		}

		switch paramPosition {
		case "formData":
			if hasBody {
				return NewError(t.String(), field.Name, "tags of `in(formData)` and `in(body)` can not exist at the same time")
			}
			hasFormData = true
		case "body":
			if hasFormData {
				return NewError(t.String(), field.Name, "tags of `in(formData)` and `in(body)` can not exist at the same time")
			}
			if hasBody {
				return NewError(t.String(), field.Name, "there should not be more than one tag `in(body)`")
			}
			hasBody = true
		case "path":
			parsedTags[KEY_REQUIRED] = KEY_REQUIRED
		// case "cookie":
		// 	switch paramTypeString {
		// 	case cookieTypeString, fasthttpCookieTypeString, stringTypeString, bytesTypeString, bytes2TypeString:
		// 	default:
		// 		return NewError( t.String(),field.Name, "invalid field type for `in(cookie)`, refer to the following: `http.Cookie`, `fasthttp.Cookie`, `string`, `[]byte` or `[]uint8`")
		// 	}
		default:
			if !TagInValues[paramPosition] {
				return NewError(t.String(), field.Name, "invalid tag `in` value, refer to the following: `path`, `query`, `formData`, `body`, `header` or `cookie`")
			}
		}
		if _, ok := parsedTags[KEY_LEN]; ok {
			switch paramTypeString {
			case "string", "[]string", "[]int", "[]int8", "[]int16", "[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64", "[]float32", "[]float64":
			default:
				return NewError(t.String(), field.Name, "invalid `len` tag for non-string or non-basetype-slice field")
			}
		}
		if _, ok := parsedTags[KEY_RANGE]; ok {
			switch paramTypeString {
			case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			case "[]int", "[]int8", "[]int16", "[]int32", "[]int64", "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64", "[]float32", "[]float64":
			default:
				return NewError(t.String(), field.Name, "invalid `range` tag for non-number field")
			}
		}
		if _, ok := parsedTags[KEY_REGEXP]; ok {
			if paramTypeString != "string" && paramTypeString != "[]string" {
				return NewError(t.String(), field.Name, "invalid `"+KEY_REGEXP+"` tag for non-string field")
			}
		}
		if a, ok := parsedTags[KEY_MAXMB]; ok {
			i, err := strconv.ParseInt(a, 10, 64)
			if err != nil {
				return NewError(t.String(), field.Name, "invalid `maxmb` tag, it must be positive integer")
			}
			if i > maxMemoryMB {
				maxMemoryMB = i
			}
		}

		fd := &Param{
			apiName:   paramsAPI.name,
			indexPath: indexPath,
			tags:      parsedTags,
			rawTag:    field.Tag,
			rawValue:  value,
		}

		if errStr, ok := fd.tags[KEY_ERR]; ok {
			fd.err = errors.New(errStr)
		}

		// fmt.Printf("%#v\n", fd.tags)

		if fd.name, ok = parsedTags[KEY_NAME]; !ok {
			fd.name = paramsAPI.paramNameMapper(field.Name)
		}
		if paramPosition == "header" {
			fd.name = textproto.CanonicalMIMEHeaderKey(fd.name)
		}

		fd.isFile = paramTypeString == fileTypeString || paramTypeString == filesTypeString || paramTypeString == fileTypeString2 || paramTypeString == filesTypeString2

		_, fd.isRequired = parsedTags[KEY_REQUIRED]
		_, hasNonzero := parsedTags[KEY_NONZERO]
		if !fd.isRequired && (hasNonzero || len(parsedTags[KEY_RANGE]) > 0) {
			fd.isRequired = true
		}

		if err = fd.makeVerifyFuncs(); err != nil {
			return NewError(t.String(), field.Name, "initial validation failed:"+err.Error())
		}

		paramsAPI.params = append(paramsAPI.params, fd)
	}
	if maxMemoryMB > 0 {
		paramsAPI.maxMemory = maxMemoryMB * MB
	} else {
		paramsAPI.maxMemory = defaultMaxMemory
	}
	return nil
}

// GetParamsAPI gets the `*ParamsAPI` object according to the type name
func GetParamsAPI(paramsAPIName string) (*ParamsAPI, error) {
	paramsAPI, ok := defaultSchema.get(paramsAPIName)
	if !ok {
		return nil, errors.New("struct `" + paramsAPIName + "` is not registered")
	}
	return paramsAPI, nil
}

// SetParamsAPI caches `*ParamsAPI`
func SetParamsAPI(paramsAPI *ParamsAPI) {
	defaultSchema.set(paramsAPI)
}

func (schema *Schema) get(paramsAPIName string) (*ParamsAPI, bool) {
	schema.RLock()
	defer schema.RUnlock()
	paramsAPI, ok := schema.lib[paramsAPIName]
	return paramsAPI, ok
}

func (schema *Schema) set(paramsAPI *ParamsAPI) {
	schema.Lock()
	schema.lib[paramsAPI.name] = paramsAPI
	defer schema.Unlock()
}

// Name gets the name
func (paramsAPI *ParamsAPI) Name() string {
	return paramsAPI.name
}

// Params gets the parameter information
func (paramsAPI *ParamsAPI) Params() []*Param {
	return paramsAPI.params
}

// Number returns the number of parameters to be bound
func (paramsAPI *ParamsAPI) Number() int {
	return len(paramsAPI.params)
}

// Raw returns the ParamsAPI's original value
func (paramsAPI *ParamsAPI) Raw() interface{} {
	return paramsAPI.rawStructPointer
}

// MaxMemory gets maxMemory
// when request Content-Type is multipart/form-data, the max memory for body.
func (paramsAPI *ParamsAPI) MaxMemory() int64 {
	return paramsAPI.maxMemory
}

// SetMaxMemory sets maxMemory for the request which Content-Type is multipart/form-data.
func (paramsAPI *ParamsAPI) SetMaxMemory(maxMemory int64) {
	paramsAPI.maxMemory = maxMemory
}

// NewReceiver creates a new struct pointer and the field's values  for its receive parameterste it.
func (paramsAPI *ParamsAPI) NewReceiver() (interface{}, []reflect.Value) {
	object := reflect.New(paramsAPI.structType)
	if len(paramsAPI.defaultValues) > 0 {
		// fmt.Printf("setting default value: %s\n", paramsAPI.structType.String())
		de := gob.NewDecoder(bytes.NewReader(paramsAPI.defaultValues))
		err := de.DecodeValue(object.Elem())
		if err != nil {
			panic(err)
		}
	}
	return object.Interface(), paramsAPI.fieldsForBinding(object.Elem())
}

func (paramsAPI *ParamsAPI) fieldsForBinding(structElem reflect.Value) []reflect.Value {
	count := len(paramsAPI.params)
	fields := make([]reflect.Value, count)
	for i := 0; i < count; i++ {
		value := structElem
		param := paramsAPI.params[i]
		for _, index := range param.indexPath {
			value = value.Field(index)
		}
		fields[i] = value
	}
	return fields
}

// BindByName binds the net/http request params to a new struct and validate it.
func BindByName(
	paramsAPIName string,
	req *http.Request,
	pathParams KV,
) (
	interface{},
	error,
) {
	paramsAPI, err := GetParamsAPI(paramsAPIName)
	if err != nil {
		return nil, err
	}
	return paramsAPI.BindNew(req, pathParams)
}

// Bind binds the net/http request params to the `structPointer` param and validate it.
// note: structPointer must be struct pointer.
func Bind(
	structPointer interface{},
	req *http.Request,
	pathParams KV,
) error {
	paramsAPI, err := GetParamsAPI(reflect.TypeOf(structPointer).String())
	if err != nil {
		return err
	}
	return paramsAPI.BindAt(structPointer, req, pathParams)
}

// BindAt binds the net/http request params to a struct pointer and validate it.
// note: structPointer must be struct pointer.
func (paramsAPI *ParamsAPI) BindAt(
	structPointer interface{},
	req *http.Request,
	pathParams KV,
) error {
	name := reflect.TypeOf(structPointer).String()
	if name != paramsAPI.name {
		return errors.New("the structPointer's type `" + name + "` does not match type `" + paramsAPI.name + "`")
	}
	return paramsAPI.BindFields(
		paramsAPI.fieldsForBinding(reflect.ValueOf(structPointer).Elem()),
		req,
		pathParams,
	)
}

// BindNew binds the net/http request params to a struct pointer and validate it.
func (paramsAPI *ParamsAPI) BindNew(
	req *http.Request,
	pathParams KV,
) (
	interface{},
	error,
) {
	structPrinter, fields := paramsAPI.NewReceiver()
	err := paramsAPI.BindFields(fields, req, pathParams)
	return structPrinter, err
}

// RawBind binds the net/http request params to the original struct pointer and validate it.
func (paramsAPI *ParamsAPI) RawBind(
	req *http.Request,
	pathParams KV,
) (
	interface{},
	error,
) {
	var fields []reflect.Value
	for _, param := range paramsAPI.params {
		fields = append(fields, param.rawValue)
	}
	err := paramsAPI.BindFields(fields, req, pathParams)
	return paramsAPI.rawStructPointer, err
}

// BindFields binds the net/http request params to a struct and validate it.
// Must ensure that the param `fields` matches `paramsAPI.params`.
func (paramsAPI *ParamsAPI) BindFields(
	fields []reflect.Value,
	req *http.Request,
	pathParams KV,
) (
	err error,
) {
	if pathParams == nil {
		pathParams = Map(map[string]string{})
	}
	if req.Form == nil {
		req.ParseMultipartForm(paramsAPI.maxMemory)
	}
	var queryValues url.Values
	defer func() {
		if p := recover(); p != nil {
			err = NewError(paramsAPI.name, "?", fmt.Sprint(p))
		}
	}()

	for i, param := range paramsAPI.params {
		value := fields[i]
		switch param.In() {
		case "path":
			paramValue, ok := pathParams.Get(param.name)
			if !ok {
				return param.myError("missing path param")
			}
			// fmt.Printf("paramName:%s\nvalue:%#v\n\n", param.name, paramValue)
			if err = convertAssign(value, []string{paramValue}); err != nil {
				return param.myError(err.Error())
			}

		case "query":
			if queryValues == nil {
				queryValues, err = url.ParseQuery(req.URL.RawQuery)
				if err != nil {
					queryValues = make(url.Values)
				}
			}
			paramValues, ok := queryValues[param.name]
			if ok {
				if err = convertAssign(value, paramValues); err != nil {
					return param.myError(err.Error())
				}
			} else if param.IsRequired() {
				return param.myError("missing query param")
			}

		case "formData":
			// Can not exist with `body` param at the same time
			if param.IsFile() {
				if req.MultipartForm != nil {
					fhs := req.MultipartForm.File[param.name]
					if len(fhs) == 0 {
						if param.IsRequired() {
							return param.myError("missing formData param")
						}
						continue
					}
					typ := value.Type()
					switch typ.String() {
					case fileTypeString:
						value.Set(reflect.ValueOf(fhs[0]))
					case fileTypeString2:
						value.Set(reflect.ValueOf(fhs[0]).Elem())
					case filesTypeString:
						value.Set(reflect.ValueOf(fhs))
					case filesTypeString2:
						fhs2 := make([]multipart.FileHeader, len(fhs))
						for i, fh := range fhs {
							fhs2[i] = *fh
						}
						value.Set(reflect.ValueOf(fhs2))
					default:
						return param.myError(
							"the param type is incorrect, reference: " +
								fileTypeString +
								"," + filesTypeString,
						)
					}
				} else if param.IsRequired() {
					return param.myError("missing formData param")
				}
				continue
			}

			paramValues, ok := req.PostForm[param.name]
			if ok {
				if err = convertAssign(value, paramValues); err != nil {
					return param.myError(err.Error())
				}
			} else if param.IsRequired() {
				return param.myError("missing formData param")
			}

		case "body":
			// Theoretically there should be at most one `body` param, and can not exist with `formData` at the same time
			var body []byte
			body, err = ioutil.ReadAll(req.Body)
			req.Body.Close()
			if err == nil {
				if err = paramsAPI.bodydecoder(value, body); err != nil {
					return param.myError(err.Error())
				}
			} else if param.IsRequired() {
				return param.myError("missing body param")
			}

		case "header":
			paramValues, ok := req.Header[param.name]
			if ok {
				if err = convertAssign(value, paramValues); err != nil {
					return param.myError(err.Error())
				}
			} else if param.IsRequired() {
				return param.myError("missing header param")
			}

		case "cookie":
			c, _ := req.Cookie(param.name)
			if c != nil {
				switch value.Type().String() {
				case cookieTypeString:
					value.Set(reflect.ValueOf(c))
				case cookieTypeString2:
					value.Set(reflect.ValueOf(c).Elem())
				default:
					if err = convertAssign(value, []string{c.Value}); err != nil {
						return param.myError(err.Error())
					}
				}
			} else if param.IsRequired() {
				return param.myError("missing cookie param")
			}
		}
		if err = param.validate(value); err != nil {
			return err
		}
	}
	return
}

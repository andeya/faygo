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

// Package swagger struct definition
package swagger

import (
	"fmt"
	"reflect"
	"strings"
)

// Version show the current swagger version
const Version = "2.0"

type (
	// Swagger object
	Swagger struct {
		Version             string                            `json:"swagger"`
		Info                *Info                             `json:"info"`
		Host                string                            `json:"host"`
		BasePath            string                            `json:"basePath"`
		Tags                []*Tag                            `json:"tags"`
		Schemes             []string                          `json:"schemes"`
		Paths               map[string]map[string]*Opera      `json:"paths,omitempty"` // {"prefix":{"method":{...}}}
		SecurityDefinitions map[string]map[string]interface{} `json:"securityDefinitions,omitempty"`
		Definitions         map[string]*Definition            `json:"definitions,omitempty"`
		ExternalDocs        map[string]string                 `json:"externalDocs,omitempty"`
	}
	// Info object
	Info struct {
		Title          string   `json:"title"`
		ApiVersion     string   `json:"version"`
		Description    string   `json:"description"`
		Contact        *Contact `json:"contact"`
		TermsOfService string   `json:"termsOfService"`
		License        *License `json:"license,omitempty"`
	}
	// Contact object
	Contact struct {
		Email string `json:"email,omitempty"`
	}
	// License object
	License struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	// Tag object
	Tag struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	// Opera object
	Opera struct {
		Tags        []string              `json:"tags"`
		Summary     string                `json:"summary"`
		Description string                `json:"description"`
		OperationId string                `json:"operationId"`
		Consumes    []string              `json:"consumes,omitempty"`
		Produces    []string              `json:"produces,omitempty"`
		Parameters  []*Parameter          `json:"parameters,omitempty"`
		Responses   map[string]*Resp      `json:"responses,omitempty"` // {"httpcode":resp}
		Security    []map[string][]string `json:"security,omitempty"`
	}
	// Parameter object
	Parameter struct {
		In               string      `json:"in"` // the position of the parameter
		Name             string      `json:"name"`
		Description      string      `json:"description"`
		Required         bool        `json:"required"`
		Type             string      `json:"type,omitempty"` // "array"|"integer"|"object"
		Items            *Items      `json:"items,omitempty"`
		Schema           *Schema     `json:"schema,omitempty"`
		CollectionFormat string      `json:"collectionFormat,omitempty"` // "multi"
		Format           string      `json:"format,omitempty"`           // "int64"
		Default          interface{} `json:"default,omitempty"`
	}
	// Items object
	Items struct {
		Ref     string      `json:"$ref,omitempty"`
		Type    string      `json:"type"`           // "string"
		Enum    interface{} `json:"enum,omitempty"` // slice
		Default interface{} `json:"default,omitempty"`
	}
	// Schema object
	Schema struct {
		Ref                  string            `json:"$ref,omitempty"`
		Type                 string            `json:"type,omitempty"` // "array"|"integer"|"object"
		Items                *Items            `json:"items,omitempty"`
		Description          string            `json:"description,omitempty"`
		AdditionalProperties map[string]string `json:"additionalProperties,omitempty"`
	}
	// Resp object
	Resp struct {
		Schema      *Schema `json:"schema,omitempty"`
		Description string  `json:"description,omitempty"`
	}
	// Definition object
	Definition struct {
		Type       string               `json:"type,omitempty"` // "object"
		Properties map[string]*Property `json:"properties,omitempty"`
		Xml        *Xml                 `json:"xml,omitempty"`
	}
	// Property object
	Property struct {
		Type        string      `json:"type,omitempty"`   // "array"|"integer"|"object"
		Format      string      `json:"format,omitempty"` // "int64"
		Description string      `json:"description,omitempty"`
		Enum        []string    `json:"enum,omitempty"`
		Example     interface{} `json:"example,omitempty"`
		Default     interface{} `json:"default,omitempty"`
	}
	// Xml object
	Xml struct {
		Name    string `json:"name"`
		Wrapped bool   `json:"wrapped,omitempty"`
	}
)

// CommonMIMETypes common MIME types
var CommonMIMETypes = []string{
	"application/json",
	"application/javascript",
	"application/xml",
	"application/x-www-form-urlencoded",
	"application/protobuf",
	"application/msgpack",
	"text/html",
	"text/plain",
	"multipart/form-data",
	"application/octet-stream",
}

// github.com/mcuadros/go-jsonschema-generator
var mapping = map[reflect.Kind]string{
	reflect.Bool:    "bool",
	reflect.Int:     "integer",
	reflect.Int8:    "integer",
	reflect.Int16:   "integer",
	reflect.Int32:   "integer",
	reflect.Int64:   "integer",
	reflect.Uint:    "integer",
	reflect.Uint8:   "integer",
	reflect.Uint16:  "integer",
	reflect.Uint32:  "integer",
	reflect.Uint64:  "integer",
	reflect.Float32: "number",
	reflect.Float64: "number",
	reflect.String:  "string",
	reflect.Slice:   "array",
	reflect.Struct:  "object",
	reflect.Map:     "object",
}

var mapping2 = map[string]string{
	"bool":    "bool",
	"int":     "integer",
	"int8":    "integer",
	"int16":   "integer",
	"int32":   "integer",
	"int64":   "integer",
	"uint":    "integer",
	"uint8":   "integer",
	"uint16":  "integer",
	"uint32":  "integer",
	"uint64":  "integer",
	"float32": "number",
	"float64": "number",
	"string":  "string",
}

// SliceInfo slice parameter information
func SliceInfo(value interface{}) (subtyp string, first interface{}, count int) {
	subtyp = fmt.Sprintf("%T", value)
	idx := strings.Index(subtyp, "]")
	subtyp = subtyp[idx+1:]
	if strings.HasPrefix(subtyp, "[]") {
		subtyp = "array"
	} else {
		subtyp = mapping2[subtyp]
		if len(subtyp) == 0 {
			subtyp = "object"
		}
	}
	rv := reflect.Indirect(reflect.ValueOf(value))
	count = rv.Len()
	if count > 0 {
		first = rv.Index(0).Interface()
	} else {
		first = reflect.New(rv.Type().Elem()).Elem().Interface()
	}
	return
}

// ParamType type of the parameter value passed in
func ParamType(value interface{}) string {
	if value == nil {
		return ""
	}
	rv, ok := value.(reflect.Type)
	if !ok {
		rv = reflect.TypeOf(value)
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	tn := rv.String()
	if tn == "multipart.FileHeader" || tn == "[]*multipart.FileHeader" || tn == "[]multipart.FileHeader" {
		return "file"
	}
	return mapping[rv.Kind()]
}

// CreateProperties creates properties
func CreateProperties(obj interface{}) map[string]*Property {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for {
		if t.Kind() != reflect.Ptr {
			break
		}
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if v.Len() > 0 {
			v = v.Index(0)
		} else {
			v = reflect.Value{}
		}
	}
	for {
		if t.Kind() != reflect.Ptr {
			break
		}
		t = t.Elem()
	}
	for {
		if v.Kind() != reflect.Ptr {
			break
		}
		v = v.Elem()
	}
	if v == (reflect.Value{}) {
		v = reflect.New(t).Elem()
	}

	ps := map[string]*Property{}
	switch t.Kind() {
	case reflect.Map:
		kvs := v.MapKeys()
		for _, kv := range kvs {
			val := v.MapIndex(kv)
			for {
				if val.Kind() != reflect.Ptr {
					break
				}
				val = val.Elem()
			}
			if val == (reflect.Value{}) {
				val = reflect.New(val.Type()).Elem()
			}
			p := &Property{
				Type:    ParamType(val.Type()),
				Format:  val.Type().Name(),
				Default: val.Interface(),
			}
			ps[kv.String()] = p
		}
		return ps

	case reflect.Struct:
		num := t.NumField()
		for i := 0; i < num; i++ {
			field := t.Field(i)
			p := &Property{
				Type:   ParamType(field.Type),
				Format: field.Type.Name(),
			}
			ft := field.Type
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr {
				fv = fv.Elem()
				ft = ft.Elem()
			}
			if !fv.CanInterface() {
				continue
			}
			if fv.Interface() == nil {
				fv = reflect.New(ft).Elem()
			}
			p.Default = fv.Interface()
			ps[field.Name] = p
		}
		return ps

	}

	return nil
}

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
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// some define
const (
	TAG_PARAM        = "param"    // request param tag name
	TAG_IGNORE_PARAM = "-"        // ignore request param tag value
	KEY_IN           = "in"       // position of param
	KEY_NAME         = "name"     // specify request param`s name
	KEY_REQUIRED     = "required" // request param is required or not
	KEY_DESC         = "desc"     // request param description
	KEY_LEN          = "len"      // length range of param's value
	KEY_RANGE        = "range"    // numerical range of param's value
	KEY_NONZERO      = "nonzero"  // param`s value can not be zero
	KEY_REGEXP       = "regexp"   // verify the value of the param with a regular expression(param value can not be null)
	KEY_MAXMB        = "maxmb"    // when request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
	KEY_ERR          = "err"      // the custom error for binding or validating

	MB                 = 1 << 20 // 1MB
	defaultMaxMemory   = 32 * MB // 32 MB
	defaultMaxMemoryMB = 32
)

// ParseTags returns the key-value in the tag string.
// If the tag does not have the conventional format,
// the value returned by ParseTags is unspecified.
func ParseTags(tag string) map[string]string {
	var values = map[string]string{}

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] != '<' {
			i++
		}
		if i >= len(tag) || tag[i] != '<' {
			break
		}
		i++

		// Skip the left Spaces
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		if i >= len(tag) {
			break
		}

		tag = tag[i:]
		if tag == "" {
			break
		}

		var name, value string
		var hadName bool
		i = 0
	PAIR:
		for i < len(tag) {
			switch tag[i] {
			case ':':
				if hadName {
					i++
					continue
				}
				name = strings.TrimRight(tag[:i], " ")
				tag = strings.TrimLeft(tag[i+1:], " ")
				hadName = true
				i = 0
			case '\\':
				i++
				// Fix the escape character of `\\<` or `\\>`
				if tag[i] == '<' || tag[i] == '>' {
					tag = tag[:i-1] + tag[i:]
				} else {
					i++
				}
			case '>':
				if !hadName {
					name = strings.TrimRight(tag[:i], " ")
				} else {
					value = strings.TrimRight(tag[:i], " ")
				}
				values[name] = value
				break PAIR
			default:
				i++
			}
		}
		if i >= len(tag) {
			break
		}
		tag = tag[i+1:]
	}
	return values
}

// func ParseTags(s string) map[string]string {
// 	c := strings.Split(s, ",")
// 	m := make(map[string]string)
// 	for _, v := range c {
// 		a := strings.IndexByte(v, '(')
// 		b := strings.LastIndexByte(v, ')')
// 		if a != -1 && b != -1 {
// 			m[v[:a]] = v[a+1 : b]
// 			continue
// 		}
// 		m[v] = ""
// 	}
// 	return m
// }

// ParseTags returns the key-value in the tag string.
// If the tag does not have the conventional format,
// the value returned by ParseTags is unspecified.
// func ParseTags(tag string) map[string]string {
// 	var values = map[string]string{}

// 	for tag != "" {
// 		// Skip leading space.
// 		i := 0
// 		for i < len(tag) && (tag[i] == ' ' || tag[i] == ',') {
// 			i++
// 		}
// 		tag = tag[i:]
// 		if tag == "" {
// 			break
// 		}

// 		// Scan to colon. A space, a quote or a control character is a syntax error.
// 		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
// 		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
// 		// as it is simpler to inspect the tag's bytes than the tag's runes.
// 		i = 0
// 		for i < len(tag) && tag[i] > ' ' && tag[i] != '<' && tag[i] != ',' && tag[i] != 0x7f {
// 			i++
// 		}
// 		if i == 0 || i+1 >= len(tag) || (tag[i] != '<' && tag[i] != ' ' && tag[i] != ',') {
// 			break
// 		}
// 		name := string(tag[:i])
// 		tag = tag[i:]
// 		if tag[0] == ' ' || tag[0] == ',' {
// 			values[name] = ""
// 			continue
// 		}
// 		// Scans a string in parentheses to find the value..
// 		i = 1
// 		for i < len(tag) && tag[i] != '>' {
// 			if tag[i] == '\\' {
// 				i++
// 				// Remove the escape character of `(` or `)`
// 				if tag[i] == '<' || tag[i] == '>' {
// 					tag = tag[:i-1] + tag[i:]
// 					continue
// 				}
// 			}
// 			i++
// 		}
// 		if i >= len(tag) {
// 			break
// 		}
// 		values[name] = string(tag[1:i])
// 		tag = tag[i+1:]
// 	}
// 	return values
// }

// Param use the struct field to define a request parameter model
type Param struct {
	apiName    string // ParamsAPI name
	name       string // param name
	indexPath  []int
	isRequired bool              // file is required or not
	isFile     bool              // is file param or not
	tags       map[string]string // struct tags for this param
	rawTag     reflect.StructTag // the raw tag
	rawValue   reflect.Value     // the raw tag value
	err        error             // the custom error for binding or validating
}

const (
	fileTypeString    = "*multipart.FileHeader"
	fileTypeString2   = "multipart.FileHeader"
	filesTypeString   = "[]*multipart.FileHeader"
	filesTypeString2  = "[]multipart.FileHeader"
	cookieTypeString  = "*http.Cookie"
	cookieTypeString2 = "http.Cookie"
	// fasthttpCookieTypeString = "fasthttp.Cookie"
	stringTypeString = "string"
	bytesTypeString  = "[]byte"
	bytes2TypeString = "[]uint8"
)

var (
	// TagInValues is values for tag 'in'
	TagInValues = map[string]bool{
		"path":     true,
		"query":    true,
		"formData": true,
		"body":     true,
		"header":   true,
		"cookie":   true,
	}
)

// Raw gets the param's original value
func (param *Param) Raw() interface{} {
	return param.rawValue.Interface()
}

// APIName gets ParamsAPI name
func (param *Param) APIName() string {
	return param.apiName
}

// Name gets parameter field name
func (param *Param) Name() string {
	return param.name
}

// In get the type value for the param
func (param *Param) In() string {
	return param.tags[KEY_IN]
}

// IsRequired tests if the param is declared
func (param *Param) IsRequired() bool {
	return param.isRequired
}

// Description gets the description value for the param
func (param *Param) Description() string {
	return param.tags[KEY_DESC]
}

// IsFile tests if the param is type *multipart.FileHeader
func (param *Param) IsFile() bool {
	return param.isFile
}

func (param *Param) validate(value reflect.Value) error {
	if value.Kind() != reflect.Slice {
		return param.validateElem(value)
	}
	var err error
	for i, count := 0, value.Len(); i < count; i++ {
		if err = param.validateElem(value.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

// Validate tests if the param conforms to it's validation constraints specified
// int the KEY_REGEXP struct tag
func (param *Param) validateElem(value reflect.Value) (err error) {
	defer func() {
		p := recover()
		if param.err != nil {
			if err != nil {
				err = param.err
			}
		} else if p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()
	// range
	if tuple, ok := param.tags[KEY_RANGE]; ok {
		var f64 float64
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f64 = float64(value.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			f64 = float64(value.Uint())
		case reflect.Float32, reflect.Float64:
			f64 = value.Float()
		}
		if err = validateRange(f64, tuple, param.name); err != nil {
			return err
		}
	}
	obj := value.Interface()
	// nonzero
	if _, ok := param.tags["nonzero"]; ok {
		if value.Kind() != reflect.Struct && obj == reflect.Zero(value.Type()).Interface() {
			return NewValidationError(ValidationErrorValueNotSet, param.name)
		}
	}
	s, isString := obj.(string)
	// length
	if tuple, ok := param.tags[KEY_LEN]; ok && isString {
		if err = validateLen(s, tuple, param.name); err != nil {
			return err
		}
	}
	// regexp
	if reg, ok := param.tags[KEY_REGEXP]; ok && isString {
		if err = validateRegexp(s, reg, param.name); err != nil {
			return err
		}
	}
	return
}

func (param *Param) myError(reason string) error {
	if param.err != nil {
		return param.err
	}
	return NewError(param.apiName, param.name, reason)
}

func parseTuple(tuple string) (string, string) {
	c := strings.Split(tuple, ":")
	var a, b string
	switch len(c) {
	case 1:
		a = c[0]
		if len(a) > 0 {
			return a, a
		}
	case 2:
		a = c[0]
		b = c[1]
		if len(a) > 0 || len(b) > 0 {
			return a, b
		}
	}
	panic("invalid validation tuple")
}

func validateLen(s, tuple, paramName string) error {
	a, b := parseTuple(tuple)
	if len(a) > 0 {
		min, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		if len(s) < min {
			return NewValidationError(ValidationErrorValueTooShort, paramName)
		}
	}
	if len(b) > 0 {
		max, err := strconv.Atoi(b)
		if err != nil {
			panic(err)
		}
		if len(s) > max {
			return NewValidationError(ValidationErrorValueTooLong, paramName)
		}
	}
	return nil
}

const accuracy = 0.0000001

func validateRange(f64 float64, tuple, paramName string) error {
	a, b := parseTuple(tuple)
	if len(a) > 0 {
		min, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return err
		}
		if math.Min(f64, min) == f64 && math.Abs(f64-min) > accuracy {
			return NewValidationError(ValidationErrorValueTooSmall, paramName)
		}
	}
	if len(b) > 0 {
		max, err := strconv.ParseFloat(b, 64)
		if err != nil {
			return err
		}
		if math.Max(f64, max) == f64 && math.Abs(f64-max) > accuracy {
			return NewValidationError(ValidationErrorValueTooBig, paramName)
		}
	}
	return nil
}

func validateRegexp(s, reg, paramName string) error {
	matched, err := regexp.MatchString(reg, s)
	if err != nil {
		return err
	}
	if !matched {
		return NewValidationError(ValidationErrorValueNotMatch, paramName)
	}
	return nil
}

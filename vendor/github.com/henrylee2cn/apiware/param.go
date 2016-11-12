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
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// use the struct field to define a request parameter model
type Param struct {
	indexPath  []int
	name       string            // param name
	isRequired bool              // file is required or not
	isFile     bool              // is file param or not
	tags       map[string]string // struct tags for this param
	rawTag     reflect.StructTag // the raw tag
	rawValue   reflect.Value     // the raw tag value
}

const (
	fileTypeString           = "multipart.FileHeader"
	cookieTypeString         = "http.Cookie"
	fasthttpCookieTypeString = "fasthttp.Cookie"
	stringTypeString         = "string"
	bytesTypeString          = "[]byte"
	bytes2TypeString         = "[]uint8"
)

var (
	// values for tag 'in'
	TagInValues = map[string]bool{
		"path":     true,
		"query":    true,
		"formData": true,
		"body":     true,
		"header":   true,
		"cookie":   true,
	}
)

// Get the param's original value
func (param *Param) Raw() interface{} {
	return param.rawValue.Interface()
}

// Get parameter field name
func (param *Param) Name() string {
	return param.name
}

// Get the type value for the param
func (param *Param) In() string {
	return param.tags["in"]
}

// IsRequired tests if the param is declared
func (param *Param) IsRequired() bool {
	return param.isRequired
}

// Get the description value for the param
func (param *Param) Description() string {
	return param.tags["desc"]
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
// int the TAG_REGEXP struct tag
func (param *Param) validateElem(value reflect.Value) (err error) {
	defer func() {
		p := recover()
		if errStr, ok := param.tags[TAG_ERR]; ok {
			if err != nil {
				err = errors.New(errStr)
			}
		} else if p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()
	// range
	if tuple, ok := param.tags["range"]; ok {
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
	if tuple, ok := param.tags["len"]; ok && isString {
		if err = validateLen(s, tuple, param.name); err != nil {
			return err
		}
	}
	// regexp
	if reg, ok := param.tags[TAG_REGEXP]; ok && isString {
		if err = validateRegexp(s, reg, param.name); err != nil {
			return err
		}
	}
	return
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

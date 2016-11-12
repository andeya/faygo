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

const (
	ValidationErrorValueNotSet = (1<<16 + iota)
	ValidationErrorValueTooSmall
	ValidationErrorValueTooBig
	ValidationErrorValueTooShort
	ValidationErrorValueTooLong
	ValidationErrorValueNotMatch
)

// Validation error type
type ValidationError struct {
	kind  int
	field string
}

// NewValidationError returns a new validation error with the specified id and
// text. The id's purpose is to distinguish different validation error types.
// Built-in validation error ids start at 65536, so you should keep your custom
// ids under that value.
func NewValidationError(id int, field string) error {
	return &ValidationError{id, field}
}

func (e *ValidationError) Error() string {
	kindStr := ""
	switch e.kind {
	case ValidationErrorValueNotSet:
		kindStr = " not set"
	case ValidationErrorValueTooBig:
		kindStr = " too big"
	case ValidationErrorValueTooLong:
		kindStr = " too long"
	case ValidationErrorValueTooSmall:
		kindStr = " too small"
	case ValidationErrorValueTooShort:
		kindStr = " too short"
	case ValidationErrorValueNotMatch:
		kindStr = " not match"
	}
	return e.field + kindStr
}

func (e *ValidationError) Kind() int {
	return e.kind
}

func (e *ValidationError) Field() string {
	return e.field
}

type Error struct {
	Api    string `json:"api"`
	Param  string `json:"param"`
	Reason string `json:"reason"`
}

func NewError(api string, param string, reason string) *Error {
	return &Error{
		Api:    api,
		Param:  param,
		Reason: reason,
	}
}

var _ error = new(Error)

func (e *Error) Error() string {
	return "[apiware] " + e.Api + " | " + e.Param + " | " + e.Reason
}

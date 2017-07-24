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

// Error a formatted error type
type Error struct {
	Api    string `json:"api"`
	Param  string `json:"param"`
	Reason string `json:"reason"`
}

// NewError creates *Error
func NewError(api string, param string, reason string) *Error {
	return &Error{
		Api:    api,
		Param:  param,
		Reason: reason,
	}
}

var _ error = new(Error)

// Error implements error interface
func (e *Error) Error() string {
	return "[apiware] " + e.Api + " | " + e.Param + " | " + e.Reason
}

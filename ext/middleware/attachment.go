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

// define common middlewares.

package middleware

import (
	"regexp"
	"strings"

	"github.com/henrylee2cn/faygo"
)

// NewAttachment has the response content downloaded as an attachment file.
// Note: if the specifiedFileExtension is empty, applies to any response content.
func NewAttachment(specifiedFileExtension ...string) faygo.HandlerFunc {
	re := regexp.MustCompile("^\\.[^\\.]+$")
	hash := make(map[string]bool, len(specifiedFileExtension))
	for _, s := range specifiedFileExtension {
		if !re.MatchString(s) {
			faygo.Fatalf("Invalid file extension: %s", s)
		}
		hash[strings.ToLower(s)] = true
	}
	return func(ctx *faygo.Context) error {
		var isAttachment bool
		if len(specifiedFileExtension) == 0 {
			isAttachment = true
		} else {
			p := ctx.Path()
			if idx := strings.LastIndex(p, "."); idx != -1 {
				isAttachment = hash[strings.ToLower(p[idx:])]
			}
		}
		if isAttachment {
			ctx.SetHeader(faygo.HeaderContentDisposition, "attachment")
		}
		return nil
	}
}

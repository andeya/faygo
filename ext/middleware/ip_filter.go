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
	"net/http"
	"strings"

	"github.com/henrylee2cn/faygo"
)

// NewIPFilter creates middleware that intercepts the specified IP prefix.
func NewIPFilter(whitelist []string, realIP bool) faygo.HandlerFunc {
	var noAccess bool
	var match []string
	var prefix []string

	if len(whitelist) == 0 {
		noAccess = true
	} else {
		for _, s := range whitelist {
			if strings.HasSuffix(s, "*") {
				prefix = append(prefix, s[:len(s)-1])
			} else {
				match = append(match, s)
			}
		}
	}

	return func(ctx *faygo.Context) error {
		if noAccess {
			ctx.Error(http.StatusForbidden, "no access")
			return nil
		}

		var ip string
		if realIP {
			ip = ctx.RealIP()
		} else {
			ip = ctx.IP()
		}
		for _, ipMatch := range match {
			if ipMatch == ip {
				ctx.Next()
				return nil
			}
		}
		for _, ipPrefix := range prefix {
			if strings.HasPrefix(ip, ipPrefix) {
				ctx.Next()
				return nil
			}
		}
		ctx.Error(http.StatusForbidden, "not allow to access: "+ip)
		return nil
	}
}

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

package thinkgo

import (
	"net/http"
	"strings"
	"time"

	"github.com/henrylee2cn/thinkgo/logging/color"
)

// Create middleware that intercepts the specified IP prefix.
func NewIPFilter(prefixList []string, realIP bool) HandlerFunc {
	return func(ctx *Context) error {
		if len(prefixList) == 0 {
			ctx.Error(http.StatusForbidden, "no access")
			return nil
		}

		var ip string
		if realIP {
			ip = ctx.RealIP()
		} else {
			ip = ctx.IP()
		}
		for _, ipPrefix := range prefixList {
			if strings.HasPrefix(ip, ipPrefix) {
				ctx.Next()
				return nil
			}
		}
		ctx.Error(http.StatusForbidden, "not allow to access: "+ip)
		return nil
	}
}

// Cross-Domain middleware
func CrossDomainFilter() HandlerFunc {
	return func(ctx *Context) error {
		ctx.SetHeader(HeaderAccessControlAllowOrigin, ctx.HeaderParam(HeaderOrigin))
		// ctx.SetHeader(HeaderAccessControlAllowOrigin, "*")
		ctx.SetHeader(HeaderAccessControlAllowCredentials, "true")
		return nil
	}
}

// Access log statistics
func AccessLogWare() HandlerFunc {
	return func(ctx *Context) error {
		var u = ctx.URI()
		start := time.Now()
		ctx.Next()
		stop := time.Now()

		method := ctx.Method()
		if u == "" {
			u = "/"
		}

		n := ctx.W.Status()
		code := color.Green(n)
		switch {
		case n >= 500:
			code = color.Red(n)
		case n >= 400:
			code = color.Magenta(n)
		case n >= 300:
			code = color.Grey(n)
		}

		ctx.Log().Infof("%15s %7s  %3s %10d %12s %-30s | ", ctx.RealIP(), method, code, ctx.W.Size(), stop.Sub(start), u)
		return nil
	}
}

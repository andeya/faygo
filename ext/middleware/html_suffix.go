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
	"path"
	"strings"

	"github.com/henrylee2cn/faygo"
)

//AutoHTMLSuffix adds smartly .html suffix to static route
func AutoHTMLSuffix() faygo.HandlerFunc {
	return func(c *faygo.Context) error {
		ps := c.PathParamAll()
		p, ok := ps.Get(faygo.FilepathKey)
		if !ok {
			return nil
		}
		if p[len(p)-1] != '/' {
			ext := path.Ext(p)
			if ext == "" || ext[0] != '.' {
				newFilepath := strings.TrimSuffix(p, ext) + ".html" + ext
				ps.Replace(faygo.FilepathKey, newFilepath, 1)
			}
		}
		return nil
	}
}

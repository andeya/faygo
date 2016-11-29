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

// API automation documentation.

package thinkgo

import (
	"fmt"
	"path"
	"strings"

	"github.com/henrylee2cn/thinkgo/swagger"
)

// register the API doc router.
func (frame *Framework) regAPIdoc() {
	frame.config.APIdoc.Path = path.Join("/", frame.config.APIdoc.Path)
	// jsonPattern := strings.TrimRight(strings.Replace(frame.config.APIdoc.Path, "/*filepath", "", -1), "/") + ".json"
	jsonPattern := "/swagger.json"
	if frame.config.APIdoc.NoLimit {
		frame.MuxAPI.NamedStaticFS("APIdoc-Swagger", frame.config.APIdoc.Path, swagger.AssetFS())
		frame.MuxAPI.NamedGET("APIdoc-Swagger-JSON", jsonPattern, newAPIdocJSONHandler())
	} else {
		allowApidoc := NewIPFilter(frame.config.APIdoc.Whitelist, frame.config.APIdoc.RealIP)
		frame.MuxAPI.NamedStaticFS("APIdoc-Swagger", frame.config.APIdoc.Path, swagger.AssetFS()).Use(allowApidoc)
		frame.MuxAPI.NamedGET("APIdoc-Swagger-JSON", jsonPattern, newAPIdocJSONHandler(), allowApidoc)
	}

	tip := `APIdoc's URL path is '` + frame.config.APIdoc.Path
	if frame.config.APIdoc.NoLimit {
		frame.syslog.Criticalf(tip + `' [free access]`)
	} else if len(frame.config.APIdoc.Whitelist) == 0 {
		frame.syslog.Criticalf(tip + `' [no access]`)
	} else if frame.config.APIdoc.RealIP {
		frame.syslog.Criticalf(tip + `' [check real ip for filter]`)
	} else {
		frame.syslog.Criticalf(tip + `' [check direct ip for filter]`)
	}
}

func newAPIdocJSONHandler() HandlerFunc {
	return func(ctx *Context) error {
		if ctx.frame.apidoc == nil {
			ctx.frame.initAPIdoc(ctx.R.Host)
		}
		ctx.frame.apidoc.Schemes = []string{ctx.Scheme()}
		ctx.frame.apidoc.Host = ctx.R.Host
		return ctx.JSON(200, ctx.frame.apidoc)
	}
}

var apiStaticParam = &swagger.Parameter{
	In:          "path",
	Name:        "filepath",
	Type:        swagger.ParamType("*"),
	Description: "any static path or file",
	Required:    true,
	Format:      fmt.Sprintf("%T", "*"),
	Default:     "",
}

func (frame *Framework) initAPIdoc(host string) {
	rootMuxAPI := frame.MuxAPI
	rootTag := &swagger.Tag{
		Name:        rootMuxAPI.Path(),
		Description: apiTagDesc(rootMuxAPI.Name()),
	}
	frame.apidoc = &swagger.Swagger{
		Version: swagger.Version,
		Info: &swagger.Info{
			Title:          strings.Title(frame.Name()) + " API",
			ApiVersion:     frame.Version(),
			Description:    frame.config.APIdoc.Desc,
			Contact:        &swagger.Contact{Email: frame.config.APIdoc.Email},
			TermsOfService: frame.config.APIdoc.TermsURL,
			License: &swagger.License{
				Name: frame.config.APIdoc.License,
				Url:  frame.config.APIdoc.LicenseURL,
			},
		},
		Host:     host,
		BasePath: "/",
		Tags:     []*swagger.Tag{rootTag},
		Schemes:  []string{"http", "https"},
		Paths:    map[string]map[string]*swagger.Opera{},
		// SecurityDefinitions: map[string]map[string]interface{}{},
		// Definitions:         map[string]Definition{},
		// ExternalDocs:        map[string]string{},
	}

	for _, child := range rootMuxAPI.Children() {
		if !child.IsGroup() {
			addpath(child, rootTag)
			continue
		}
		childTag := &swagger.Tag{
			Name:        child.Path(),
			Description: apiTagDesc(child.Name()),
		}
		frame.apidoc.Tags = append(frame.apidoc.Tags, childTag)
		for _, grandson := range child.Children() {
			if !grandson.IsGroup() {
				addpath(grandson, childTag)
				continue
			}
			grandsonTag := &swagger.Tag{
				Name:        grandson.Path(),
				Description: apiTagDesc(grandson.Name()),
			}
			frame.apidoc.Tags = append(frame.apidoc.Tags, grandsonTag)
			for _, progeny := range grandson.Progeny() {
				if !progeny.IsGroup() {
					addpath(progeny, grandsonTag)
					continue
				}
			}
		}
	}
}

// 添加API操作项
func addpath(mux *MuxAPI, tag *swagger.Tag) {
	operas := map[string]*swagger.Opera{}
	pid := apiCreatePath(mux.Path())
	summary := apiSummary(mux.Name())
	desc := apiDesc(mux.Name())
	for _, method := range mux.Methods() {
		if method == "CONNECT" || method == "TRACE" {
			continue
		}
		if method == "WS" {
			method = "GET"
		}
		o := &swagger.Opera{
			Tags:        []string{tag.Name},
			Summary:     summary,
			Description: desc,
			OperationId: pid + "-" + method,
			Consumes:    swagger.CommonMIMETypes,
			Produces:    swagger.CommonMIMETypes,
			Responses:   make(map[string]*swagger.Resp, 1),
			// Security: []map[string][]string{},
		}

		for _, param := range mux.ParamInfos() {
			p := &swagger.Parameter{
				In:          param.In,
				Name:        param.Name,
				Description: param.Desc,
				Required:    param.Required,
				// Items:       &Items{},
				// Schema:      &Schema{},
			}
			typ := swagger.ParamType(param.Model)
			switch p.In {
			case "cookie":
				continue
			default:
				switch typ {
				case "file":
					o.Consumes = []string{"multipart/form-data"}
					p.Type = typ

				case "array":
					subtyp, first, count := swagger.SliceInfo(param.Model)
					switch subtyp {
					case "object":
						ref := apiDefinitions(mux, p.Name, method, param.Model)
						p.Schema = &swagger.Schema{
							Type: typ,
							Items: &swagger.Items{
								Ref: "#/definitions/" + ref,
							},
						}

					default:
						p.Type = typ
						p.Items = &swagger.Items{
							Type:    subtyp,
							Default: first,
						}
						if count > 0 {
							p.Items.Enum = param.Model
						}
						p.CollectionFormat = "multi"
					}

				case "object":
					ref := apiDefinitions(mux, p.Name, method, param.Model)
					p.Schema = &swagger.Schema{
						Type: typ,
						Ref:  "#/definitions/" + ref,
					}

				default:
					p.Type = typ
					p.Format = fmt.Sprintf("%T", param.Model)
					p.Default = param.Model
				}
			}

			o.Parameters = append(o.Parameters, p)
		}

		// static file
		if strings.HasSuffix(pid, "/{filepath}") {
			o.Parameters = append(o.Parameters, apiStaticParam)
		}

		// response
		for _, r := range mux.Returns() {
			var resp = &swagger.Resp{
				Description: r.Description,
				Schema: &swagger.Schema{
					Type: swagger.ParamType(r.ExampleValue),
				},
			}
			status := fmt.Sprint(r.Code)
			name := "http" + status
			switch resp.Schema.Type {
			case "array":
				subtyp, first, count := swagger.SliceInfo(r.ExampleValue)
				switch subtyp {
				case "object":
					ref := apiDefinitions(mux, name, method, r.ExampleValue)
					resp.Schema.Items = &swagger.Items{
						Ref: "#/definitions/" + ref,
					}

				default:
					resp.Schema.Items = &swagger.Items{
						Type:    subtyp,
						Default: first,
					}
					if count > 0 {
						resp.Schema.Items.Enum = r.ExampleValue
					}
				}

			case "object":
				ref := apiDefinitions(mux, name, method, r.ExampleValue)
				resp.Schema.Ref = "#/definitions/" + ref
			}
			o.Responses[status] = resp
		}
		if len(o.Responses) == 0 {
			o.Responses["200"] = &swagger.Resp{Description: "ok"}
		}
		operas[strings.ToLower(method)] = o
	}
	if _operas, ok := mux.frame.apidoc.Paths[pid]; ok {
		for k, v := range operas {
			_operas[k] = v
		}
	} else {
		mux.frame.apidoc.Paths[pid] = operas
	}
}

func apiDefinitions(mux *MuxAPI, pname, method string, format interface{}) (ref string) {
	upath := mux.Path()
	ref = strings.Replace(path.Join(upath[1:], pname, method), "/", "@", -1)
	def := &swagger.Definition{
		Type: "object",
		Xml:  &swagger.Xml{Name: ref},
	}
	def.Properties = swagger.CreateProperties(format)
	if mux.frame.apidoc.Definitions == nil {
		mux.frame.apidoc.Definitions = map[string]*swagger.Definition{}
	}
	mux.frame.apidoc.Definitions[ref] = def
	return
}

func apiCreatePath(u string) string {
	a := strings.Split(u, "/*")
	if len(a) > 1 {
		u = path.Join(a[0], "{filepath}")
	}
	s := strings.Split(u, "/:")
	p := s[0]
	if len(s) == 1 {
		return p
	}
	for _, param := range s[1:] {
		p = path.Join(p, "{"+param+"}")
	}
	return p
}

func apiTagDesc(desc string) string {
	return strings.TrimSpace(desc)
}

func apiSummary(desc string) string {
	return strings.TrimSpace(strings.Split(strings.TrimSpace(desc), "\n")[0])
}

func apiDesc(name string) string {
	return "<pre style=\"line-height:18px;\">" + name + "</pre>"
}

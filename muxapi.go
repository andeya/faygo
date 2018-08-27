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
//
//
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//  Syntax    Type
//  :name     named parameter
//  *name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//  Path: /blog/:category/:post
//
//  Requests:
//   /blog/go/request-routers            match: category="go", post="request-routers"
//   /blog/go/request-routers/           no match, but the router would redirect
//   /blog/go/                           no match
//   /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all parameters must always be the final path element.
//  Path: /files/*filepath
//
//  Requests:
//   /files/                             match: filepath="/"
//   /files/LICENSE                      match: filepath="/LICENSE"
//   /files/templates/article.html       match: filepath="/templates/article.html"
//   /files                              no match, but the router would redirect
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the Handle func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//  // by the name of the parameter
//  user := ps.ByName("user") // defined by :user or *user
//
//  // by the index of the parameter. This way you can also get the name (key)
//  thirdKey   := ps[2].Key   // the name of the 3rd parameter
//  thirdValue := ps[2].Value // the value of the 3rd parameter
//

package faygo

import (
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
)

type (
	// MuxAPI the visible api for the serveMux, in order to prepare for routing.
	MuxAPI struct {
		name       string
		pattern    string
		path       string
		methods    []string
		handlers   []Handler
		paramInfos []ParamInfo
		notes      []Notes
		parent     *MuxAPI
		children   []*MuxAPI
		frame      *Framework
	}
	// Methodset is the methods string of request
	Methodset string
)

// RESTfulMethodList is the list of all RESTful methods
var RESTfulMethodList = []string{
	"CONNECT",
	"DELETE",
	"GET",
	"HEAD",
	"OPTIONS",
	"PATCH",
	"POST",
	"PUT",
	"TRACE",
}

func newMuxAPI(frame *Framework, name string, methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	muxapi := &MuxAPI{
		name:       name,
		pattern:    pattern,
		methods:    methodset.Methods(),
		handlers:   handlers,
		paramInfos: []ParamInfo{},
		notes:      []Notes{},
		children:   []*MuxAPI{},
		frame:      frame,
	}
	return muxapi
}

// Methods parses out the list of methods.
// List of common methods:
//  CONNECT
//  DELETE
//  GET
//  HEAD
//  OPTIONS
//  PATCH
//  POST
//  PUT
//  TRACE
//  "*"——CONNECT/DELETE/GET/HEAD/OPTIONS/PATCH/POST/PUT/TRACE
func (m *Methodset) Methods() []string {
	s := strings.ToUpper(string(*m))
	if strings.Contains(s, "*") {
		methods := make([]string, len(RESTfulMethodList))
		copy(methods, RESTfulMethodList)
		return methods
	}
	methods := []string{}
	for _, method := range RESTfulMethodList {
		if strings.Contains(s, method) {
			methods = append(methods, method)
		}
	}
	return methods
}

// HasMethod checks whether the specified method exists or not.
func (mux *MuxAPI) HasMethod(method string) bool {
	method = strings.ToUpper(method)
	for _, m := range mux.methods {
		if method == m {
			return true
		}
	}
	return false
}

// Group adds a subordinate subgroup node to the current muxAPI grouping node.
// notes: handler cannot be nil.
func (mux *MuxAPI) Group(pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI("", "", pattern, handlers...)
}

// NamedGroup adds a subordinate subgroup node with the name to the current muxAPI grouping node.
// notes: handler cannot be nil.
func (mux *MuxAPI) NamedGroup(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "", pattern, handlers...)
}

// IsGroup returns whether the muxapi node is group or not.
func (mux *MuxAPI) IsGroup() bool {
	return len(mux.methods) == 0
}

// API adds a subordinate node to the current muxAPI grouping node.
// notes: handler cannot be nil.
func (mux *MuxAPI) API(methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI("", methodset, pattern, handlers...)
}

// NamedAPI adds a subordinate node with the name to the current muxAPI grouping node.
// notes: handler cannot be nil.
func (mux *MuxAPI) NamedAPI(name string, methodset Methodset, pattern string, handlers ...Handler) *MuxAPI {
	for _, h := range handlers {
		if h == nil {
			errStr := "handler cannot be nil:" + reflect.TypeOf(h).String()
			mux.frame.Log().Panicf("%s\n", errStr)
		}
	}
	pattern = path.Join("/", pattern)
	var child = newMuxAPI(mux.frame, name, methodset, pattern, handlers...)
	mux.children = append(mux.children, child)
	child.parent = mux
	return child
}

// GET is a shortcut for muxAPI.API("GET", pattern, handlers...)
func (mux *MuxAPI) GET(pattern string, handlers ...Handler) *MuxAPI {
	return mux.API("GET", pattern, handlers...)
}

// HEAD is a shortcut for muxAPI.API("HEAD", pattern, handlers...)
func (mux *MuxAPI) HEAD(pattern string, handlers ...Handler) *MuxAPI {
	return mux.API("HEAD", pattern, handlers...)
}

// OPTIONS is a shortcut for muxAPI.API("OPTIONS", pattern, handlers...)
func (mux *MuxAPI) OPTIONS(pattern string, handlers ...Handler) *MuxAPI {
	return mux.API("OPTIONS", pattern, handlers...)
}

// POST is a shortcut for muxAPI.API("POST", pattern, handlers...)
func (mux *MuxAPI) POST(pattern string, handlers ...Handler) *MuxAPI {
	return mux.API("POST", pattern, handlers...)
}

// PUT is a shortcut for muxAPI.API("PUT", pattern, handlers...)
func (mux *MuxAPI) PUT(pattern string, handlers ...Handler) *MuxAPI {
	return mux.API("PUT", pattern, handlers...)
}

// PATCH is a shortcut for muxAPI.API("PATCH", pattern, handlers...)
func (mux *MuxAPI) PATCH(pattern string, handlers ...Handler) *MuxAPI {
	return mux.API("PATCH", pattern, handlers...)
}

// DELETE is a shortcut for muxAPI.API("DELETE", pattern, handlers...)
func (mux *MuxAPI) DELETE(pattern string, handlers ...Handler) *MuxAPI {
	return mux.API("DELETE", pattern, handlers...)
}

// NamedGET is a shortcut for muxAPI.NamedAPI(name, "GET", pattern, handlers...)
func (mux *MuxAPI) NamedGET(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "GET", pattern, handlers...)
}

// NamedHEAD is a shortcut for muxAPI.NamedAPI(name, "HEAD", pattern, handlers...)
func (mux *MuxAPI) NamedHEAD(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "HEAD", pattern, handlers...)
}

// NamedOPTIONS is a shortcut for muxAPI.NamedAPI(name, "OPTIONS", pattern, handlers...)
func (mux *MuxAPI) NamedOPTIONS(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "OPTIONS", pattern, handlers...)
}

// NamedPOST is a shortcut for muxAPI.NamedAPI(name, "POST", pattern, handlers...)
func (mux *MuxAPI) NamedPOST(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "POST", pattern, handlers...)
}

// NamedPUT is a shortcut for muxAPI.NamedAPI(name, "PUT", pattern, handlers...)
func (mux *MuxAPI) NamedPUT(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "PUT", pattern, handlers...)
}

// NamedPATCH is a shortcut for muxAPI.NamedAPI(name, "PATCH", pattern, handlers...)
func (mux *MuxAPI) NamedPATCH(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "PATCH", pattern, handlers...)
}

// NamedDELETE is a shortcut for muxAPI.NamedAPI(name, "DELETE", pattern, handlers...)
func (mux *MuxAPI) NamedDELETE(name string, pattern string, handlers ...Handler) *MuxAPI {
	return mux.NamedAPI(name, "DELETE", pattern, handlers...)
}

// FilepathKey path key for static router pattern.
const FilepathKey = "filepath"

// NamedStaticFS serves files from the given file system fs.
// The pattern must end with "/*filepath", files are then served from the local
// pattern /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     frame.StaticFS("/src/*filepath", Dir("/var/www", true, true)
func (mux *MuxAPI) NamedStaticFS(name, pattern string, fs FileSystem) *MuxAPI {
	if fs == nil {
		errStr := "For file server, fs (http.FileSystem) cannot be nil"
		mux.frame.Log().Panicf("%s\n", errStr)
	}
	if len(pattern) < 10 || pattern[len(pattern)-10:] != "/*"+FilepathKey {
		pattern = path.Join(pattern, "/*"+FilepathKey)
	}
	handler := func(fileServer Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			ctx.R.URL.Path = ctx.pathParams.ByName(FilepathKey)
			return fileServer.Serve(ctx)
		})
	}(global.fsManager.FileServer(fs))
	return mux.NamedAPI(name, "GET", pattern, handler)
}

// StaticFS is similar to NamedStaticFS, but no name.
func (mux *MuxAPI) StaticFS(pattern string, fs FileSystem) *MuxAPI {
	return mux.NamedStaticFS("fileserver", pattern, fs)
}

// NamedStatic is similar to NamedStaticFS, but the second parameter is the local file path.
func (mux *MuxAPI) NamedStatic(name, pattern string, root string, nocompressAndNocache ...bool) *MuxAPI {
	os.MkdirAll(root, 0777)
	return mux.NamedStaticFS(name, pattern, DirFS(root, nocompressAndNocache...))
}

// Static is similar to NamedStatic, but no name.
func (mux *MuxAPI) Static(pattern string, root string, nocompressAndNocache ...bool) *MuxAPI {
	return mux.NamedStatic(root, pattern, root, nocompressAndNocache...)
}

// Use inserts the middlewares at the left end of the node's handler chain.
// notes: handler cannot be nil.
func (mux *MuxAPI) Use(handlers ...Handler) *MuxAPI {
	_handlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		if h == nil {
			errStr := "For using middleware, handler cannot be nil:" + reflect.TypeOf(h).String()
			mux.frame.Log().Panicf("%s\n", errStr)
		}
		if !IsHandlerWithoutPath(h, mux.frame.config.Router.NoDefaultParams) {
			errStr := "For using middleware, the handlers can not bind the path parameter:" + reflect.TypeOf(h).String()
			mux.frame.Log().Panicf("%s\n", errStr)
		}
		_handlers[i] = h
	}
	mux.handlers = append(_handlers, mux.handlers...)
	return mux
}

// comb mux.handlers, mux.paramInfos, mux.notes and mux.path,.
// sort children by path.
// note: can only be executed once before HTTP serving.
func (mux *MuxAPI) comb() {
	mux.paramInfos = mux.paramInfos[:0]
	mux.notes = mux.notes[:0]
	for i, handler := range mux.handlers {
		h, err := ToAPIHandler(handler, mux.frame.config.Router.NoDefaultParams)
		if err != nil {
			if err == ErrNotStructPtr || err == ErrNoParamHandler {
				// Get the information for apidoc
				if doc, ok := handler.(APIDoc); ok {
					docinfo := doc.Doc()
					if docinfo.Note != "" || docinfo.Return != nil {
						mux.notes = append(mux.notes, Notes{Note: docinfo.Note, Return: docinfo.Return})
					}
					for _, param := range docinfo.MoreParams {
						// The path parameter must be a required parameter.
						if param.In == "path" {
							param.Required = true
						}
						mux.paramInfos = append(mux.paramInfos, param)
					}
				}
				continue
			}
			errStr := "[Faygo-ToAPIHandler] " + err.Error()
			mux.frame.Log().Panicf("%s\n", errStr)
		}

		if h.paramsAPI.MaxMemory() == defaultMultipartMaxMemory {
			h.paramsAPI.SetMaxMemory(mux.frame.config.multipartMaxMemory)
		}
		// Get the information for apidoc
		docinfo := h.Doc()
		if docinfo.Note != "" || docinfo.Return != nil {
			mux.notes = append(mux.notes, Notes{Note: docinfo.Note, Return: docinfo.Return})
		}
		for _, param := range docinfo.MoreParams {
			// The path parameter must be a required parameter.
			if param.In == "path" {
				param.Required = true
			}
			mux.paramInfos = append(mux.paramInfos, param)
		}

		mux.handlers[i] = h
	}

	// check path params defined, and panic if there is any error.
	mux.checkPathParams()

	mux.path = mux.pattern
	if mux.parent != nil {
		mux.path = path.Join(mux.parent.path, mux.path)
		mux.notes = append(mux.parent.notes, mux.notes...)
		mux.paramInfos = append(mux.parent.paramInfos, mux.paramInfos...)
		mux.handlers = append(mux.parent.handlers, mux.handlers...)
	}

	// Get distinct and sorted parameters information.
	mux.paramInfos = distinctAndSortedParamInfos(mux.paramInfos)

	if len(mux.children) == 0 {
		// Check for body parameter conflicts
		mux.checkBodyParamConflicts()
	} else {
		for _, child := range mux.children {
			child.comb()
		}
		sort.Sort(MuxAPIs(mux.children))
	}
}

// check path params defined, and panic if there is any error.
func (mux *MuxAPI) checkPathParams() {
	var numPathParams uint8
	for _, paramInfo := range mux.paramInfos {
		if paramInfo.In != "path" {
			continue
		}
		if !strings.Contains(mux.pattern, "/:"+paramInfo.Name) && !strings.Contains(mux.pattern, "/*"+paramInfo.Name) {
			mux.frame.Log().Panicf(
				"[Faygo-checkPathParams] the router pattern `%s` does not match the path param:\n%#v",
				mux.pattern,
				paramInfo,
			)
		}
		numPathParams++
	}
	if countPathParams(mux.pattern) < numPathParams {
		mux.frame.Log().Panicf(
			"[Faygo-checkPathParams] the router pattern `%s` does not match the path params:\n%#v",
			mux.pattern,
			mux.paramInfos,
		)
	}
}

// check path params defined, and panic if there is any error.
func (mux *MuxAPI) checkBodyParamConflicts() {
	var hasBody bool
	var hasFormData bool
	for _, paramInfo := range mux.paramInfos {
		switch paramInfo.In {
		case "formData":
			if hasBody {
				errStr := "[Faygo-checkBodyParamConflicts] handler struct tags of `in(formData)` and `in(body)` can not exist at the same time:\nURL path: " + mux.path
				mux.frame.Log().Panicf("%s\n", errStr)
			}
			hasFormData = true
		case "body":
			if hasFormData {
				errStr := "[Faygo-checkBodyParamConflicts] handler struct tags of `in(formData)` and `in(body)` can not exist at the same time:\nURL path: " + mux.path
				mux.frame.Log().Panicf("%s\n", errStr)
			}
			if hasBody {
				errStr := "[Faygo-checkBodyParamConflicts] there should not be more than one handler struct tag `in(body)`:\nURL path: " + mux.path
				mux.frame.Log().Panicf("%s\n", errStr)
			}
			hasBody = true
		}
	}
}

// Methods returns the methods of muxAPI node.
func (mux *MuxAPI) Methods() []string {
	return mux.methods
}

// Path returns the path of muxAPI node.
func (mux *MuxAPI) Path() string {
	return mux.path
}

// Name returns the name of muxAPI node.
func (mux *MuxAPI) Name() string {
	return mux.name
}

// ParamInfos returns the paramInfos of muxAPI node.
func (mux *MuxAPI) ParamInfos() []ParamInfo {
	return mux.paramInfos
}

// Notes returns the notes of muxAPI node.
func (mux *MuxAPI) Notes() []Notes {
	return mux.notes
}

// Parent returns the parent of muxAPI node.
func (mux *MuxAPI) Parent() *MuxAPI {
	return mux.parent
}

// Children returns the children of muxAPI node.
func (mux *MuxAPI) Children() []*MuxAPI {
	return mux.children
}

// Progeny returns an ordered list of all subordinate nodes.
func (mux *MuxAPI) Progeny() []*MuxAPI {
	nodes := []*MuxAPI{}
	for _, child := range mux.children {
		child.family(&nodes)
	}
	return nodes
}

// Family returns an ordered list of tree nodes.
func (mux *MuxAPI) Family() []*MuxAPI {
	nodes := []*MuxAPI{mux}
	for _, child := range mux.children {
		child.family(&nodes)
	}
	return nodes
}

func (mux *MuxAPI) family(nodes *[]*MuxAPI) {
	*nodes = append(*nodes, mux)
	for _, child := range mux.children {
		child.family(nodes)
	}
}

// HandlerProgeny returns an ordered list of subordinate nodes used to register router.
func (mux *MuxAPI) HandlerProgeny() []*MuxAPI {
	if !mux.IsGroup() {
		return []*MuxAPI{mux}
	}
	nodes := []*MuxAPI{}
	for _, child := range mux.children {
		nodes = append(nodes, child.HandlerProgeny()...)
	}
	return nodes
}

// MuxAPIs is the array of muxAPIs for sorting
type MuxAPIs []*MuxAPI

// Len returns the length of muxAPIs
func (ends MuxAPIs) Len() int {
	return len(ends)
}

// Less returns the smaller muxAPI.
func (ends MuxAPIs) Less(i, j int) bool {
	return ends[i].path <= ends[j].path
}

// Swap swaps the two muxAPIs
func (ends MuxAPIs) Swap(i, j int) {
	ends[i], ends[j] = ends[j], ends[i]
}

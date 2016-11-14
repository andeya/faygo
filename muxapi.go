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

package thinkgo

import (
	"net/http"
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
		returns    Returns
		parent     *MuxAPI
		children   []*MuxAPI
		frame      *Framework
	}
	// method set for request
	Methodset string
)

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
		returns:    Returns{},
		children:   []*MuxAPI{},
		frame:      frame,
	}
	return muxapi
}

/*
 * Parse out the list of methods.
 *
 * List of common methods:
 * CONNECT
 * DELETE
 * GET
 * HEAD
 * OPTIONS
 * PATCH
 * POST
 * PUT
 * TRACE
 *
 * Special identification：
 * "REST"——CONNECT/DELETE/GET/HEAD/OPTIONS/PATCH/POST/PUT/TRACE
 * "WS"——web socket
 */
func (m *Methodset) Methods() []string {
	s := strings.ToUpper(string(*m))
	if strings.Contains(s, "REST") {
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
			mux.frame.Log().Critical(errStr)
			panic(errStr)
		}
	}
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

// NamedStaticFS serves files from the given file system fs.
// The pattern must end with "/*filepath", files are then served from the local
// pattern /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     frame.StaticFS("/src/*filepath", http.Dir("/var/www"))
func (mux *MuxAPI) NamedStaticFS(name, pattern string, fs http.FileSystem) *MuxAPI {
	if fs == nil {
		errStr := "For file server, fs (http.FileSystem) cannot be nil"
		mux.frame.Log().Critical(errStr)
		panic(errStr)
	}
	if len(pattern) < 10 || pattern[len(pattern)-10:] != "/*filepath" {
		pattern = path.Join(pattern, "/*filepath")
	}
	handler := func(fileServer Handler) Handler {
		return HandlerFunc(func(ctx *Context) error {
			ctx.R.URL.Path = ctx.pathParams.ByName("filepath")
			return fileServer.Serve(ctx)
		})
	}(mux.frame.fileServerManager.FileServer(fs))
	return mux.NamedAPI(name, "GET", pattern, handler)
}

// StaticFS is similar to NamedStaticFS, but no name.
func (mux *MuxAPI) StaticFS(pattern string, fs http.FileSystem) *MuxAPI {
	return mux.NamedStaticFS("fileserver", pattern, fs)
}

// NamedStatic is similar to NamedStaticFS, but the second parameter is the local file path.
func (mux *MuxAPI) NamedStatic(name, pattern string, root string) *MuxAPI {
	os.MkdirAll(root, 0777)
	return mux.NamedStaticFS(name, pattern, http.Dir(root))
}

// Static is similar to NamedStatic, but no name.
func (mux *MuxAPI) Static(pattern string, root string) *MuxAPI {
	return mux.NamedStatic(root, pattern, root)
}

// Insert the middlewares at the left end of the node's handler chain.
// notes: handler cannot be nil.
func (mux *MuxAPI) Use(handlers ...HandlerWithoutPath) *MuxAPI {
	_handlers := make([]Handler, len(handlers))
	for i, h := range handlers {
		if h == nil {
			errStr := "For using middleware, handler cannot be nil:" + reflect.TypeOf(h).String()
			mux.frame.Log().Critical(errStr)
			panic(errStr)
		}
		if !IsHandlerWithoutPath(h) {
			errStr := "For using middleware, the handlers can not bind the path parameter:" + reflect.TypeOf(h).String()
			mux.frame.Log().Critical(errStr)
			panic(errStr)
		}
		_handlers[i] = h
	}
	mux.handlers = append(_handlers, mux.handlers...)
	return mux
}

// comb mux.handlers, mux.paramInfos, mux.returns and mux.path,.
// sort children by path.
// note: can only be executed once before HTTP serving.
func (mux *MuxAPI) comb() {
	mux.paramInfos = mux.paramInfos[:0]
	mux.returns = mux.returns[:0]
	for i, handler := range mux.handlers {
		apiHandler := ToAPIHandler(handler)
		if apiHandler == nil {
			continue
		}
		h, err := newHandlerStruct(apiHandler, mux.frame.paramMapping)
		if err != nil {
			errStr := "[Thinkgo-newHandlerStruct] " + err.Error()
			mux.frame.Log().Critical(errStr)
			panic(errStr)
		}
		if h.paramsAPI.Number() == 0 {
			continue
		}
		if h.paramsAPI.MaxMemory() == defaultMultipartMaxMemory {
			h.paramsAPI.SetMaxMemory(mux.frame.config.multipartMaxMemory)
		}
		mux.handlers[i] = h

		// Get the information for apidoc
		mux.paramInfos = append(mux.paramInfos, h.paramInfos()...)
		mux.returns = append(mux.returns, h.returns()...)
	}

	// check path params defined, and panic if there is any error.
	mux.checkPathParams()

	mux.path = mux.pattern
	if mux.parent != nil {
		mux.path = path.Join(mux.parent.path, mux.path)
		mux.returns = append(mux.parent.returns, mux.returns...)
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
	for _, paramInfo := range mux.paramInfos {
		if paramInfo.In != "path" {
			continue
		}
		count := strings.Count(mux.pattern, "/:"+paramInfo.Name) + strings.Count(mux.pattern, "/*"+paramInfo.Name)
		if count != 1 {
			errStr := "[Thinkgo-checkPathParams] the router pattern does not match the path param:\nname: " +
				paramInfo.Name + "\ndesc:" + paramInfo.Desc
			mux.frame.Log().Critical(errStr)
			panic(errStr)
		}
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
				errStr := "[Thinkgo-checkBodyParamConflicts] handler struct tags of `in(formData)` and `in(body)` can not exist at the same time:\nURL path: " + mux.path
				mux.frame.Log().Critical(errStr)
				panic(errStr)
			}
			hasFormData = true
		case "body":
			if hasFormData {
				errStr := "[Thinkgo-checkBodyParamConflicts] handler struct tags of `in(formData)` and `in(body)` can not exist at the same time:\nURL path: " + mux.path
				mux.frame.Log().Critical(errStr)
				panic(errStr)
			}
			if hasBody {
				errStr := "[Thinkgo-checkBodyParamConflicts] there should not be more than one handler struct tag `in(body)`:\nURL path: " + mux.path
				mux.frame.Log().Critical(errStr)
				panic(errStr)
			}
			hasBody = true
		}
	}
}

func (mux *MuxAPI) Methods() []string {
	return mux.methods
}

func (mux *MuxAPI) Path() string {
	return mux.path
}

func (mux *MuxAPI) Name() string {
	return mux.name
}

func (mux *MuxAPI) ParamInfos() []ParamInfo {
	return mux.paramInfos
}

func (mux *MuxAPI) Returns() []Return {
	return mux.returns
}

func (mux *MuxAPI) Parent() *MuxAPI {
	return mux.parent
}

func (mux *MuxAPI) Children() []*MuxAPI {
	return mux.children
}

// Get an ordered list of all subordinate nodes.
func (mux *MuxAPI) Progeny() []*MuxAPI {
	nodes := []*MuxAPI{}
	for _, child := range mux.children {
		nodes = append(nodes, child.Progeny()...)
	}
	return nodes
}

// Get an ordered list of subordinate nodes used to register router.
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

type MuxAPIs []*MuxAPI

func (ends MuxAPIs) Len() int {
	return len(ends)
}

func (ends MuxAPIs) Less(i, j int) bool {
	return ends[i].path <= ends[j].path
}

func (ends MuxAPIs) Swap(i, j int) {
	ends[i], ends[j] = ends[j], ends[i]
}

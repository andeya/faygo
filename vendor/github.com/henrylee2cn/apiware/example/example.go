package main

import (
	"github.com/henrylee2cn/apiware"
	"strings"
)

var myApiware = apiware.New(pathDecodeFunc, nil, nil)

var pattern = "/test/:id"

func pathDecodeFunc(urlPath, pattern string) apiware.KV {
	idx := map[int]string{}
	for k, v := range strings.Split(pattern, "/") {
		if !strings.HasPrefix(v, ":") {
			continue
		}
		idx[k] = v[1:]
	}
	pathParams := make(map[string]string, len(idx))
	for k, v := range strings.Split(urlPath, "/") {
		name, ok := idx[k]
		if !ok {
			continue
		}
		pathParams[name] = v
	}
	return apiware.Map(pathParams)
}

func main() {
	// Check whether these structs meet the requirements of apiware, and register them
	err := myApiware.Register(
		new(httpTestApiware),
	)
	if err != nil {
		panic(err)
	}

	// http server
	httpServer(":8080")
}

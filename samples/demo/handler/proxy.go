package handler

import (
	"github.com/henrylee2cn/faygo"
)

type Search int

func (Search) Serve(ctx *faygo.Context) error {
	return ctx.ReverseProxy("https://cn.bing.com/search", false)
}

func (Search) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "reverse proxy",
		Return: nil,
		Params: []faygo.ParamInfo{
			{
				Name:     "q",
				In:       "query",
				Required: true,
				Model:    "golang",
				Desc:     "bing search",
			},
		},
	}
}

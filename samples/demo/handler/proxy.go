package handler

import (
	"github.com/henrylee2cn/thinkgo"
)

type Search int

func (Search) Serve(ctx *thinkgo.Context) error {
	return ctx.ReverseProxy("https://cn.bing.com/search", false)
}

func (Search) Doc() thinkgo.Doc {
	return thinkgo.Doc{
		Note:   "reverse proxy",
		Return: nil,
		Params: []thinkgo.ParamInfo{
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

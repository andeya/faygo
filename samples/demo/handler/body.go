package handler

import (
	"github.com/henrylee2cn/thinkgo"
)

type Body struct {
	Json map[string]interface{} `param:"<in:body>"`
}

func (b *Body) Serve(ctx *thinkgo.Context) error {
	return ctx.JSON(200, b.Json, true)
}

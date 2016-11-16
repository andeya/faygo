package handler

import (
	"github.com/henrylee2cn/thinkgo"
)

type Render struct {
	Title     string   `param:"in(query),nonzero"`
	Paragraph []string `param:"in(query),name(p),len(1:10)" regexp:"(^[\\w]*$)"`
}

func (r *Render) Serve(ctx *thinkgo.Context) error {
	return ctx.Render(200, thinkgo.JionStatic("test_render.html"), thinkgo.Map{
		"title": r.Title,
		"p":     r.Paragraph,
	})
}

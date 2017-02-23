package handler

import (
	"github.com/henrylee2cn/faygo"
	// "time"
)

type Render struct {
	Title     string   `param:"<in:query> <nonzero>"`
	Paragraph []string `param:"<in:query> <name:p> <len: 1:10> <regexp: ^[\\w]*$>"`
}

func (r *Render) Serve(ctx *faygo.Context) error {
	return ctx.Render(200, faygo.JoinStatic("render.html"), faygo.Map{
		"title": r.Title,
		"p":     r.Paragraph,
	})
}

func init() {
	faygo.RenderVar("__PUBLIC__", "/syso")
}

func Index() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		// time.Sleep(10e9)
		return ctx.Render(200, "../../_syso/index.html", faygo.Map{
			"TITLE":   "faygo",
			"VERSION": faygo.VERSION,
			"CONTENT": "Welcome To Faygo",
			"AUTHOR":  "HenryLee",
		})
	}
}

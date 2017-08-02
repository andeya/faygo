package handler

import (
	"github.com/henrylee2cn/faygo"
	// "time"
)

type Render struct {
	Title      string `param:"<in:query> <nonzero>"`
	Paragraph  string `param:"<in:query> <name:p> <len: 1:10> <regexp: ^[\\w]*$>"`
	TestHeader string `param:"<in:header><name:Test-Header> <required><test header>"`
}

func (r *Render) Serve(ctx *faygo.Context) error {
	return ctx.Render(200, faygo.JoinStatic("render.html"), faygo.Map{
		"title": r.Title,
		"p":     r.Paragraph + r.TestHeader,
	})
}

func init() {
	faygo.RenderVar("__PUBLIC__", "/syso")
}

func Index() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		// time.Sleep(2e9)
		return ctx.Render(200, "../../_syso/index.html", faygo.Map{
			"TITLE":   "faygo",
			"VERSION": faygo.VERSION,
			"CONTENT": "Welcome To Faygo",
			"AUTHOR":  "HenryLee",
		})
	}
}

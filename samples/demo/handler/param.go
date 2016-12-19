package handler

import (
	// "mime/multipart"
	"net/http"
	"sync"

	"github.com/henrylee2cn/thinkgo"
)

type Param struct {
	Id        int      `param:"in(path),required,desc(ID),range(0:10)"`
	Num       float32  `param:"in(query),name(n),range(0.1:10)"`
	Title     string   `param:"in(query),nonzero"`
	Paragraph []string `param:"in(query),name(p),len(1:10)" regexp:"(^[\\w]*$)"`
	// Picture         multipart.FileHeader `param:"in(formData),name(pic),maxmb(30)"`
	Cookie       http.Cookie `param:"in(cookie),name(thinkgo)"`
	CookieString string      `param:"in(cookie),name(thinkgo)"`
}

var once sync.Once

// Implement the handler interface
func (p *Param) Serve(ctx *thinkgo.Context) error {
	ctx.Log().Info(ctx.R.Host)
	// name, id := ctx.GetSession("name"), ctx.GetSession("id")
	once.Do(func() {
		println("SetSession...")
		ctx.SetSession("name", "henry")
		ctx.SetSession("id", 123)
		ctx.SetCookie("thinkgo", "henrylee")
	})

	return ctx.JSON(200, p, true)
	// return ctx.String(200, "name: %v\nid: %d", name, id)
}

// Implementation notes of a response.
func (p *Param) Notes() thinkgo.Notes {
	return thinkgo.Notes{
		Note: "param desc",
		Return: thinkgo.JSONMsg{
			Code: 1,
			Info: "success",
		},
	}
}

package handler

import (
	// "mime/multipart"
	"net/http"
	"sync"

	"github.com/henrylee2cn/thinkgo"
)

type Index struct {
	Id        int      `param:"in(path),required,desc(ID),range(0:10)"`
	Num       float32  `param:"in(query),name(n),range(0.1:10)"`
	Title     string   `param:"in(query),nonzero"`
	Paragraph []string `param:"in(query),name(p),len(1:10)" regexp:"(^[\\w]*$)"`
	// Picture         multipart.FileHeader `param:"in(formData),name(pic),maxmb(30)"`
	Cookie          http.Cookie `param:"in(cookie),name(apiwareid)"`
	CookieString    string      `param:"in(cookie),name(apiwareid)"`
	thinkgo.Returns `param:"-" json:"-"`
}

var once sync.Once

func (i *Index) Serve(ctx *thinkgo.Context) error {
	// name, id := ctx.GetSession("name"), ctx.GetSession("id")
	// once.Do(func() {
	// 	println("SetSession...")
	// 	ctx.SetSession("name", "henry")
	// 	ctx.SetSession("id", 123)
	// })

	return ctx.JSON(200, i)
	// return ctx.String(200, "name: %v\nid: %d", name, id)
}

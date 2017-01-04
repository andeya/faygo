package handler

import (
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/henrylee2cn/thinkgo"
)

type Param struct {
	Id           int                   `param:"<in:path> <required> <desc:ID> <range: 0:10>"`
	Num          float32               `param:"<in:query> <required> <name:n> <range: 0.1:10> <err: query param 'n' must be number in 0.1~10>"`
	Title        string                `param:"<in:query> <nonzero>"`
	Paragraph    []string              `param:"<in:query> <name:p> <len: 1:10> <regexp: ^[\\w]*$>"`
	Picture      *multipart.FileHeader `param:"<in:formData> <name:pic> <maxmb:30>"`
	Cookie       http.Cookie           `param:"<in:cookie> <name:thinkgo>"`
	CookieString string                `param:"<in:cookie> <name:thinkgo>"`
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

	info, err := ctx.SaveFile("pic", false)
	ctx.Log().Infof("ctx.SaveFile: filename %s  url %s, size %d, err %v", p.Picture.Filename, info.Url, info.Size, err)
	return ctx.JSON(200,
		thinkgo.Map{
			"Struct Params":    p,
			"Additional Param": ctx.PathParam("additional"),
		}, true)
	// return ctx.String(200, "name: %v\nid: %d", name, id)
}

// Doc returns the API's note, result or parameters information.
func (p *Param) Doc() thinkgo.Doc {
	return thinkgo.Doc{
		Note: "param desc",
		Return: thinkgo.JSONMsg{
			Code: 1,
			Info: "success",
		},
		Params: []thinkgo.ParamInfo{
			{
				Name:  "additional",
				In:    "path",
				Model: "a",
				Desc:  "defined by the `Doc()` method",
			},
		},
	}
}

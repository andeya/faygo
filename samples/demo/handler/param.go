package handler

import (
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/henrylee2cn/faygo"
)

type Param struct {
	Id           int                   `param:"<in:path> <required> <desc:ID> <range: 0:10>"`
	Title        string                `param:"<in:query> <nonzero>"`
	Num          float32               `param:"<in:formData> <required> <name:n> <range: 0.1:10> <err: query param 'n' must be number in 0.1~10>"`
	Paragraph    []string              `param:"<in:formData> <name:p> <len: 1:20> <regexp: ^[\\w&=]*$>"`
	Picture      *multipart.FileHeader `param:"<in:formData> <name:pic> <maxmb:30>"`
	Cookie       *http.Cookie          `param:"<in:cookie> <name:faygo>"`
	CookieString string                `param:"<in:cookie> <name:faygo>"`
}

var once sync.Once

// Implement the handler interface
func (p *Param) Serve(ctx *faygo.Context) error {
	ctx.Log().Info(ctx.R.Host)
	once.Do(func() {
		println("Set session...")
		ctx.SetSession("name", "henry")
		ctx.SetCookie("faygo", "henrylee")
	})
	ctx.Log().Infof("Get session name=%v", ctx.GetSession("name"))

	info, err := ctx.SaveFile("pic", false)
	if err == nil {
		ctx.Log().Infof("ctx.SaveFile: filename %s  url %s, size %d", p.Picture.Filename, info.Url, info.Size)
	}
	return ctx.JSON(200,
		faygo.Map{
			"Struct Params":    p,
			"Additional Param": ctx.PathParam("additional"),
		}, true)
	// return ctx.String(200, "name=%v", name)
}

// Doc returns the API's note, result or parameters information.
func (p *Param) Doc() faygo.Doc {
	return faygo.Doc{
		Note: "param desc",
		Return: faygo.JSONMsg{
			Code: 1,
			Info: "success",
		},
		Params: []faygo.ParamInfo{
			{
				Name:  "additional",
				In:    "path",
				Model: "a",
				Desc:  "defined by the `Doc()` method",
			},
		},
	}
}

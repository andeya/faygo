package main

import (
	"github.com/henrylee2cn/thinkgo"
	"time"
)

type Index struct {
	Id        int      `param:"in(path),required,desc(ID),range(0:10)"`
	Title     string   `param:"in(query),nonzero"`
	Paragraph []string `param:"in(query),name(p),len(1:10)" regexp:"(^[\\w]*$)"`
	Cookie    string   `param:"in(cookie),name(thinkgoID)"`
	// Picture         multipart.FileHeader `param:"in(formData),name(pic),maxmb(30)"`
}

func (i *Index) Serve(ctx *thinkgo.Context) error {
	if ctx.CookieParam("thinkgoID") == "" {
		ctx.SetCookie("thinkgoID", time.Now().String())
	}
	return ctx.JSON(200, i)
}

func main() {
	thinkgo.Root().GET("/index/:id", new(Index))
	// or thinkgo.Route(thinkgo.GET("/index/:id", new(Index)))
	thinkgo.Run()
}

/*
http GET:
    http://localhost:8080/index/1?title=test&p=abc&p=xyz
response:
    {
      "Id": 1,
      "Title": "test",
      "Paragraph": [
        "abc",
        "xyz"
      ],
      "Cookie": "2016-11-13 01:14:40.9038005 +0800 CST"
    }
*/

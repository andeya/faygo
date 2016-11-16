package main

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/router"
)

func main() {
	thinkgo.Global.SetUpload("./upload/0")
	thinkgo.Global.SetStatic("./static/0")
	thinkgo.Init("defaultapp")
	go thinkgo.Run()

	app := thinkgo.New("testapp", "0.1")
	router.Route(app)
	app.Run()
}

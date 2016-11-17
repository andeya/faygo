package main

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/router"
)

func main() {
	thinkgo.Global.SetUpload("./upload/0")
	thinkgo.Global.SetStatic("./static/0")
	app1 := thinkgo.New("testapp1", "0.1")
	router.Route1(app1)
	go app1.Run()

	app2 := thinkgo.New("testapp2", "0.1")
	router.Route2(app2)
	app2.Run()
}

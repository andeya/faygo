package main

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/router"
)

func main() {
	thinkgo.Global.SetUpload("./upload/0")
	// thinkgo.Global.SetStatic("./static")
	app1 := thinkgo.New("myapp1", "1.0")
	router.Route1(app1)
	go app1.Run()

	app2 := thinkgo.New("myapp2", "1.0")
	router.Route2(app2)
	app2.Run()
}

package main

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/samples/demo/router"
)

func main() {
	thinkgo.SetUpload("./upload/0", false, false)
	// thinkgo.SetStatic("./static", false, false)
	{
		app1 := thinkgo.New("myapp1", "1.0")
		router.Route1(app1)
	}
	{
		app2 := thinkgo.New("myapp2", "1.0")
		router.Route2(app2)
	}
	thinkgo.Run()
}

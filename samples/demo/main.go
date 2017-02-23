package main

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/samples/demo/router"
	"time"
)

func main() {
	thinkgo.SetShutdown(time.Minute, func() error {
		thinkgo.Debug("finalizer 等待5s...")
		time.Sleep(5 * time.Second)
		thinkgo.Debug("finalizer 5s到时！")
		return nil
	})
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

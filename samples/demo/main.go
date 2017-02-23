package main

import (
	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/samples/demo/router"
	"time"
)

func main() {
	faygo.SetShutdown(time.Minute, func() error {
		faygo.Debug("finalizer 等待3s...")
		time.Sleep(3 * time.Second)
		faygo.Debug("finalizer 3s到时！")
		return nil
	})
	faygo.SetUpload("./upload/0", false, false)
	// faygo.SetStatic("./static", false, false)
	{
		app1 := faygo.New("myapp1", "1.0")
		router.Route1(app1)
	}
	{
		app2 := faygo.New("myapp2", "1.0")
		router.Route2(app2)
	}
	faygo.Run()
}

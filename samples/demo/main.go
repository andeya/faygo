package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/samples/demo/router"
)

// Browse 'http://localhost:8080/apidoc' and 'http://localhost:8081/apidoc' to test.

// run type 1
func main() {
	// pprof
	// http://localhost:7777/debug/pprof
	go pprofServer()

	faygo.SetShutdown(time.Minute, func() error {
		faygo.Debug("Before services are closing1: wait for 1s...")
		time.Sleep(1 * time.Second)
		faygo.Debug("Before services are closed1!")
		return nil
	}, func() error {
		faygo.Debug("After services are closing2: wait for 1s...")
		time.Sleep(1 * time.Second)
		faygo.Debug("After services are closed2!")
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

// run type 2
func main2() {
	// pprof
	// http://localhost:7777/debug/pprof
	go pprofServer()

	faygo.SetShutdown(time.Minute, func() error {
		faygo.Debug("Before services are closing1: wait for 1s...")
		time.Sleep(1 * time.Second)
		faygo.Debug("Before services are closed1!")
		return nil
	}, func() error {
		faygo.Debug("After services are closing2: wait for 1s...")
		time.Sleep(1 * time.Second)
		faygo.Debug("After services are closed2!")
		return nil
	})

	faygo.SetUpload("./upload/0", false, false)
	// faygo.SetStatic("./static", false, false)
	{
		app1 := faygo.New("myapp1", "1.0")
		router.Route1(app1)
		go app1.Run()
	}
	{
		app2 := faygo.New("myapp2", "1.0")
		router.Route2(app2)
		go app2.Run()

	}
	select {}
}

// only for pprof
// http://localhost:7777/debug/pprof
func pprofServer() {
	http.ListenAndServe("0.0.0.0:7777", nil)
}

package main

import (
	"github.com/henrylee2cn/faygo"
)

func main() {
	cfg := faygo.NewDefaultConfig()
	app := faygo.NewWithConfig(cfg, "myconfig", "v1")
	app.Run()
}

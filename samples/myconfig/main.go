package main

import (
	"github.com/andeya/faygo"
)

func main() {
	cfg := faygo.NewDefaultConfig()
	app := faygo.NewWithConfig(cfg, "myconfig", "v1")
	app.Run()
}

package main

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/samples/directsql/router"
)

func main() {
	{
		directsql := thinkgo.New("directsqldemo", "1.0")
		router.Route(directsql)
		// go directsql.Run()
	}
	thinkgo.Run()
}

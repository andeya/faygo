package main

import (
	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/henrylee2cn/thinkgo"
	_ "github.com/henrylee2cn/thinkgo/samples/directsql/common"
	"github.com/henrylee2cn/thinkgo/samples/directsql/router"
)

func main() {
	{
		directsql := thinkgo.New("directsqldemo", "1.0")
		router.Route(directsql)
	}
	thinkgo.Run()
}

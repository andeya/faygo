package main

import (
	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/henrylee2cn/faygo"
	_ "github.com/henrylee2cn/faygo/samples/directsql/common"
	"github.com/henrylee2cn/faygo/samples/directsql/router"
)

func main() {
	{
		directsql := faygo.New("directsqldemo", "1.0")
		router.Route(directsql)
	}
	faygo.Run()
}

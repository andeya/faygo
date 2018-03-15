package main

import (
	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/henrylee2cn/faygo"
	_ "github.com/henrylee2cn/faygo/samples/directsqlx/common"
	"github.com/henrylee2cn/faygo/samples/directsqlx/router"
)

func main() {
	{
		directsqlx := faygo.New("directsqlx demo", "1.0")
		router.Route(directsqlx)
	}
	faygo.Run()
}

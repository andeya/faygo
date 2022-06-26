package main

import (
	"github.com/andeya/faygo"
	_ "github.com/andeya/faygo/samples/directsql/common"
	"github.com/andeya/faygo/samples/directsql/router"
	_ "github.com/go-sql-driver/mysql" // mysql driver
)

func main() {
	{
		directsql := faygo.New("directsqldemo", "1.0")
		router.Route(directsql)
	}
	faygo.Run()
}

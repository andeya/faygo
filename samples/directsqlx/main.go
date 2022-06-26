package main

import (
	"github.com/andeya/faygo"
	_ "github.com/andeya/faygo/samples/directsqlx/common"
	"github.com/andeya/faygo/samples/directsqlx/router"
	_ "github.com/go-sql-driver/mysql" // mysql driver
)

func main() {
	{
		directsqlx := faygo.New("directsqlx demo", "1.0")
		router.Route(directsqlx)
	}
	faygo.Run()
}

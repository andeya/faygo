package model

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/henrylee2cn/thinkgo/conf"
)

var GormMysql = new(gorm.DB)

func init() {
	var err error
	*GormMysql, err = gorm.Open("mysql", conf.MYSQL_URL)
	if err != nil {
		log.Println(err)
	}

	err = GormMysql.DB().Ping()
	if err != nil {
		log.Println(err)
	}

	// 最大空闲连接数
	GormMysql.DB().SetMaxIdleConns(10)
	// 最大打开连接数
	GormMysql.DB().SetMaxOpenConns(100)
	// Disable table name's pluralization
	GormMysql.SingularTable(true)
}

package gorm

import (
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
	// _ "github.com/jinzhu/gorm/dialects/mysql" //github.com/go-sql-driver/mysql
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite" //github.com/mattn/go-sqlite3

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/utils"
)

type DBService struct {
	Default *gorm.DB            // the default database engine
	List    map[string]*gorm.DB // database engine list
}

var dbService = func() (serv *DBService) {
	serv = &DBService{
		List: map[string]*gorm.DB{},
	}

	defer func() {
		if serv.Default == nil {
			time.Sleep(2e9)
		}
	}()

	err := loadDBConfig()
	if err != nil {
		thinkgo.Error(err.Error())
	}

	for _, conf := range dbConfigs {
		engine, err := gorm.Open(conf.Driver, conf.Connstring)
		if err != nil {
			thinkgo.Error(err.Error())
			continue
		}
		engine.SetLogger(thinkgo.NewLog())
		engine.LogMode(conf.ShowSql)

		engine.DB().SetMaxOpenConns(conf.MaxOpenConns)
		engine.DB().SetMaxIdleConns(conf.MaxIdleConns)

		if conf.Driver == "sqlite3" && !utils.FileExists(conf.Connstring) {
			os.MkdirAll(filepath.Dir(conf.Connstring), 0777)
			f, err := os.Create(conf.Connstring)
			if err != nil {
				thinkgo.Error(err.Error())
			} else {
				f.Close()
			}
		}

		serv.List[conf.Name] = engine
		if DEFAULTDB_NAME == conf.Name {
			serv.Default = engine
		}
	}
	return
}()

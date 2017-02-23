package gorm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"

	// _ "github.com/jinzhu/gorm/dialects/mssql"    //github.com/denisenkom/go-mssqldb
	_ "github.com/jinzhu/gorm/dialects/mysql"    //github.com/go-sql-driver/mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" //github.com/lib/pq
	// _ "github.com/jinzhu/gorm/dialects/sqlite"   //github.com/mattn/go-sqlite3

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/utils"
)

// DBService is a database engine object.
type DBService struct {
	Default *gorm.DB            // the default database engine
	List    map[string]*gorm.DB // database engine list
}

var dbService = func() (serv *DBService) {
	serv = &DBService{
		List: map[string]*gorm.DB{},
	}

	var errs []string
	defer func() {
		if len(errs) > 0 {
			panic("[gorm] " + strings.Join(errs, "\n"))
		}
		if serv.Default == nil {
			faygo.Panicf("[gorm] the `default` database engine must be configured and enabled")
		}
	}()

	err := loadDBConfig()
	if err != nil {
		faygo.Panicf("[gorm]", err.Error())
		return
	}

	for _, conf := range dbConfigs {
		if !conf.Enable {
			continue
		}
		engine, err := gorm.Open(conf.Driver, conf.Connstring)
		if err != nil {
			faygo.Critical("[gorm]", err.Error())
			errs = append(errs, err.Error())
			continue
		}
		engine.SetLogger(faygo.NewLog())
		engine.LogMode(conf.ShowSql)

		engine.DB().SetMaxOpenConns(conf.MaxOpenConns)
		engine.DB().SetMaxIdleConns(conf.MaxIdleConns)

		if conf.Driver == "sqlite3" && !utils.FileExists(conf.Connstring) {
			os.MkdirAll(filepath.Dir(conf.Connstring), 0777)
			f, err := os.Create(conf.Connstring)
			if err != nil {
				faygo.Critical("[gorm]", err.Error())
				errs = append(errs, err.Error())
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

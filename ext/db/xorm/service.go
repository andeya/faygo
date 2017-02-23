package xorm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	// _ "github.com/denisenkom/go-mssqldb" //mssql
	_ "github.com/go-sql-driver/mysql" //mysql
	_ "github.com/lib/pq"              //postgres
	// _ "github.com/mattn/go-oci8"         //oracle(need to install the pkg-config utility)
	// _ "github.com/mattn/go-sqlite3"      //sqlite

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/utils"
)

// DBService is a database engine object.
type DBService struct {
	Default *xorm.Engine            // the default database engine
	List    map[string]*xorm.Engine // database engine list
}

var dbService = func() (serv *DBService) {
	serv = &DBService{
		List: map[string]*xorm.Engine{},
	}
	var errs []string
	defer func() {
		if len(errs) > 0 {
			panic("[xorm] " + strings.Join(errs, "\n"))
		}
		if serv.Default == nil {
			faygo.Panicf("[xorm] the `default` database engine must be configured and enabled")
		}
	}()

	err := loadDBConfig()
	if err != nil {
		faygo.Panicf("[xorm]", err.Error())
		return
	}

	for _, conf := range dbConfigs {
		if !conf.Enable {
			continue
		}
		engine, err := xorm.NewEngine(conf.Driver, conf.Connstring)
		if err != nil {
			faygo.Critical("[xorm]", err.Error())
			errs = append(errs, err.Error())
			continue
		}
		err = engine.Ping()
		if err != nil {
			faygo.Critical("[xorm]", err.Error())
			errs = append(errs, err.Error())
			continue
		}
		engine.SetLogger(iLogger)
		engine.SetMaxOpenConns(conf.MaxOpenConns)
		engine.SetMaxIdleConns(conf.MaxIdleConns)
		engine.SetDisableGlobalCache(conf.DisableCache)
		engine.ShowSQL(conf.ShowSql)
		engine.ShowExecTime(conf.ShowExecTime)

		if (conf.TableFix == "prefix" || conf.TableFix == "suffix") && len(conf.TableSpace) > 0 {
			var impr core.IMapper
			if conf.TableSnake {
				impr = core.SnakeMapper{}
			} else {
				impr = core.SameMapper{}
			}
			if conf.TableFix == "prefix" {
				engine.SetTableMapper(core.NewPrefixMapper(impr, conf.TableSpace))
			} else {
				engine.SetTableMapper(core.NewSuffixMapper(impr, conf.TableSpace))
			}
		}

		if (conf.ColumnFix == "prefix" || conf.ColumnFix == "suffix") && len(conf.ColumnSpace) > 0 {
			var impr core.IMapper
			if conf.ColumnSnake {
				impr = core.SnakeMapper{}
			} else {
				impr = core.SameMapper{}
			}
			if conf.ColumnFix == "prefix" {
				engine.SetTableMapper(core.NewPrefixMapper(impr, conf.ColumnSpace))
			} else {
				engine.SetTableMapper(core.NewSuffixMapper(impr, conf.ColumnSpace))
			}
		}

		if conf.Driver == "sqlite3" && !utils.FileExists(conf.Connstring) {
			os.MkdirAll(filepath.Dir(conf.Connstring), 0777)
			f, err := os.Create(conf.Connstring)
			if err != nil {
				faygo.Critical("[xorm]", err.Error())
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

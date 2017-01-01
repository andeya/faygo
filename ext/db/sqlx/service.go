package sqlx

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"

	// _ "github.com/denisenkom/go-mssqldb" //mssql
	_ "github.com/go-sql-driver/mysql" //mysql
	_ "github.com/lib/pq"              //postgres
	// _ "github.com/mattn/go-oci8"         //oracle(need to install the pkg-config utility)
	// _ "github.com/mattn/go-sqlite3"      //sqlite

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/utils"
)

// DBService is a database engine object.
type DBService struct {
	Default *sqlx.DB            // the default database engine
	List    map[string]*sqlx.DB // database engine list
}

var dbService = func() (serv *DBService) {
	serv = &DBService{
		List: map[string]*sqlx.DB{},
	}

	defer func() {
		if serv.Default == nil {
			time.Sleep(2e9)
		}
	}()

	err := loadDBConfig()
	if err != nil {
		thinkgo.Error(err.Error())
		return
	}

	for _, conf := range dbConfigs {
		db, err := sqlx.Connect(conf.Driver, conf.Connstring)
		if err != nil {
			thinkgo.Error(err.Error())
			continue
		}

		db.SetMaxOpenConns(conf.MaxOpenConns)
		db.SetMaxIdleConns(conf.MaxIdleConns)

		var strFunc = strings.ToLower
		if conf.ColumnSnake {
			strFunc = utils.SnakeString
		}

		// Create a new mapper which will use the struct field tag "json" instead of "db"
		db.Mapper = reflectx.NewMapperFunc(conf.StructTag, strFunc)

		if conf.Driver == "sqlite3" && !utils.FileExists(conf.Connstring) {
			os.MkdirAll(filepath.Dir(conf.Connstring), 0777)
			f, err := os.Create(conf.Connstring)
			if err != nil {
				thinkgo.Error(err.Error())
			} else {
				f.Close()
			}
		}

		serv.List[conf.Name] = db
		if DEFAULTDB_NAME == conf.Name {
			serv.Default = db
		}
	}
	return
}()

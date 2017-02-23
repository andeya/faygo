package sqlx

import (
	"database/sql"
	"errors"

	"github.com/henrylee2cn/faygo"
	"github.com/jmoiron/sqlx"
)

// MustDB g the specified database engine,
// or the default DB if no name is specified.
func MustDB(name ...string) *sqlx.DB {
	if len(name) == 0 {
		return dbService.Default
	}
	db, ok := dbService.List[name[0]]
	if !ok {
		faygo.Panicf("[sqlx] the database engine `%s` is not configured", name[0])
	}
	return db
}

// DB is similar to MustDB, but safe.
func DB(name ...string) (*sqlx.DB, bool) {
	if len(name) == 0 {
		return dbService.Default, true
	}
	engine, ok := dbService.List[name[0]]
	return engine, ok
}

// List gets the list of database engines
func List() map[string]*sqlx.DB {
	return dbService.List
}

// MustConfig gets the configuration information for the specified database,
// or returns the default if no name is specified.
func MustConfig(name ...string) DBConfig {
	if len(name) == 0 {
		return *defaultConfig
	}
	config, ok := dbConfigs[name[0]]
	if !ok {
		faygo.Panicf("[sqlx] the database engine `%s` is not configured", name[0])
	}
	return *config
}

// Config is similar to MustConfig, but safe.
func Config(name ...string) (DBConfig, bool) {
	if len(name) == 0 {
		return *defaultConfig, true
	}
	config, ok := dbConfigs[name[0]]
	if !ok {
		return DBConfig{}, false
	}
	return *config, true
}

// Callback uses the `default` database for non-transactional operations.
func Callback(fn func(DBTX) error, tx ...*sqlx.Tx) error {
	if fn == nil {
		return nil
	}
	if len(tx) > 0 && tx[0] != nil {
		return fn(tx[0])
	}
	return fn(MustDB())
}

// CallbackByName uses the specified database for non-transactional operations.
func CallbackByName(dbName string, fn func(DBTX) error, tx ...*sqlx.Tx) error {
	if fn == nil {
		return nil
	}
	if len(tx) > 0 && tx[0] != nil {
		return fn(tx[0])
	}
	engine, ok := DB(dbName)
	if !ok {
		return errors.New("[sqlx] the database engine `" + dbName + "` is not configured")
	}
	return fn(engine)
}

// TransactCallback uses the default database for transactional operations.
// note: if an error is returned, the rollback method should be invoked outside the function.
func TransactCallback(fn func(*sqlx.Tx) error, tx ...*sqlx.Tx) (err error) {
	if fn == nil {
		return
	}
	var _tx *sqlx.Tx
	if len(tx) > 0 {
		_tx = tx[0]
	}
	if _tx == nil {
		_tx, err = MustDB().Beginx()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				_tx.Rollback()
			} else {
				_tx.Commit()
			}
		}()
	}
	err = fn(_tx)
	return
}

// TransactCallbackByName uses the `specified` database for transactional operations.
// note: if an error is returned, the rollback method should be invoked outside the function.
func TransactCallbackByName(dbName string, fn func(*sqlx.Tx) error, tx ...*sqlx.Tx) (err error) {
	if fn == nil {
		return
	}
	var _tx *sqlx.Tx
	if len(tx) > 0 {
		_tx = tx[0]
	}
	if _tx == nil {
		engine, ok := DB(dbName)
		if !ok {
			return errors.New("[sqlx] the database engine `" + dbName + "` is not configured")
		}
		_tx, err = engine.Beginx()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				_tx.Rollback()
			} else {
				_tx.Commit()
			}
		}()
	}
	err = fn(_tx)
	return
}

// DBTX contains all the exportable methods of * sqlx.TX
type DBTX interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error

	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

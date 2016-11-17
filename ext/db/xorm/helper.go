package xorm

import (
	"errors"

	"github.com/go-xorm/xorm"
)

// Gets the specified database engine,
// or the default DB if no name is specified.
func MustDB(name ...string) *xorm.Engine {
	db, _ := DB(name...)
	return db
}

// DB is similar to MustDB, but safe.
func DB(name ...string) (*xorm.Engine, bool) {
	if len(name) == 0 {
		return dbService.Default, true
	}
	engine, ok := dbService.List[name[0]]
	return engine, ok
}

// List gets the list of database engines
func List() map[string]*xorm.Engine {
	return dbService.List
}

// Gets the connection string for the specified database,
// or returns the default if no name is specified.
func MustConnstring(name ...string) string {
	conn, _ := Connstring(name...)
	return conn
}

// Connstring is similar to MustConnstring, but safe.
func Connstring(name ...string) (string, bool) {
	if len(name) == 0 {
		return defaultConfig.Connstring, true
	}
	config, ok := dbConfigs[name[0]]
	if !ok {
		return "", false
	}
	return config.Connstring, true
}

type Table interface {
	TableName() string
}

// A callback function that uses the `default` database for non-transactional operations.
func Callback(fn func(*xorm.Session) error, session ...*xorm.Session) error {
	if fn == nil {
		return nil
	}
	var sess *xorm.Session
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		sess = MustDB().NewSession()
		defer sess.Close()
	}
	return fn(sess)
}

// A callback function that uses the specified database for non-transactional operations.
func CallbackByName(dbName string, fn func(*xorm.Session) error, session ...*xorm.Session) error {
	if fn == nil {
		return nil
	}
	var sess *xorm.Session
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		engine, ok := DB(dbName)
		if !ok {
			return errors.New("the specified database engine " + dbName + " is not configured")
		}
		sess = engine.NewSession()
		defer sess.Close()
	}
	return fn(sess)
}

// A callback function that uses the default database for transactional operations.
// note: if an error is returned, the rollback method should be invoked outside the function.
func TransactCallback(fn func(*xorm.Session) error, session ...*xorm.Session) (err error) {
	if fn == nil {
		return
	}
	var sess *xorm.Session
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		sess = MustDB().NewSession()
		defer sess.Close()
		err = sess.Begin()
		if err != nil {
			return
		}
		defer func() {
			if err != nil {
				sess.Rollback()
				return
			}
			err = sess.Commit()
		}()
	}
	err = fn(sess)
	return
}

// A callback function that uses the `specified` database for transactional operations.
// note: if an error is returned, the rollback method should be invoked outside the function.
func TransactCallbackByName(dbName string, fn func(*xorm.Session) error, session ...*xorm.Session) (err error) {
	if fn == nil {
		return
	}
	var sess *xorm.Session
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		engine, ok := DB(dbName)
		if !ok {
			return errors.New("the specified database engine " + dbName + " is not configured")
		}
		sess = engine.NewSession()
		defer sess.Close()
		err = sess.Begin()
		if err != nil {
			return
		}
		defer func() {
			if err != nil {
				sess.Rollback()
				return
			}
			err = sess.Commit()
		}()
	}
	err = fn(sess)
	return
}

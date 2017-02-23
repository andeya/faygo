package gorm

import (
	"errors"

	"github.com/henrylee2cn/faygo"
	"github.com/jinzhu/gorm"
)

// MustDB gets the specified database engine,
// or the default DB if no name is specified.
func MustDB(name ...string) *gorm.DB {
	if len(name) == 0 {
		return dbService.Default
	}
	db, ok := dbService.List[name[0]]
	if !ok {
		faygo.Panicf("[gorm] the database engine `%s` is not configured", name[0])
	}
	return db
}

// DB is similar to MustDB, but safe.
func DB(name ...string) (*gorm.DB, bool) {
	if len(name) == 0 {
		return dbService.Default, true
	}
	engine, ok := dbService.List[name[0]]
	return engine, ok
}

// List gets the list of database engines
func List() map[string]*gorm.DB {
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
		faygo.Panicf("[gorm] the database engine `%s` is not configured", name[0])
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
func Callback(fn func(*gorm.DB) error, session ...*gorm.DB) error {
	if fn == nil {
		return nil
	}
	var sess *gorm.DB
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		sess = MustDB().New()
		defer sess.Close()
	}
	return fn(sess)
}

// CallbackByName uses the specified database for non-transactional operations.
func CallbackByName(dbName string, fn func(*gorm.DB) error, session ...*gorm.DB) error {
	if fn == nil {
		return nil
	}
	var sess *gorm.DB
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		engine, ok := DB(dbName)
		if !ok {
			return errors.New("[gorm] the database engine `" + dbName + "` is not configured")
		}
		sess = engine.New()
		defer sess.Close()
	}
	return fn(sess)
}

// TransactCallback uses the default database for transactional operations.
// note: if an error is returned, the rollback method should be invoked outside the function.
func TransactCallback(fn func(*gorm.DB) error, session ...*gorm.DB) (err error) {
	if fn == nil {
		return
	}
	var sess *gorm.DB
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		sess = MustDB().New().Begin()
		defer func() {
			if err != nil {
				sess.Rollback()
			} else {
				sess.Commit()
			}
			sess.Close()
		}()
	}
	err = fn(sess)
	return
}

// TransactCallbackByName uses the `specified` database for transactional operations.
// note: if an error is returned, the rollback method should be invoked outside the function.
func TransactCallbackByName(dbName string, fn func(*gorm.DB) error, session ...*gorm.DB) (err error) {
	if fn == nil {
		return
	}
	var sess *gorm.DB
	if len(session) > 0 {
		sess = session[0]
	}
	if sess == nil {
		engine, ok := DB(dbName)
		if !ok {
			return errors.New("[gorm] the database engine `" + dbName + "` is not configured")
		}
		sess = engine.New().Begin()
		defer func() {
			if err != nil {
				sess.Rollback()
			} else {
				sess.Commit()
			}
			sess.Close()
		}()
	}
	err = fn(sess)
	return
}

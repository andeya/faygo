package gorm

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/jinzhu/gorm"
)

// Gets the specified database engine,
// or the default DB if no name is specified.
func MustDB(name ...string) *gorm.DB {
	db, ok := DB(name...)
	if !ok {
		_name := "default"
		if len(name) == 0 {
			_name = name[0]
		}
		thinkgo.Panicf("the database engine `%s` is not configured", _name)
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

// Gets the connection string for the specified database,
// or returns the default if no name is specified.
func MustConnstring(name ...string) string {
	conn, ok := Connstring(name...)
	if !ok {
		_name := "default"
		if len(name) == 0 {
			_name = name[0]
		}
		thinkgo.Panicf("the database engine `%s` is not configured", _name)
	}
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

package gorm

import (
	"github.com/jinzhu/gorm"
)

// Gets the specified database engine,
// or the default DB if no name is specified.
func MustDB(name ...string) *gorm.DB {
	db, _ := DB(name...)
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

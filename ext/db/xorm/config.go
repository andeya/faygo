package xorm

import (
	"os"
	"path/filepath"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/ini"
)

// DBConfig is database connection config
type DBConfig struct {
	Name         string `ini:"-"`
	Driver       string `ini:"driver"` // driver: mssql | odbc(mssql) | mysql | mymysql | postgres | sqlite3 | oci8 | goracle
	Connstring   string `ini:"connstring"`
	MaxOpenConns int    `ini:"max_open_conns"`
	MaxIdleConns int    `ini:"max_idle_conns"`
	ShowSql      bool   `ini:"show_sql"`
	TableFix     string `ini:"table_fix"`    // the table namespace is located in the prefix or suffix: prefix | suffix
	TableSpace   string `ini:"table_space"`  // table namespace
	TableSnake   bool   `ini:"table_snake"`  // the table name uses the snake style or remains unchanged
	ColumnFix    string `ini:"column_fix"`   // column namespace is located in the prefix or suffix: prefix | suffix
	ColumnSpace  string `ini:"column_space"` // column namespace
	ColumnSnake  bool   `ini:"column_snake"` // the column name uses the snake style or remains unchanged
	DisableCache bool   `ini:"disable_cache"`
	ShowExecTime bool   `ini:"show_exec_time"`
}

const (
	DBCONFIG_FILE  = thinkgo.CONFIG_DIR + "xorm.ini"
	DATABASE_DIR   = "database/"
	DEFAULTDB_NAME = "default"
)

var (
	dbConfigs     = map[string]*DBConfig{DEFAULTDB_NAME: defaultConfig}
	defaultConfig = &DBConfig{
		Name:         DEFAULTDB_NAME,
		Driver:       "sqlite3",
		Connstring:   DATABASE_DIR + "sqlite.db",
		MaxOpenConns: 0,
		MaxIdleConns: 0,
		TableFix:     "prefix",
		TableSpace:   "",
		TableSnake:   true,
		ColumnFix:    "prefix",
		ColumnSpace:  "",
		ColumnSnake:  true,
		DisableCache: false,
		ShowExecTime: false,
		ShowSql:      false,
	}
)

func loadDBConfig() error {
	os.MkdirAll(filepath.Dir(DBCONFIG_FILE), 0777)
	cfg, err := ini.LooseLoad(DBCONFIG_FILE)
	if err != nil {
		return err
	}
	for _, section := range cfg.Sections() {
		if section.Name() == ini.DEFAULT_SECTION {
			continue
		}
		var dbConfig *DBConfig
		if section.Name() == DEFAULTDB_NAME {
			dbConfig = defaultConfig
		} else {
			dbConfig = &DBConfig{Name: section.Name()}
		}
		err := section.MapTo(dbConfig)
		if err != nil {
			return err
		}
		dbConfigs[dbConfig.Name] = dbConfig
	}
	_, err = cfg.GetSection(DEFAULTDB_NAME)
	if err != nil {
		sec, _ := cfg.NewSection(DEFAULTDB_NAME)
		err := sec.ReflectFrom(defaultConfig)
		if err != nil {
			return err
		}
	}
	return cfg.SaveTo(DBCONFIG_FILE)
}

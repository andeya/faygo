package xorm

import (
	"github.com/go-xorm/core"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/logging"
)

// ILogger logger
type ILogger struct {
	logging *logging.Logger
	level   core.LogLevel
	showSQL bool
}

var iLogger = func() *ILogger {
	log := &ILogger{
		logging: faygo.NewLog(),
	}
	log.logging.ExtraCalldepth++
	return log
}()

// Debug DEBUG level log
func (i *ILogger) Debug(v ...interface{}) {
	i.logging.Debug(v...)
}

// Debug DEBUG level log with format
func (i *ILogger) Debugf(format string, v ...interface{}) {
	i.logging.Debugf(format, v...)
}

// Debug ERROR level log
func (i *ILogger) Error(v ...interface{}) {
	i.logging.Error(v...)
}

// Debug ERROR level log with format
func (i *ILogger) Errorf(format string, v ...interface{}) {
	i.logging.Errorf(format, v...)
}

// Debug INFO level log
func (i *ILogger) Info(v ...interface{}) {
	i.logging.Info(v...)
}

// Debug INFO level log with format
func (i *ILogger) Infof(format string, v ...interface{}) {
	i.logging.Infof(format, v...)
}

// Debug WARN level log
func (i *ILogger) Warn(v ...interface{}) {
	i.logging.Warn(v...)
}

// Debug WARN level log with format
func (i *ILogger) Warnf(format string, v ...interface{}) {
	i.logging.Warnf(format, v...)
}

// Level returns log level
func (i *ILogger) Level() core.LogLevel {
	return core.LOG_UNKNOWN
}

// SetLevel sets log level
func (i *ILogger) SetLevel(l core.LogLevel) {}

func (i *ILogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		i.showSQL = true
		return
	}
	i.showSQL = show[0]
}

// IsShowSQL returns if it wills show SQL
func (i *ILogger) IsShowSQL() bool {
	return i.showSQL
}

// func conversionLevel(level logging.Level) core.LogLevel {
// 	switch level {
// 	case logging.CRITICAL, logging.ERROR:
// 		return core.LOG_ERR
// 	case logging.WARNING:
// 		return core.LOG_WARNING
// 	case logging.NOTICE, logging.INFO:
// 		return core.LOG_INFO
// 	case logging.DEBUG:
// 		return core.LOG_DEBUG
// 	default:
// 		return core.LOG_UNKNOWN
// 	}
// }

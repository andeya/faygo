package xorm

import (
	"github.com/go-xorm/core"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/logging"
)

type ILogger struct {
	logging *logging.Logger
	level   core.LogLevel
	showSQL bool
}

var iLogger = func() *ILogger {
	log := &ILogger{
		logging: thinkgo.NewLog(),
	}
	log.logging.ExtraCalldepth++
	return log
}()

func (i *ILogger) Debug(v ...interface{}) {
	i.logging.Debug(v...)
}

func (i *ILogger) Debugf(format string, v ...interface{}) {
	i.logging.Debugf(format, v...)
}

func (i *ILogger) Error(v ...interface{}) {
	i.logging.Error(v...)
}

func (i *ILogger) Errorf(format string, v ...interface{}) {
	i.logging.Errorf(format, v...)
}

func (i *ILogger) Info(v ...interface{}) {
	i.logging.Info(v...)
}

func (i *ILogger) Infof(format string, v ...interface{}) {
	i.logging.Infof(format, v...)
}

func (i *ILogger) Warn(v ...interface{}) {
	i.logging.Warn(v...)
}
func (i *ILogger) Warnf(format string, v ...interface{}) {
	i.logging.Warnf(format, v...)
}

func (i *ILogger) Level() core.LogLevel {
	return core.LOG_UNKNOWN
}

func (i *ILogger) SetLevel(l core.LogLevel) {}

func (i *ILogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		i.showSQL = true
		return
	}
	i.showSQL = show[0]
}

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

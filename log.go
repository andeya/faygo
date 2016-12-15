// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package thinkgo

import (
	"log"
	"strings"

	"github.com/henrylee2cn/thinkgo/logging"
	"github.com/henrylee2cn/thinkgo/logging/color"
)

// NewLog gets a global logger
func NewLog() *logging.Logger {
	newlog := *Global.bizlog
	newlog.ExtraCalldepth--
	return &newlog
}

// Sys logs a system message using CRITICAL as log level.
func Sys(args ...interface{}) {
	Global.syslog.Critical(args...)
}

// Sys logs a system message using CRITICAL as log level.
func Sysf(format string, args ...interface{}) {
	Global.syslog.Criticalf(format, args...)
}

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	Global.bizlog.Fatal(args...)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	Global.bizlog.Fatalf(format, args...)
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func Panic(args ...interface{}) {
	Global.bizlog.Panic(args...)
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func Panicf(format string, args ...interface{}) {
	Global.bizlog.Panicf(format, args...)
}

// Critical logs a message using CRITICAL as log level.
func Critical(args ...interface{}) {
	Global.bizlog.Critical(args...)
}

// Criticalf logs a message using CRITICAL as log level.
func Criticalf(format string, args ...interface{}) {
	Global.bizlog.Criticalf(format, args...)
}

// Error logs a message using ERROR as log level.
func Error(args ...interface{}) {
	Global.bizlog.Error(args...)
}

// Errorf logs a message using ERROR as log level.
func Errorf(format string, args ...interface{}) {
	Global.bizlog.Errorf(format, args...)
}

// Warning logs a message using WARNING as log level.
func Warning(args ...interface{}) {
	Global.bizlog.Warning(args...)
}

// Warningf logs a message using WARNING as log level.
func Warningf(format string, args ...interface{}) {
	Global.bizlog.Warningf(format, args...)
}

// Notice logs a message using NOTICE as log level.
func Notice(args ...interface{}) {
	Global.bizlog.Notice(args...)
}

// Noticef logs a message using NOTICE as log level.
func Noticef(format string, args ...interface{}) {
	Global.bizlog.Noticef(format, args...)
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
	Global.bizlog.Info(args...)
}

// Infof logs a message using INFO as log level.
func Infof(format string, args ...interface{}) {
	Global.bizlog.Infof(format, args...)
}

// Debug logs a message using DEBUG as log level.
func Debug(args ...interface{}) {
	Global.bizlog.Debug(args...)
}

// Debugf logs a message using DEBUG as log level.
func Debugf(format string, args ...interface{}) {
	Global.bizlog.Debugf(format, args...)
}

var (
	consoleLogBackend = &logging.LogBackend{
		Logger: log.New(color.NewColorableStdout(), "", 0),
		Color:  true,
	}
	fileBackend *logging.FileBackend
)

func (global *GlobalSetting) initLogger() {
	fileBackend = func() *logging.FileBackend {
		fileBackend, err := logging.NewDefaultFileBackend(global.logDir + "thinkgo.log")
		if err != nil {
			panic(err)
		}
		return fileBackend
	}()
	consoleFormat := logging.MustStringFormatter("[%{time:01/02 15:04:05}] %{message}")
	consoleBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(consoleLogBackend, consoleFormat))
	level, err := logging.LogLevel(global.config.Log.ConsoleLevel)
	if err != nil {
		panic(err)
	}
	consoleBackendLevel.SetLevel(level, "")
	global.syslog = logging.MustGetLogger("globalsys")
	global.syslog.SetBackend(consoleBackendLevel)

	var consoleFormatString string
	var fileFormatString string
	// switch frame.config.RunMode {
	// case RUNMODE_DEV:
	consoleFormatString = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{shortfile}>"
	fileFormatString = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{shortfile}>"
	// case RUNMODE_PROD:
	// consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{shortfile}>"
	// fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{module} #%{shortfile}>"
	// }
	global.bizlog = global.newLogger(
		"globalbiz",
		consoleFormatString,
		fileFormatString,
	)
	global.bizlog.ExtraCalldepth++
}

func (frame *Framework) initSysLogger() {
	var consoleFormat string
	var fileFormat string
	// switch frame.config.RunMode {
	// case RUNMODE_DEV:
	consoleFormat = "[%{time:01/02 15:04:05}] \x1b[46m[SYS]\x1b[0m %{message} <%{module}>"
	fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [SYS] %{message} <%{module}>"
	// case RUNMODE_PROD:
	// consoleFormat = "[%{time:01/02 15:04:05}] \x1b[46m[SYS]\x1b[0m %{message} <%{module}>"
	// fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [SYS] %{message} <%{module}>"
	// }
	frame.syslog = Global.newLogger(
		strings.ToLower(frame.NameWithVersion()),
		consoleFormat,
		fileFormat,
	)
}

func (frame *Framework) initBizLogger() {
	var consoleFormat string
	var fileFormat string
	// switch frame.config.RunMode {
	// case RUNMODE_DEV:
	consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{shortfile}>"
	fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{module} #%{shortfile}>"
	// case RUNMODE_PROD:
	// consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{shortfile}>"
	// fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{module} #%{shortfile}>"
	// }
	frame.bizlog = Global.newLogger(
		strings.ToLower(frame.NameWithVersion()),
		consoleFormat,
		fileFormat,
	)
}

func (global *GlobalSetting) newLogger(module string, consoleFormatString, fileFormatString string) *logging.Logger {
	consoleLevel, err := logging.LogLevel(global.config.Log.ConsoleLevel)
	if err != nil {
		panic(err)
	}
	fileLevel, err := logging.LogLevel(global.config.Log.FileLevel)
	if err != nil {
		panic(err)
	}
	consoleFormat := logging.MustStringFormatter(consoleFormatString)
	fileFormat := logging.MustStringFormatter(fileFormatString)
	backends := []logging.Backend{}

	if global.config.Log.ConsoleEnable {
		consoleBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(consoleLogBackend, consoleFormat))
		consoleBackendLevel.SetLevel(consoleLevel, "")
		backends = append(backends, consoleBackendLevel)
	}

	if global.config.Log.FileEnable {
		fileBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(fileBackend, fileFormat))
		fileBackendLevel.SetLevel(fileLevel, "")
		backends = append(backends, fileBackendLevel)
	}

	newLog := logging.MustGetLogger(module)
	switch len(backends) {
	case 1:
		newLog.SetBackend(backends[0].(logging.LeveledBackend))
	default:
		newLog.SetBackend(logging.MultiLogger(backends...))
	}
	return newLog
}

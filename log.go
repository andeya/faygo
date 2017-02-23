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

package faygo

import (
	"log"
	"os"
	"strings"

	"github.com/henrylee2cn/faygo/logging"
	"github.com/henrylee2cn/faygo/logging/color"
)

// NewLog gets a global logger
func NewLog() *logging.Logger {
	newlog := *global.bizlog
	newlog.ExtraCalldepth--
	return &newlog
}

var (
	consoleLogBackend = &logging.LogBackend{
		Logger: log.New(color.NewColorableStdout(), "", 0),
		Color:  true,
	}
	fileBackend *logging.FileBackend
)

func (global *GlobalVariables) initLogger() {
	if global.config.Log.FileEnable {
		fileBackend = func() *logging.FileBackend {
			fileBackend, err := logging.NewDefaultFileBackend(global.logDir+"faygo.log", global.config.Log.AsyncLen)
			if err != nil {
				panic(err)
			}
			return fileBackend
		}()
	} else {
		os.MkdirAll(global.logDir, 0777)
	}
	consoleFormat := logging.MustStringFormatter("[%{time:01/02 15:04:05}] %{message}")
	consoleBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(consoleLogBackend, consoleFormat))
	level, err := logging.LogLevel(global.config.Log.ConsoleLevel)
	if err != nil {
		panic(err)
	}
	consoleBackendLevel.SetLevel(level, "")
	global.syslog = logging.NewLogger("globalsys")
	global.syslog.SetBackend(consoleBackendLevel)

	var consoleFormatString string
	var fileFormatString string
	// switch frame.config.RunMode {
	// case RUNMODE_DEV:
	consoleFormatString = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{longfile}>"
	fileFormatString = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{longfile}>"
	// case RUNMODE_PROD:
	// consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{longfile}>"
	// fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{module} #%{longfile}>"
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
	consoleFormat = "[%{time:01/02 15:04:05}] %{message} <%{module}>"
	fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] %{message} <%{module}>"
	// case RUNMODE_PROD:
	// consoleFormat = "[%{time:01/02 15:04:05}] \x1b[46m[SYS]\x1b[0m %{message} <%{module}>"
	// fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [SYS] %{message} <%{module}>"
	// }
	frame.syslog = global.newLogger(
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
	consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{longfile}>"
	fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{module} #%{longfile}>"
	// case RUNMODE_PROD:
	// consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{longfile}>"
	// fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [%{level:.1s}] %{message} <%{module} #%{longfile}>"
	// }
	frame.bizlog = global.newLogger(
		strings.ToLower(frame.NameWithVersion()),
		consoleFormat,
		fileFormat,
	)
}

func (global *GlobalVariables) newLogger(module string, consoleFormatString, fileFormatString string) *logging.Logger {
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

	newLog := logging.NewLogger(module)
	switch len(backends) {
	case 1:
		newLog.SetBackend(backends[0].(logging.LeveledBackend))
	default:
		newLog.SetBackend(logging.MultiLogger(backends...))
	}
	return newLog
}

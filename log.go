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
	"sync"

	"github.com/henrylee2cn/thinkgo/logging"
	"github.com/henrylee2cn/thinkgo/logging/color"
)

var (
	consoleLogBackend = &logging.LogBackend{
		Logger: log.New(color.NewColorableStdout(), "", 0),
		Color:  true,
	}
	fileBackend = func() *logging.FileBackend {
		fileBackend, err := logging.NewDefaultFileBackend(LOG_DIR + "/thinkgo.log")
		if err != nil {
			panic(err)
		}
		return fileBackend
	}()
	globalSysLogger     *logging.Logger
	globalSysLoggerOnce sync.Once
)

func initGlobalSysLogger(config LogConfig) {
	globalSysLoggerOnce.Do(func() {
		consoleFormat := logging.MustStringFormatter("[%{time:01/02 15:04:05}] %{message}")
		consoleBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(consoleLogBackend, consoleFormat))
		level, err := logging.LogLevel("warning")
		if err != nil {
			panic(err)
		}
		consoleBackendLevel.SetLevel(level, "")
		globalSysLogger = logging.MustGetLogger("global")
		globalSysLogger.SetBackend(consoleBackendLevel)
	})
}

func (frame *Framework) initSysLogger() {
	initGlobalSysLogger(frame.config.Log)
	var consoleFormat string
	var fileFormat string
	switch frame.config.RunMode {
	case RUNMODE_DEV:
		consoleFormat = "[%{time:01/02 15:04:05}] \x1b[46m[SYS]\x1b[0m %{message} <%{module}>"
		fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [SYS] %{message} <%{module}>"
	case RUNMODE_PROD:
		consoleFormat = "[%{time:01/02 15:04:05}] \x1b[46m[SYS]\x1b[0m %{message} <%{module}>"
		fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [SYS] %{message} <%{module}>"
	}
	frame.syslog = frame.newLogger(
		strings.ToLower(frame.NameWithVersion()),
		consoleFormat,
		fileFormat,
	)
}

func (frame *Framework) initBizLogger() {
	var consoleFormat string
	var fileFormat string
	switch frame.config.RunMode {
	case RUNMODE_DEV:
		consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{shortfile}>"
		fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{shortfile}>"
	case RUNMODE_PROD:
		consoleFormat = "[%{time:01/02 15:04:05}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{shortfile}>"
		fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] %{color}[%{level:.1s}]%{color:reset} %{message} <%{module} #%{shortfile}>"
	}
	frame.bizlog = frame.newLogger(
		strings.ToLower(frame.NameWithVersion()),
		consoleFormat,
		fileFormat,
	)
}

func (frame *Framework) newLogger(module string, consoleFormatString, fileFormatString string) *logging.Logger {
	consoleLevel, err := logging.LogLevel(frame.config.Log.ConsoleLevel)
	if err != nil {
		panic(err)
	}
	fileLevel, err := logging.LogLevel(frame.config.Log.FileLevel)
	if err != nil {
		panic(err)
	}
	consoleFormat := logging.MustStringFormatter(consoleFormatString)
	fileFormat := logging.MustStringFormatter(fileFormatString)
	backends := []logging.Backend{}

	if frame.config.Log.ConsoleEnable {
		consoleBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(consoleLogBackend, consoleFormat))
		consoleBackendLevel.SetLevel(consoleLevel, "")
		backends = append(backends, consoleBackendLevel)
	}

	if frame.config.Log.FileEnable {
		fileBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(fileBackend, fileFormat))
		fileBackendLevel.SetLevel(fileLevel, "")
		backends = append(backends, fileBackendLevel)
	}

	newLog := logging.MustGetLogger(module)
	switch len(backends) {
	case 1:
		newLog.SetBackend(backends[0].(logging.LeveledBackend))
	case 2:
		newLog.SetBackend(logging.MultiLogger(backends...))
	default:
		globalSysLogger.Warning("config: log::enable_console and log::enable_file can not be disabled at the same time, so automatically open console log.")
		frame.config.Log.ConsoleEnable = true
		configFileName, _ := createConfigFilenameAndVersion(frame.name, frame.version)
		err := syncConfigToFile(configFileName, &frame.config)
		if err != nil {
			globalSysLogger.Critical("[C] config: log::enable_console and log::enable_file must be at least one for true.")
			return globalSysLogger
		}
		return frame.newLogger(module, consoleFormatString, fileFormatString)
	}
	return newLog
}

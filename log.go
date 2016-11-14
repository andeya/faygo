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

func (frame *Framework) initSysLogger() {
	var consoleFormat string
	var fileFormat string
	switch frame.config.RunMode {
	case RUNMODE_DEV:
		consoleFormat = "[%{time:01/02 15:04:05}] \x1b[46m[SYS]\x1b[0m %{message} <%{module} #%{shortfile}>"
		fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [SYS] %{message} <%{module} #%{shortfile}>"
	case RUNMODE_PROD:
		consoleFormat = "[%{time:01/02 15:04:05}] \x1b[46m[SYS]\x1b[0m %{message} <%{module} #%{shortfile}>"
		fileFormat = "[%{time:2006/01/02T15:04:05.999Z07:00}] [SYS] %{message} <%{module} #%{shortfile}>"
	}
	frame.syslog = newLogger(
		strings.TrimSuffix(strings.ToLower(frame.name+"_"+frame.version), "_"),
		frame.config.Log,
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
	frame.bizlog = newLogger(
		strings.TrimSuffix(strings.ToLower(frame.name+"_"+frame.version), "_"),
		frame.config.Log,
		consoleFormat,
		fileFormat,
	)
}

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
)

func newLogger(module string, config LogConfig, consoleFormatString, fileFormatString string) *logging.Logger {
	consoleLevel, err := logging.LogLevel(config.ConsoleLevel)
	if err != nil {
		panic(err)
	}
	fileLevel, err := logging.LogLevel(config.FileLevel)
	if err != nil {
		panic(err)
	}
	consoleFormat := logging.MustStringFormatter(consoleFormatString)
	fileFormat := logging.MustStringFormatter(fileFormatString)
	backends := []logging.Backend{}

	if config.ConsoleEnable {
		consoleBackendLevel := logging.AddModuleLevel(logging.NewBackendFormatter(consoleLogBackend, consoleFormat))
		consoleBackendLevel.SetLevel(consoleLevel, "")
		backends = append(backends, consoleBackendLevel)
	}

	if config.FileEnable {
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
		panic("[config] log::enable_console and log::enable_file must be at least one for true")
	}
	return newLog
}

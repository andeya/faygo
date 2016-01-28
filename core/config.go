// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/henrylee2cn/thinkgo/core/config"
	"github.com/henrylee2cn/thinkgo/core/log"
)

type Config struct {
	AppName       string // 应用名称
	Debug         bool   // 是否开启调试模式
	LogLevel      log.Level
	HttpAddr      string // 应用监听地址，默认为空，监听所有的网卡 IP
	HttpPort      int    // 应用监听端口，默认为 8080
	TplSuffix     string // 模板后缀名
	TplLeft       string // 模板左定界符
	TplRight      string // 模板右定界符
	DefaultModule string // 默认模块的名称
}

func getConfig() Config {
	iniconf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("\n  请确保在项目目录下运行，且存在配置文件 conf/app.conf")
		os.Exit(1)
	}

	var logLevel log.Level
	switch strings.ToUpper(iniconf.String("loglevel")) {
	case "TRACE":
		logLevel = log.TRACE
	case "DEBUG":
		logLevel = log.DEBUG
	case "INFO":
		logLevel = log.INFO
	case "NOTICE":
		logLevel = log.NOTICE
	case "WARN":
		logLevel = log.WARN
	case "ERROR":
		logLevel = log.ERROR
	case "FATAL":
		logLevel = log.FATAL
	case "OFF":
		logLevel = log.OFF
	default:
		logLevel = log.DEBUG
	}
	defaultModule := iniconf.DefaultString("defmodule", "home")
	return Config{
		AppName:       iniconf.DefaultString("appname", "thinkgo"),
		Debug:         iniconf.DefaultBool("debug", true),
		LogLevel:      logLevel,
		HttpAddr:      iniconf.DefaultString("httpaddr", "0.0.0.0"),
		HttpPort:      iniconf.DefaultInt("httpport", 8080),
		TplSuffix:     iniconf.DefaultString("tplsuffex", ".html"),
		TplLeft:       iniconf.DefaultString("tplleft", "{{{"),
		TplRight:      iniconf.DefaultString("tplright", "}}}"),
		DefaultModule: SnakeString(strings.Trim(defaultModule, "/")),
	}
}

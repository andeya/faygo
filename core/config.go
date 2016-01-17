// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package core

import (
	"github.com/henrylee2cn/thinkgo/core/config"
	"log"
)

type Config struct {
	AppName   string // 应用名称
	RunMode   string // 运行模式 "release"/"debug"/"test"
	HttpAddr  string // 应用监听地址，默认为空，监听所有的网卡 IP
	HttpPort  int    // 应用监听端口，默认为 8080
	TplSuffex string // 模板后缀名
	TplLeft   string // 模板左定界符
	TplRight  string // 模板右定界符
}

func readConfig() Config {
	iniconf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		log.Panicln(err)
	}
	return Config{
		AppName:   iniconf.DefaultString("appname", "thinkgo"),
		RunMode:   iniconf.DefaultString("runmode", "debug"),
		HttpAddr:  iniconf.DefaultString("httpaddr", "0.0.0.0"),
		HttpPort:  iniconf.DefaultInt("httpport", 8080),
		TplSuffex: iniconf.DefaultString("tplsuffex", ".html"),
		TplLeft:   iniconf.DefaultString("tplleft", "{{{"),
		TplRight:  iniconf.DefaultString("tplright", "}}}"),
	}
}

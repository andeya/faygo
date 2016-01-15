package conf

import (
	"github.com/henrylee2cn/thinkgo/core"
)

const (
	NAME        = "MoreChat 管理平台"
	CLASS       = "系统框架"
	DESCRIPTION = "MoreChat 管理平台，采用插件式架构，开发者可以进行轻松扩展。"
)

var BASE_DATA = core.H{
	"TITLE":                NAME,
	"PAGE_HEADER":          "Admin",
	"OPTIONAL_DESCRIPTION": "MoreChat 管理平台首页",
	"USER_NAME":            "Guest",
	"ROLE":                 "Guest",
	"LEVEL": map[string]string{
		"K": "Home",
		"V": "javascript:;",
	},
	"HERE":           "Admin",
	"NAME":           core.App.Name,
	"VERSION":        core.App.Version,
	"COMPANY":        core.App.Company,
	"COPYRIGHT_YEAR": core.App.CopyrightYear,
	"CONTACT":        core.App.Contact,
	"DESCRIPTION":    core.App.Description,
}

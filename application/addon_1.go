package application

import (
	"github.com/henrylee2cn/thinkgo/application/addon_1/conf"
	. "github.com/henrylee2cn/thinkgo/application/addon_1/controller"
	. "github.com/henrylee2cn/thinkgo/core"
)

func init() {
	AddonPrepare(&Module{
		Name:        conf.NAME,
		Class:       conf.CLASS,
		Description: conf.DESCRIPTION,
	}).SetThemes(
		// 自动设置传入的第1个主题为当前主题
		&Theme{
			Name:        "default",
			Description: "default",
			Src:         map[string]string{},
		},
	).
		// 指定当前主题
		UseTheme("default").
		// 中间件
		// Use().
		// 注册路由
		GET("/index", &IndexController{}).
		GET("/compose", &IndexController{}).
		GET("/ReadMail/:id", &IndexController{}) // 将被转为read_mail/:id
}

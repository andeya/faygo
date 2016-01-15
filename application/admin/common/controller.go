package common

import (
	"github.com/henrylee2cn/thinkgo/application/admin/conf"
	"github.com/henrylee2cn/thinkgo/core"
)

type BaseController struct {
	core.BaseController
}

func (this *BaseController) HTML(code ...int) {
	for k, v := range conf.BASE_DATA {
		this.Data.TrySet(k, v)
	}
	this.Data["ADDONS"] = core.App.Addons.Slice
	this.Data.TrySet("USER_IMAGE", this.Data["__PUBLIC__"].(string)+"/img/user2-160x160.jpg")
	this.BaseController.HTML(code...)
}

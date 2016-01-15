package common

import (
	"github.com/henrylee2cn/thinkgo/application/addon_1/conf"
	"github.com/henrylee2cn/thinkgo/core"
)

type BaseController struct {
	core.BaseController
}

func (this *BaseController) HTML(code ...int) {
	for k, v := range conf.BASE_DATA {
		this.Data.TrySet(k, v)
	}
	this.BaseController.HTML(code...)
}

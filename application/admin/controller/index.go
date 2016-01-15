package controller

import (
	"github.com/henrylee2cn/thinkgo/application/admin/common"
)

type IndexController struct {
	common.BaseController
}

func (this *IndexController) Index() {
	id := this.Query("addon")
	if id == "" {
		id = this.Param("addon")
	}
	var iFrame string
	if id != "" {
		iFrame = "/" + id + "/index/index"
	}
	this.Data["iFrame"] = iFrame
	this.HTML()
}

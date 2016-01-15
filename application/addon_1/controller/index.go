package controller

import (
	"github.com/henrylee2cn/thinkgo/application/addon_1/common"
)

type IndexController struct {
	common.BaseController
}

func (this *IndexController) Index() {
	this.HTML()
}

func (this *IndexController) Compose() {
	this.HTML()
}

func (this *IndexController) ReadMail() {
	this.HTML()
}

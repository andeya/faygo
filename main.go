package main

import (
	_ "github.com/henrylee2cn/thinkgo/application"
	_ "github.com/henrylee2cn/thinkgo/common/deploy"
	"github.com/henrylee2cn/thinkgo/core"
)

func main() {
	core.ThinkGoDefault().Run(":8080")
}

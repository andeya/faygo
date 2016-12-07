package middleware

import (
	// "errors"
	"github.com/henrylee2cn/thinkgo"
)

func Root2Index(ctx *thinkgo.Context) error {
	if ctx.Path() == "/index" {
		ctx.Stop()
		return nil
		// return errors.New("Please access the root directory `/`")
	}
	if ctx.Path() == "/" {
		ctx.ModifyPath("/index")
	}
	return nil
}

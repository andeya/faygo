package middleware

import (
	// "net/http"

	"github.com/henrylee2cn/faygo"
)

func Root2Index(ctx *faygo.Context) error {
	// Direct access to `/index` is not allowed
	if ctx.Path() == "/index" {
		ctx.Stop()
		// ctx.Error(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return nil
	}

	if ctx.Path() == "/" {
		ctx.ModifyPath("/index")
	}
	return nil
}

/**
 * desc : 登录后起始页
 * author:畅雨
 * date:  2016.05.16
 * history:
 *
 */

package handler

import (
	"github.com/henrylee2cn/faygo"
)

func Index() faygo.HandlerFunc {
	return func(c *faygo.Context) error {
		return c.Render(200, "view/index.html", nil)
	}
}
func Pongo2() faygo.HandlerFunc {
	return func(c *faygo.Context) error {
		return c.Render(200, "view/pongo2.tpl", nil)
	}
}

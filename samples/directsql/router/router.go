/**
 * desc : 系统功能路由注册
 * author:畅雨
 * date:  2016.05.16
 * history:
 *
 */
package router

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/ext/db/directsql"
	"github.com/henrylee2cn/thinkgo/samples/directsql/handler"
)

// Register the system router in a chain style
func Route(frame *thinkgo.Framework) {
	thinkgo.SetUpload("./upload", false, false)
	thinkgo.SetStatic("./static", false, false)
	// Home page
	frame.NamedAPI("Home", "GET", "/", handler.Index())
	frame.NamedAPI("Pongo2", "GET", "/pongo2", handler.Pongo2())
	// bos 执行SQL定义的路由
	frame.NamedAPI("DirectSQL", "POST", "/bos/*path", directsql.DirectSQL())
	frame.NamedGET("DirectSQL ModelSql Reload", "/bom/reloadall", directsql.DirectSQLReloadAll())
	frame.NamedGET("DirectSQL ModelSql Reload", "/bom/reload/*path", directsql.DirectSQLReloadModel())
	/*frame.NamedStaticFS("render ", "/tpl", thinkgo.RenderFS(
		"./view",
		".tpl", // "*"
		thinkgo.Map{"title": "tpl page"},
	))*/

}

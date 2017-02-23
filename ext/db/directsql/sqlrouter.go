/*
  动态SQL路由注册
  特别说明：将router定义代码放到 sysrouter.go中
*/
package directsql

/*
import (
	"github.com/henrylee2cn/faygo"
)

//注册路由
func SysRoute(frame *faygo.Framework) {
	frame.Static("/sys", "./sys_view")

	file := frame.NamedGroup("文件管理", "/file", middleware.CheckLogin())
	{
		file.NamedPOST("文件上传", "/upload", service.FileUpload())
		// file.NamedPOST("文件上传加强版", "/uploadpro", service.FileUploadPro())
		file.NamedPOST("文件删除", "/delete", &service.FileDelete{})
		// file.NamedPOST("文件删除加强版", "/deletepro", &service.FileDeletePro{})
	}
	_bos := frame.NamedGroup("后台管理", "/bos", middleware.CheckGotoLogin())
	{
		_bos.NamedAPI("DirectSQL", "REST", "/*path", directsql.DirectSQL())
	}
	// bos 执行SQL定义的路由
	frame.NamedAPI("DirectSQL", "GET POST", "/bos/*path", directsql.DirectSQL())                       //.Use(middleware.CheckLogin())
	frame.NamedGET("DirectSQL ModelSql Reload", "/bom/reloadall", directsql.DirectSQLReloadAll())      //.Use(middleware.CheckLogin())
	frame.NamedGET("DirectSQL ModelSql Reload", "/bom/reload/*path", directsql.DirectSQLReloadModel()) //.Use(middleware.CheckLogin())

	_admin := frame.NamedGroup("后台管理", "/admin", middleware.CheckGotoLogin())
	{
		_admin.NamedAPI("后台首页", "REST", "/index", admin.Index())
	}
}
}*/

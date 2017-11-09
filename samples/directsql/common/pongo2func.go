/*
 功能：pongo2模板数据库访问函数
 日期：2017.01.06

*/
package common

import (
	"encoding/json"
	"strings"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/directsql"
	"github.com/henrylee2cn/faygo/pongo2"
)

/* tpl
   {% for o in SimpleData("biz/demo","selecttpl","{`para`:`我是sql参数值`}") %}
       <p>名称:{{ o.cnname }}</p>
   {% endfor %}

*/
//单个查询 参数 map[string]interface{} 返回 []map[string]interface{}
func SimpleData(modelId, sqlId string, mp string) *pongo2.Value {
	//参数处理
	//faygo.Debug("SimpleData mp:", mp)
	para := make(map[string]interface{})
	if err := json.Unmarshal([]byte(strings.Replace(mp, "`", `"`, -1)), &para); err != nil {
		faygo.Error(err.Error())
		return pongo2.AsValue(err.Error())
	}
	faygo.Debug("SimpleData para:", para)
	//执行sql获取结果
	result, err := directsql.SelectMapToMap(modelId, sqlId, para)
	if err != nil {
		faygo.Error(err.Error())
		return pongo2.AsValue(err)
	}
	faygo.Debug("SimpleData result :", result)
	return pongo2.AsValue(result)
}

//单个查询 参数 map[string]interface{} 返回 []map[string]interface{}
func SimpleData2(modelId, sqlId string, mp string) []map[string]interface{} {
	//参数处理
	//faygo.Debug("SimpleData mp:", mp)
	result := make([]map[string]interface{}, 0)
	para := make(map[string]interface{})
	if err := json.Unmarshal([]byte(strings.Replace(mp, "`", `"`, -1)), &para); err != nil {
		faygo.Error(err.Error())
		//result=append(result,err)
	}
	faygo.Debug("SimpleData para:", para)
	//执行sql获取结果
	result, err := directsql.SelectMapToMap(modelId, sqlId, para)
	if err != nil {
		faygo.Error(err.Error())
		//result=append(result,err.Error())
	}
	faygo.Debug("SimpleData result :", result)
	return result
}

/* tpl
   {{ GetName("李四ssssss","张三","王五") }}
*/

func Test(str1, str2, str3 string) string {
	return "我们是" + str1 + str2 + str3
}
func init() {
	// 测试函数
	faygo.RenderVar("Test", Test)
	//获取单表数据 返回 []map[string]interface{}
	faygo.RenderVar("SimpleData", SimpleData)
	faygo.RenderVar("SimpleData2", SimpleData2)
}

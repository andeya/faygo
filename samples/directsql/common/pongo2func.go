/*
 功能：pongo2模板数据库访问函数
 日期：2017.01.06

*/
package common

import (
	"encoding/json"
	"strings"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/ext/db/directsql"
	"github.com/henrylee2cn/thinkgo/pongo2"
)

/* tpl
   {% for o in SimpleData("biz/demo","selecttpl","{`para`:`我是sql参数值`}") %}
       <p>名称:{{ o.cnname }}</p>
   {% endfor %}

*/
//单个查询 参数 map[string]interface{} 返回 []map[string]interface{}
func SimpleData(modelId, sqlId string, mp string) *pongo2.Value {
	//参数处理
	//thinkgo.Debug("SimpleData mp:", mp)
	para := make(map[string]interface{})
	if err := json.Unmarshal([]byte(strings.Replace(mp, "`", `"`, -1)), &para); err != nil {
		thinkgo.Error(err.Error())
		return pongo2.AsValue(err.Error())
	}
	thinkgo.Debug("SimpleData para:", para)
	//执行sql获取结果
	result, err := directsql.SelectMapToMap(modelId, sqlId, para)
	if err != nil {
		thinkgo.Error(err.Error())
		return pongo2.AsValue(err)
	}
	thinkgo.Debug("SimpleData result :", result)
	return pongo2.AsValue(result)
}

//单个查询 参数 map[string]interface{} 返回 []map[string]interface{}
func SimpleData2(modelId, sqlId string, mp string) []map[string]interface{} {
	//参数处理
	//thinkgo.Debug("SimpleData mp:", mp)
	result := make([]map[string]interface{}, 0)
	para := make(map[string]interface{})
	if err := json.Unmarshal([]byte(strings.Replace(mp, "`", `"`, -1)), &para); err != nil {
		thinkgo.Error(err.Error())
		//result=append(result,err)
	}
	thinkgo.Debug("SimpleData para:", para)
	//执行sql获取结果
	result, err := directsql.SelectMapToMap(modelId, sqlId, para)
	if err != nil {
		thinkgo.Error(err.Error())
		//result=append(result,err.Error())
	}
	thinkgo.Debug("SimpleData result :", result)
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
	thinkgo.RenderVar("Test", Test)
	//获取单表数据 返回 []map[string]interface{}
	thinkgo.RenderVar("SimpleData", SimpleData)
	thinkgo.RenderVar("SimpleData2", SimpleData2)
}

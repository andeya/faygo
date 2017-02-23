/**
* desc : 请求处理handle
* author:畅雨
* date:  2016.10.27
* history:
         2016.10.17：优化 select与pagingselect返回结果，当数据为空时返回[]
*
*/
package directsql

import (
	"bytes"
	"encoding/json"
	"html/template"

	"github.com/henrylee2cn/faygo"
)

//DirectSQL handler 定义
func DirectSQL() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		//1.根据路径获取sqlentity:去掉/bos/,再拆分成 modelId，sqlId
		modelId, sqlId := trimBeforeSplitRight(ctx.Path(), '/', 2)
		faygo.Debug("Model file: " + modelId + "  - sqlId:" + sqlId)
		//2.获取ModelSql
		m := findModel(modelId)
		if m == nil {
			faygo.Error("Error: model file does not exist,") //("错误：未定义的Model文件: " + modelId)
			return ctx.JSONMsg(404, 404, "Error:model file does not exist: "+modelId)
		}
		//3.获取Sql
		se := m.findSql(sqlId)
		if se == nil { //
			faygo.Error("Error: sql is not defined in the model file, " + modelId + "/" + sqlId) //错误：Model文件中未定义sql:
			return ctx.JSONMsg(404, 404, "Error: sql is not defined in the model file, "+modelId+"/"+sqlId)
		}
		//4.根据SQL类型分别处理执行并返回结果信息
		switch se.Sqltype {
		case ST_PAGINGSELECT: //分页选择SQL，分頁查詢結果cache第一次查询的结果
			//.1 获取POST参数並轉換
			var jsonpara map[string]interface{}
			err := ctx.BindJSON(&jsonpara) //从Body获取JSON参数
			if err != nil {
				faygo.Debug("Info:POST para is empty," + err.Error())
				//如果参数为空则会触发EOF错误,不应该退出因为可能本来就没有参数，也就是jsonpara仍旧为空，需要创建该变量，后续sql中的参数处理需要
				if jsonpara == nil {
					jsonpara = make(map[string]interface{})
				}
			}
			//. 常規參數處理
			var callback string
			if v, ok := jsonpara["callback"]; ok {
				s, ok := v.(string)
				if ok {
					callback = s
					delete(jsonpara, "callback")
				}
			}
			//.2 判断是否是缓存的并存在有效缓存，存在则直接从缓存返回
			if se.Cached {
				//构造缓存查询key
				key := modelId + "/" + sqlId
				//缓存识别的后缀
				sf, err := json.Marshal(jsonpara)
				if err != nil {
					sf = nil
				}
				suffix := string(sf)
				//如果OK则直接返回缓存
				if ok, jsonb := GetCache(key, suffix); ok {
					faygo.Debug("Directsql getCache:[" + key + " - " + suffix + "] result from cache.")
					//发送JSON(P)响应
					return sendJSON(ctx, callback, jsonb)
				}
			}
			//.3 检查sql语句配置个数
			if len(se.Cmds) != 2 {
				faygo.Error("Error: paging query must define two sql nodes, one for total number and one for data query!") //错误：分页查询必须定义2个SQL节点，一个获取总页数另一个用于查询数据！
				return ctx.JSONMsg(404, 404, "Error: paging query must define two sql nodes, one for total number and one for data query!")
			}

			//.5 参数验证并处理(参数定义到真正查询结果的cmd下)－OK
			_, err = dealwithParameter(se.Cmds[1].Parameters, jsonpara, ctx)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(400, 400, err.Error())
			}
			//.6 執行並返回結果
			data, err := m.pagingSelectMap(se, jsonpara)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//.7 如果需要缓存则缓存结果集(cached=true 并且缓存不存在或失效才会执行)
			jsonb, err := intface2json(data)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//结果集为空响应
			if data.Total == 0 {
				err := ctx.JSONBlob(200, []byte(`{"total":0,"data":[]}`))
				if err != nil {
					return err
				}
				return nil
			}
			//如果需要缓存则
			if se.Cached {
				//构造缓存查询key
				key := modelId + "/" + sqlId
				//缓存识别的后缀
				sf, err := json.Marshal(jsonpara)
				if err != nil {
					sf = nil
				}
				suffix := string(sf)
				SetCache(key, suffix, jsonb, se.Cachetime)
				faygo.Debug("Directsql setCache:[" + key + "] result to cache.")
			}
			//发送JSON(P)响应
			return sendJSON(ctx, callback, jsonb)

			//一般选择SQL,嵌套选择暂时未实现跟一般选择一样,增强的插入后返回sysid ----OK
		case ST_SELECT, ST_NESTEDSELECT:
			//.1 获取POST参数並轉換
			var jsonpara map[string]interface{}
			err := ctx.BindJSON(&jsonpara) //从Body获取JSON参数
			if err != nil {
				faygo.Debug("Info: POST para is empty," + err.Error())
				//如果参数为空则会触发EOF错误,不应该退出因为可能本来就没有参数，也就是jsonpara仍旧为空，需要创建该变量，后续sql中的参数处理需要
				if jsonpara == nil {
					jsonpara = make(map[string]interface{})
				}
			}
			//. 常規參數處理
			var callback string
			if v, ok := jsonpara["callback"]; ok {
				s, ok := v.(string)
				if ok {
					callback = s
					delete(jsonpara, "callback")
				}
			}
			//.2 判断是否是缓存的并存在有效缓存，存在则直接从缓存返回
			if se.Cached {
				//构造缓存查询key
				key := modelId + "/" + sqlId
				//缓存识别的后缀
				sf, err := json.Marshal(jsonpara)
				if err != nil {
					sf = nil
				}
				suffix := string(sf)
				//如果OK则直接返回缓存
				if ok, jsonb := GetCache(key, suffix); ok {
					faygo.Debug("Directsql getCache:[" + key + "] result from cache.")
					//发送JSON(P)响应
					return sendJSON(ctx, callback, jsonb)
				}
			}
			//.3 参数验证并处理,－OK
			_, err = dealwithParameter(se.Cmds[0].Parameters, jsonpara, ctx)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(400, 400, err.Error())
			}
			//.4 執行並返回結果
			data, err := m.selectMap(se, jsonpara)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//.7 如果需要缓存则缓存结果集(cached=true 并且缓存不存在或失效才会执行)
			jsonb, err := intface2json(data)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//结果集为空响应
			if len(data) == 0 {
				err := ctx.JSONBlob(200, []byte(`[]`))
				if err != nil {
					return err
				}
				return nil
			}
			//如果需要缓存则
			if se.Cached {
				//构造缓存查询key
				key := modelId + "/" + sqlId
				//缓存识别的后缀
				sf, err := json.Marshal(jsonpara)
				if err != nil {
					sf = nil
				}
				suffix := string(sf)
				SetCache(key, suffix, jsonb, se.Cachetime)
				faygo.Debug("Directsql setCache:[" + key + "] result to cache.")
			}
			//发送JSON(P)响应
			return sendJSON(ctx, callback, jsonb)
			//如果没有结果集则输出[]
			/*if len(data) == 0 {
				_, err := ctx.Write([]byte("[]"))
				if err != nil {
					return err
				}
				return nil
			}
			//
			if len(callback) > 0 {
				return ctx.JSONP(200, callback, data)
			}
			return ctx.JSON(200, data)*/

		case ST_MULTISELECT: //返回多結果集選擇
			//.1 获取POST参数並轉換
			var jsonpara map[string]interface{}
			err := ctx.BindJSON(&jsonpara) //从Body获取JSON参数
			if err != nil {
				faygo.Info("Info:POST para is empty," + err.Error())
				//return ctx.JSONMsg(404, 404, err.Error())
				//如果参数为空则会触发EOF错误,不应该退出因为可能本来就没有参数，也就是jsonpara仍旧为空，需要创建该变量，后续sql中的参数处理需要
				if jsonpara == nil {
					jsonpara = make(map[string]interface{})
				}
			}
			//. 常規參數處理
			var callback string
			if v, ok := jsonpara["callback"]; ok {
				s, ok := v.(string)
				if ok {
					callback = s
					delete(jsonpara, "callback")
				}
			}
			//.2 判断是否是缓存的并存在有效缓存，存在则直接从缓存返回
			if se.Cached {
				//构造缓存查询key
				key := modelId + "/" + sqlId
				//缓存识别的后缀
				sf, err := json.Marshal(jsonpara)
				if err != nil {
					sf = nil
				}
				suffix := string(sf)
				//如果OK则直接返回缓存
				if ok, jsonb := GetCache(key, suffix); ok {
					faygo.Debug("GetCache:[" + key + " - " + suffix + "] result from cache.")
					//发送JSON(P)响应
					return sendJSON(ctx, callback, jsonb)
				}
			}
			//.3 参数验证并处理－OK
			for _, cmd := range se.Cmds {
				//未配置参数则直接忽略
				if len(cmd.Parameters) == 0 {
					continue
				}
				_, err = dealwithParameter(cmd.Parameters, jsonpara, ctx)
				if err != nil {
					faygo.Error(err.Error())
					return ctx.JSONMsg(400, 400, err.Error())
				}
			}
			//.4 執行並返回結果
			data, err := m.multiSelectMap(se, jsonpara)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//.7 如果需要缓存则缓存结果集(cached=true 并且缓存不存在或失效才会执行)
			jsonb, err := intface2json(data)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//结果集为空响应
			if len(data) == 0 {
				err := ctx.JSONBlob(200, []byte(`[]`))
				if err != nil {
					return err
				}
				return nil
			}
			//如果需要缓存则
			if se.Cached {
				//构造缓存查询key
				key := modelId + "/" + sqlId
				//缓存识别的后缀
				sf, err := json.Marshal(jsonpara)
				if err != nil {
					sf = nil
				}
				suffix := string(sf)
				SetCache(key, suffix, jsonb, se.Cachetime)
				faygo.Debug("Directsql setCache:[" + key + " - " + suffix + "] result to cache.")
			}
			//发送JSON(P)响应
			return sendJSON(ctx, callback, jsonb)

		case ST_EXEC: //执行SQL(插入、删除、更新sql)
			//.1.获取 Ajax post json 参数
			var jsonpara map[string]interface{}
			err := ctx.BindJSON(&jsonpara) //从Body获取 json参数
			if err != nil {
				faygo.Info("Info: POST para is empty," + err.Error())
				//return ctx.JSONMsg(404, 404, err.Error())
				//如果参数为空则会触发EOF错误,不应该退出因为可能本来就没有参数，也就是jsonpara仍旧为空，需要创建该变量，后续sql中的参数处理需要
				if jsonpara == nil {
					jsonpara = make(map[string]interface{})
				}
			}
			//.2.SQL定义参数验证并处理---OK，服务端生成的uuid返回给客户端
			result, err := dealwithParameter(se.Cmds[0].Parameters, jsonpara, ctx)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(400, 400, err.Error())
			}
			//.3.执行sql
			err = m.execMap(se, jsonpara)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//return ctx.JSON(200, result)
			//如果存在服务端生成的uuid参数的则返回到客户端
			if (result != nil) && (len(result) > 0) {
				return ctx.JSONMsg(200, 200, result)
			} else {
				return ctx.JSONMsg(200, 200, "Info: Exec sql ok!")
			}

		case ST_BATCHEXEC: //批量执行--原来的批量插入
			//.1.获取 Ajax post json 参数
			var jsonpara []map[string]interface{}
			err := ctx.BindJSON(&jsonpara) //从Body获取 json参数
			if err != nil {
				faygo.Info("Info: POST para is empty," + err.Error())
				//return ctx.JSONMsg(404, 404, err.Error())
				//如果参数为空则会触发EOF错误,不应该退出因为可能本来就没有参数，也就是jsonpara仍旧为空，需要创建该变量，后续sql中的参数处理需要
				if jsonpara == nil {
					jsonpara = make([]map[string]interface{}, 0)
				}
			}
			//.2.SQL定义参数验证并处理---OK，考虑将通过参数默认值处理的数据返回给客户端
			//将在服务端生成的uuid返回到客户端的变量
			var results []map[string]interface{}
			//未配置参数则直接忽略
			if len(se.Cmds[0].Parameters) > 0 {
				results = make([]map[string]interface{}, 0)
				for _, jp := range jsonpara {
					result, err := dealwithParameter(se.Cmds[0].Parameters, jp, ctx)
					if err != nil {
						faygo.Error(err.Error())
						return ctx.JSONMsg(400, 400, err.Error())
					}
					//
					if len(result) > 0 {
						results = append(results, result)
					}

				}
			}
			//.3.执行sql并返回结果
			err = m.bacthExecMap(se, jsonpara)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//如果存在服务端生成的uuid参数的则返回到客户端
			if (results != nil) && (len(results) > 0) {
				return ctx.JSONMsg(200, 200, results)
			} else {
				return ctx.JSONMsg(200, 200, "Bacth exec sql ok!")
			}

		case ST_BATCHMULTIEXEC: //批量複合語句
			//.1.获取 Ajax post json 参数
			var jsonpara map[string][]map[string]interface{}
			err := ctx.BindJSON(&jsonpara) //从Body获取 json参数
			if err != nil {
				faygo.Info("Info: POST para is empty," + err.Error())
				//return ctx.JSONMsg(404, 404, err.Error())
				//如果参数为空则会触发EOF错误,不应该退出因为可能本来就没有参数，也就是jsonpara仍旧为空，需要创建该变量，后续sql中的参数处理需要
				if jsonpara == nil {
					jsonpara = make(map[string][]map[string]interface{})
				}
			}
			//.2.SQL定义参数验证并处理---OK，考虑将通过参数默认值处理的数据返回给客户端
			//将在服务端生成的uuid返回到客户端的变量
			var results map[string][]map[string]interface{}
			results = make(map[string][]map[string]interface{})
			//循環每個sql定義
			for _, cmd := range se.Cmds {
				//未配置参数则直接忽略
				if len(cmd.Parameters) == 0 {
					continue
				}
				//循環其批量參數
				if sp, ok := jsonpara[cmd.Pin]; ok {
					result1 := make([]map[string]interface{}, 0)
					for _, p := range sp {
						result2, err := dealwithParameter(cmd.Parameters, p, ctx)
						if err != nil {
							faygo.Error(err.Error())
							return ctx.JSONMsg(400, 400, err.Error())
						}
						if len(result2) > 0 {
							result1 = append(result1, result2)
						}
						if len(result1) > 0 {
							results[cmd.Pin] = result1
						}
					}
				}
			}
			//.3.执行sql并返回结果
			err = m.bacthMultiExecMap(se, jsonpara)
			if err != nil {
				faygo.Error(err.Error())
				return ctx.JSONMsg(404, 404, err.Error())
			}
			//如果存在服务端生成的uuid参数的则返回到客户端
			if (results != nil) && (len(results) > 0) {
				return ctx.JSONMsg(200, 200, results)
			} else {
				return ctx.JSONMsg(200, 200, "Bacth Multi Exec sql ok!")
			}
		}
		return ctx.JSONMsg(404, 404, "Undefined sqltype!")
	}
}

//重新载入全部ModelSql配置文件
func DirectSQLReloadAll() faygo.HandlerFunc {
	return func(c *faygo.Context) error {
		ReloadAll()
		return c.JSONMsg(200, 200, "Info: Reload all modelsqls file ok!")
	}
}

//重新载入单个ModelSql配置文件
func DirectSQLReloadModel() faygo.HandlerFunc {
	return func(c *faygo.Context) error {
		//ctx.Path(), '/', 2) 去掉 /bom/reload/
		err := ReloadModel(trimBefore(c.Path(), '/', 3))
		if err != nil {
			return err
		}
		return c.JSONMsg(200, 200, "Info: Reload the modelsql file ok!")
	}
}

//发送JSON(P)响应
func sendJSON(ctx *faygo.Context, callback string, b []byte) error {
	//发送JSONP响应
	if len(callback) > 0 {
		callback = template.JSEscapeString(callback)
		callbackContent := bytes.NewBufferString(" if(window." + callback + ")" + callback)
		callbackContent.WriteString("(")
		callbackContent.Write(b)
		callbackContent.WriteString(");\r\n")
		return ctx.Bytes(200, faygo.MIMEApplicationJavaScriptCharsetUTF8, callbackContent.Bytes())
	}
	//正常有数据JSON响应
	return ctx.JSONBlob(200, b)
}

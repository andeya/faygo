/*
 功能：动态SQL执行函数供其他包调用单元
 日期：
 更新：2016.10.18
     增加 PagingSelectMapToMap func
*/
package directsql

import (
	"database/sql"
	"errors"
	"github.com/go-xorm/core"
	"github.com/henrylee2cn/faygo"
	"reflect"
)

var notFoundError = func(sqlid string) error {
	return errors.New("错误:未定义的sql:" + sqlid)
}
var notMatchError = func() error {
	return errors.New("错误:调用的语句的sqltype与该函数不匹配！")
}

//查询 根据modelId，sqlId ，mp:map[string]interface{}命名参数,返回*core.Rows
func SelectMapToRows(modelId, sqlId string, mp map[string]interface{}) (*core.Rows, error) {
	//获取Sqlentity,db
	se, db := findSqlAndDB(modelId, sqlId)
	if se == nil {
		return nil, notFoundError(modelId + "/" + sqlId)
	}
	//判断类型不是Select 就返回错误
	if se.Sqltype != ST_SELECT {
		return nil, notMatchError()
	}
	return db.QueryMap(se.Cmds[0].Sql, &mp)
}

//查询  根据modelId，sqlId ,SQL参数 map  返回 []map[string]interface{}
func SelectMapToMap(modelId, sqlId string, mp map[string]interface{}) ([]map[string]interface{}, error) {
	rows, err := SelectMapToRows(modelId, sqlId, mp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows2mapObjects(rows)
}

//查询 根据modelId，sqlId ，SQL参数是map, 返回 []struct
//目前使用比较繁琐：st －－是结构体的一个空实例，返回的是 改结构体的实例的slice，再使用返还结果时还的需要转换下类型。
func SelectMapToStruct(modelId, sqlId string, mp map[string]interface{}, st interface{}) (*[]interface{}, error) {

	s := reflect.ValueOf(st).Elem()
	leng := s.NumField()
	onerow := make([]interface{}, leng)
	for i := 0; i < leng; i++ {
		onerow[i] = s.Field(i).Addr().Interface()
	}
	result := make([]interface{}, 0)
	rows, err := SelectMapToRows(modelId, sqlId, mp)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(onerow...)
		if err != nil {
			return nil, err
			//panic(err)
		}
		result = append(result, s.Interface())
	}

	return &result, nil
}

//查询 根据modelId，sqlId ，SQL参数是map,dest 是待填充的返回结果 []*Struct ---未完成
func SelectMapToStructPro(modelId, sqlId string, mp map[string]interface{}, dest interface{}) error {
	return nil
}

//执行返回多個結果集的多個查询根据modelId，sqlId ，SQLmp:map[string]interface{}命名参数 返回结果 map[string]*Rows
func MultiSelectMapToRows(modelId, sqlId string, mp map[string]interface{}) (map[string]*core.Rows, error) {
	result := make(map[string]*core.Rows)
	//获取Sqlentity,db
	se, db := findSqlAndDB(modelId, sqlId)
	if se == nil {
		return nil, notFoundError(modelId + "/" + sqlId)
	}
	//判断类型不是MULTISELECT 就返回错误
	if se.Sqltype != ST_MULTISELECT {
		return nil, notMatchError()
	}
	//循環每個sql定義
	for i, cmd := range se.Cmds {
		faygo.Debug("MultiSelectMap :" + cmd.Sql)
		rows, err := db.QueryMap(cmd.Sql, &mp)
		if err != nil {
			return nil, err
		}
		if len(cmd.Rout) == 0 {
			result["data"+string(i)] = rows
		} else {
			result[cmd.Rout] = rows
		}
	}
	return result, nil
}

//执行分页查询SQL  mp 是MAP类型命名参数 返回结果 int,[]map[string][]interface{}
func PagingSelectMapToMap(modelId, sqlId string, mp map[string]interface{}) (*PagingSelectResult, error) {
	se, db := findSqlAndDB(modelId, sqlId)
	//获取总页数，約定該SQL放到第二條，並且只返回一條記錄一個字段
	trows, err := db.QueryMap(se.Cmds[0].Sql, &mp)
	if err != nil {
		return nil, err
	}
	defer trows.Close()
	for trows.Next() {
		var total = make([]int, 1)
		err := trows.ScanSlice(&total)
		if err != nil {
			return nil, err
		}
		if len(total) != 1 {
			return nil, errors.New("错误：获取总页数的SQL执行结果非唯一记录！")
		}
		//2.获取当前页數據，約定該SQL放到第二條
		rows, err := db.QueryMap(se.Cmds[1].Sql, &mp)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		result, err := rows2mapObjects(rows)
		if err != nil {
			return nil, err
		}
		return &PagingSelectResult{Total: total[0], Data: result}, nil //最終的結果
	}
	return nil, err
}

//多個查询 返回 map[string][]map[string]interface{}
func MultiSelectMapToMap(modelId, sqlId string, mp map[string]interface{}) (map[string][]map[string]interface{}, error) {
	multirows, err := MultiSelectMapToRows(modelId, sqlId, mp)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]map[string]interface{})
	for key, rows := range multirows {
		single, err := rows2mapObjects(rows)
		if err != nil {
			return nil, err
		}
		result[key] = single
	}
	return result, nil
}

//执行EXEC (UPDATE、DELETE、INSERT)，mp 是MAP类型命名参数 返回结果 sql.Result
func ExecMap(modelId, sqlId string, mp map[string]interface{}) (sql.Result, error) {
	//获取Sqlentity,db
	se, db := findSqlAndDB(modelId, sqlId)
	if se == nil {
		return nil, notFoundError(modelId + "/" + sqlId)
	}
	//判断类型不是EXEC(UPDATE、DELETE、INSERT)就返回错误
	if se.Sqltype != ST_EXEC {
		return nil, notMatchError()
	}
	return db.ExecMap(se.Cmds[0].Sql, &mp)
}

//执行EXEC (UPDATE、DELETE、INSERT)，SQL参数是struct  返回结果 sql.Result
func ExecStruct(modelId, sqlId string, st interface{}) (sql.Result, error) {
	//获取Sqlentity,db
	se, db := findSqlAndDB(modelId, sqlId)
	if se == nil {
		return nil, notFoundError(modelId + "/" + sqlId)
	}
	//判断类型不是EXEC 就返回错误
	if se.Sqltype != ST_EXEC {
		return nil, notMatchError()
	}
	return db.ExecStruct(se.Cmds[0].Sql, &st)
}

//批量执行 UPDATE、INSERT、DELETE、mp 是MAP类型命名参数
func BacthExecMap(modelId, sqlId string, sp []map[string]interface{}) error {
	//获取Sqlentity,db
	se, db := findSqlAndDB(modelId, sqlId)
	if se == nil {
		return notFoundError(modelId + "/" + sqlId)
	}
	//判断类型不是BATCHEXEC 就返回错误
	if se.Sqltype != ST_BATCHEXEC {
		return notMatchError()
	}
	return transact(db, func(tx *core.Tx) error {
		for _, p := range sp {
			faygo.Debug("BacthExecMap :" + se.Cmds[0].Sql)
			if _, err := tx.ExecMap(se.Cmds[0].Sql, &p); err != nil {
				return err
			}
		}
		return nil
	})
}

//批量执行 BacthComplex、mp 是MAP类型命名参数,事务中依次执行
func BacthMultiExecMap(modelId, sqlId string, mp map[string][]map[string]interface{}) error {
	//获取Sqlentity,db
	se, db := findSqlAndDB(modelId, sqlId)
	if se == nil {
		return notFoundError(modelId + "/" + sqlId)
	}
	//判断类型不是 ST_BATCHMULTIEXEC 就返回错误
	if se.Sqltype != ST_BATCHMULTIEXEC {
		return notMatchError()
	}
	return transact(db, func(tx *core.Tx) error {
		//循環每個sql定義
		for _, cmd := range se.Cmds {
			//循環其批量參數
			if sp, ok := mp[cmd.Pin]; ok {
				for _, p := range sp {
					faygo.Debug("BacthMultiExecMap :" + cmd.Sql)
					if _, err := tx.ExecMap(cmd.Sql, &p); err != nil {
						return err
					}
				}
			} else {
				return errors.New("错误：传入的参数与SQL节点定义的sql.pin名称不匹配！")
			}
		}
		return nil
	})
}

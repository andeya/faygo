/**
* func   : 动态SQL执行引擎
* author : 畅雨参数校验函数
* date   : 2016.06.13
* desc   : 关联关系全部通Id进行关联
* history :

		 -2018.12.18
		    -配置参数增加，批量SQL的事务处理：全部一个事务还是分开不同的事务中。
         - 2017.03.22
			-增加两种类型的sql，处理二进制对象保存到数据库和从数据库获取：
         	 ST_GETBLOB  //10 获取BLOB (binary large object)，二进制大对象从数据库
	         ST_SETBLOB  //11 保存BLOB (binary large object)，二进制大对象到数据库
         -2016.11.20 将执行sql的execMap修改该可以执行多个配置的cmd，采用相同的参数
*/
package directsql

import (
	"errors"
	"fmt"

	"xorm.io/core"
	"github.com/henrylee2cn/faygo"
)

//根据sqlid获取 *TSql
func (m *TModel) findSql(sqlid string) *TSql {
	if se, ok := m.Sqls[sqlid]; ok {
		//faygo.Debug("SqlId: " + sqlid)
		return se
	}
	return nil
}

//执行普通的单个查询SQL  mp 是MAP类型命名参数 map[string]interface{},返回结果 []map[string][]interface{}
func (m *TModel) selectMap(se *TSql, mp map[string]interface{}) ([]map[string]interface{}, error) {
	faygo.Debug("selectMap parameters :", mp)
	//执行sql
	rows, err := m.DB.QueryMap(se.Cmds[0].Sql, &mp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows2mapObjects(rows)
}

//分頁查詢的返回結果
type PagingSelectResult struct {
	Total int                      `json:"total"`
	Data  []map[string]interface{} `json:"data"`
}

//执行分页查询SQL  mp 是MAP类型命名参数 返回结果 int,[]map[string][]interface{}
func (m *TModel) pagingSelectMap(se *TSql, mp map[string]interface{}) (*PagingSelectResult, error) {
	faygo.Debug("pagingSelectMap parameters :", mp)
	//获取总页数，約定該SQL放到第二條，並且只返回一條記錄一個字段
	trows, err := m.DB.QueryMap(se.Cmds[0].Sql, &mp)
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
		rows, err := m.DB.QueryMap(se.Cmds[1].Sql, &mp)
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

//执行返回多個結果集的多個查询SQL， mp 是MAP类型命名参数 返回结果 map[string][]map[string][]string
func (m *TModel) multiSelectMap(se *TSql, mp map[string]interface{}) (map[string][]map[string]interface{}, error) {
	result := make(map[string][]map[string]interface{})
	faygo.Debug("MultiSelectMap parameters :", mp)
	//循環每個sql定義
	for i, cmd := range se.Cmds {
		faygo.Debug("MultiSelectMap :" + cmd.Sql)
		rows, err := m.DB.QueryMap(cmd.Sql, &mp)
		if err != nil {
			return nil, err
		}
		single, err := rows2mapObjects(rows)
		if err != nil {
			return nil, err
		}
		rows.Close()
		if len(cmd.Rout) == 0 {
			result["data"+string(i)] = single
		} else {
			result[cmd.Rout] = single
		}
	}
	return result, nil
}

//执行单个查询SQL返回JSON父子嵌套結果集 mp 是MAP类型命名参数 map[string]interface{},返回结果 []map[string][]interface{}
//根据 Idfield、Pidfield 构建嵌套的 map 结果集
func (m *TModel) nestedSelectMap(se *TSql, mp map[string]interface{}) ([]map[string]interface{}, error) {
	faygo.Debug("NestedSelectMap :" + se.Cmds[0].Sql)
	rows, err := m.DB.QueryMap(se.Cmds[0].Sql, &mp)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return rows2mapObjects(rows)
}

// execMap的返回结果定义
type Execresult struct {
	LastInsertId int64  `json:"lastinsertdd"`
	RowsAffected int64  `json:"rowsaffected"`
	Info         string `json:"info"`
}

//执行 UPDATE、DELETE、INSERT，mp 是 map[string]interface{}, 返回结果 execresult
/*func (m *TModel) execMap(se *TSql, mp map[string]interface{}) (*Execresult, error) {
	faygo.Debug("ExecMap :" + se.Cmds[0].Sql)
	faygo.Debug("map paras :", mp)
	Result, err := m.DB.ExecMap(se.Cmds[0].Sql, &mp)
	if err != nil {
		return nil, err
	}
	LIId, _ := Result.LastInsertId()
	RAffected, _ := Result.RowsAffected()
	return &Execresult{LastInsertId: LIId, RowsAffected: RAffected, Info: "Exec sql ok!"}, nil
}
*/
//说明，将执行sql的execMap修改该可以执行多个配置的cmd，采用相同的参数---2016.11.20
//执行 UPDATE、DELETE、INSERT，mp 是 map[string]interface{}，可以配置多个sql语句，使用相同的参数执行。
func (m *TModel) execMap(se *TSql, mp map[string]interface{}) error {
	faygo.Debug("ExecMap parameters :", mp)
	return transact(m.DB, func(tx *core.Tx) error {
		//循環每個sql定義
		for _, cmd := range se.Cmds {
			faygo.Debug("ExecMap sql:" + cmd.Sql)
			_, err := tx.ExecMap(cmd.Sql, &mp)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

//批量执行 UPDATE、INSERT、sp 是MAP类型命名参数
func (m *TModel) bacthExecMap(se *TSql, sp []map[string]interface{}) error {
	faygo.Debug("BacthExecMap parameters :", sp)
	// 如果执行sql的语句组使用每次循环使用单独的事务，则如下
	if se.Eachtran {
		for _, p := range sp {
			transact(m.DB, func(tx *core.Tx) error {
				////////////////////////////////
				//循環每個sql定義
				for _, cmd := range se.Cmds {
					faygo.Debug("BacthExecMap sql:" + cmd.Sql)
					_, err := tx.ExecMap(cmd.Sql, &p)
					if err != nil {
						return err
					}
				}
				return nil
				////////////////////////////////////
			})
		}
		//所有循环使用同一个事务
	} else {
		return transact(m.DB, func(tx *core.Tx) error {
			for _, p := range sp {
				////////////////////////////////
				//循環每個sql定義
				for _, cmd := range se.Cmds {
					faygo.Debug("BacthExecMap sql:" + cmd.Sql)
					_, err := tx.ExecMap(cmd.Sql, &p)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	}
	return nil
}

/*根据循环参数在一个事务中执行sql语句组（一个sql下多个cmd）
func (m *TModel) bacthExecMap(se *TSql, sp []map[string]interface{}) error {
	faygo.Debug("BacthExecMap parameters :", sp)
	return transact(m.DB, func(tx *core.Tx) error {
		for _, p := range sp {
			////////////////////////////////
			//循環每個sql定義
			for _, cmd := range se.Cmds {
				faygo.Debug("BacthExecMap sql:" + cmd.Sql)
				_, err := tx.ExecMap(cmd.Sql, &p)
				if err != nil {
					return err
				}
			}
			////////////////////////////////////
		}
		return nil
	})
}*/

/*原来的之只能根据参数循环执行一个exec语句
func (m *TModel) bacthExecMap(se *TSql, sp []map[string]interface{}) error {
	faygo.Debug("BacthExecMap parameters :", sp)
	return transact(m.DB, func(tx *core.Tx) error {
		for _, p := range sp {
			faygo.Debug("BacthExecMap :" + se.Cmds[0].Sql)
			if _, err := tx.ExecMap(se.Cmds[0].Sql, &p); err != nil {
				return err
			}
		}
		return nil
	})
}
*/
//批量执行 BacthMultiExec、mp 是map[string][]map[string]interface{}参数,事务中依次执行
func (m *TModel) bacthMultiExecMap(se *TSql, mp map[string][]map[string]interface{}) error {
	faygo.Debug("BacthMultiExecMap parameters : ", mp)
	// 如果执行sql的语句组使用每次循环使用单独的事务，规则如下：
	// 比如有4个SQL ，其中 1，2使用同一个参数组则1，2组合在一起每次循环使用一个事务，多次循环多个事务
	// ，3，4分别使用不同的参数组则各自也在不同的事务，规则同1，2
	if se.Eachtran {
		//根据参数循环
		for key, sp := range mp {
			for _, p := range sp {
				transact(m.DB, func(tx *core.Tx) error {
					//循環每個sql定義
					for _, cmd := range se.Cmds {
						//将使用相同参数的在一个事务执行
						if cmd.Pin == key {
							faygo.Debug("BacthMultiExecMap-EachTran :" + cmd.Sql)
							if _, err := tx.ExecMap(cmd.Sql, &p); err != nil {
								return err
							}
						}
					}
					return nil
				})
			}
		}
	} else { //原来的方式，在所有sql同一个事务中执行
		return transact(m.DB, func(tx *core.Tx) error {
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
	return nil
}

/* 原来代码 批量执行 BacthMultiExec、mp 是map[string][]map[string]interface{}参数,在同一个事务中依次执行
func (m *TModel) bacthMultiExecMap(se *TSql, mp map[string][]map[string]interface{}) error {
	return transact(m.DB, func(tx *core.Tx) error {
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
*/
/* / BinaryData的类型是blob/longblob
func addData(db *sql.DB, data []byte) error {
    _, err := db.Exec("INSERT INTO MyTable(BinaryData) VALUES(?)", data)
    return err
}*/
func (m *TModel) setBLOB(se *TSql, mp map[string]interface{}) error {
	faygo.Debug("setBLOB :" + se.Cmds[0].Sql)
	_, err := m.DB.ExecMap(se.Cmds[0].Sql, &mp)
	if err != nil {
		return err
	}
	return nil
}

//执行普通的查询SQL返回一个二进制字段的值   mp 是MAP类型命名参数 map[string]interface{},返回结果 []byte
func (m *TModel) getBLOB(se *TSql, mp map[string]interface{}) ([]byte, error) {
	faygo.Debug("getBLOB parameters :", mp)
	//执行sql
	rows, err := m.DB.QueryMap(se.Cmds[0].Sql, &mp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//字段slice
	fields, err := rows.Columns()
	if err != nil {
		faygo.Error(err)
		return nil, err
	}
	//只能返回一个字段，即要获取的二进制字段
	if len(fields) != 1 {
		return nil, errors.New("error: getBLOB sql only return one field")
	}
	if rows.Next() {
		blobvalue := make([]byte, 1024)
		err = rows.Scan(&blobvalue)
		if err != nil {
			faygo.Error(err)
			return nil, err
		}
		return blobvalue, nil
	}
	return nil, errors.New("error: getBLOB sql result is empty!")
}

//ransaction handler 封装在一个事务中执行多个SQL语句
func transact(db *core.DB, txFunc func(*core.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	return txFunc(tx)
}

/*
 功能：辅助函数
*/
package directsql

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-xorm/core"
)

//-------解析参数的函数------------
//将s从左边开始c出现的第n(>=1)次的位置之前的去掉 比如 'aa / bb / cc' -> 'bb / cc'
func trimBefore(s string, c byte, n int) string {
	r := 1
	for i := 0; i <= len(s)-1; i++ {
		if s[i] == c {
			if r == n {
				return s[i+1:]
			}
			r++
		}
	}
	return s
}

// 将s 根据从右边第一个出现的c进行分割成两个stirng,比如 'aa / bb / cc' -> 'aa / bb','cc'
func SplitRight(s string, c byte) (left, right string) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}

//先将s从左边开始c出现的第n(>=1)次的位置之前的去掉，在从右边第一个出现的c进行分割成两个stirng
func trimBeforeSplitRight(s string, c byte, n int) (left, right string) {
	r := 1

	for i := 0; i <= len(s)-1; i++ {
		if s[i] == c {
			if r == n {
				return SplitRight(s[i+1:], c)
			}
			r++
		}
	}
	return s, ""
}

//------------ rows转换相关函数 -------------------------
func reflect2value(rawValue *reflect.Value) (str string, err error) {
	aa := reflect.TypeOf((*rawValue).Interface())
	vv := reflect.ValueOf((*rawValue).Interface())
	switch aa.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = strconv.FormatInt(vv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = strconv.FormatUint(vv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		str = strconv.FormatFloat(vv.Float(), 'f', -1, 64)
	case reflect.String:
		str = vv.String()
	case reflect.Array, reflect.Slice:
		switch aa.Elem().Kind() {
		case reflect.Uint8:
			data := rawValue.Interface().([]byte)
			str = string(data)
		default:
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	// time type
	case reflect.Struct:
		if aa.ConvertibleTo(core.TimeType) {
			str = vv.Convert(core.TimeType).Interface().(time.Time).Format(time.RFC3339Nano)
		} else {
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	case reflect.Bool:
		str = strconv.FormatBool(vv.Bool())
	case reflect.Complex128, reflect.Complex64:
		str = fmt.Sprintf("%v", vv.Complex())
	/* TODO: unsupported types below
	   case reflect.Map:
	   case reflect.Ptr:
	   case reflect.Uintptr:
	   case reflect.UnsafePointer:
	   case reflect.Chan, reflect.Func, reflect.Interface:
	*/
	default:
		err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
	}
	return
}

func value2Bytes(rawValue *reflect.Value) (data []byte, err error) {
	var str string
	str, err = reflect2value(rawValue)
	if err != nil {
		return
	}
	data = []byte(str)
	return
}

func value2String(rawValue *reflect.Value) (data string, err error) {
	data, err = reflect2value(rawValue)
	if err != nil {
		return
	}
	return
}

func rows2Strings(rows *core.Rows) (resultsSlice []map[string]string, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2mapStr(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func rows2maps(rows *core.Rows) (resultsSlice []map[string][]byte, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2map(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func row2map(rows *core.Rows, fields []string) (resultsMap map[string][]byte, err error) {
	result := make(map[string][]byte)
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			//fmt.Println("ignore ...", key, rawValue)
			continue
		}

		if data, err := value2Bytes(&rawValue); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}
	}
	return result, nil
}

func row2mapStr(rows *core.Rows, fields []string) (resultsMap map[string]string, err error) {
	result := make(map[string]string)
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			//fmt.Println("ignore ...", key, rawValue)
			continue
		}

		if data, err := value2String(&rawValue); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}
	}
	return result, nil
}

//转换 interface{} 到 JSON
func JSONString(v interface{}, Indent bool) (string, error) {
	var result []byte
	var err error
	if Indent {
		result, err = json.MarshalIndent(v, "", "  ")
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		return "", err
	}

	if string(result) == "null" {
		return "", nil
	}
	return string(result), nil
}

//转换接口为字节数组格式的json
func intface2json(i interface{}) ([]byte, error) {
	var (
		b   []byte
		err error
	)
	b, err = json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//-----------------------------------------------------------------------------------------
func reflect2object(rawValue *reflect.Value) (value interface{}, err error) {
	aa := reflect.TypeOf((*rawValue).Interface())
	vv := reflect.ValueOf((*rawValue).Interface())
	switch aa.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = vv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = vv.Uint()
	case reflect.Float32, reflect.Float64:
		value = vv.Float()
	case reflect.String:
		value = vv.String()
	case reflect.Array, reflect.Slice:
		switch aa.Elem().Kind() {
		case reflect.Uint8:
			data := rawValue.Interface().([]byte)
			value = string(data)
		default:
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	//时间类型
	case reflect.Struct:
		if aa.ConvertibleTo(core.TimeType) {
			value = vv.Convert(core.TimeType).Interface().(time.Time)
		} else {
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	case reflect.Bool:
		value = vv.Bool()
	case reflect.Complex128, reflect.Complex64:
		value = vv.Complex()
	/* TODO: unsupported types below
	   case reflect.Map:
	   case reflect.Ptr:
	   case reflect.Uintptr:
	   case reflect.UnsafePointer:
	   case reflect.Chan, reflect.Func, reflect.Interface:
	*/
	default:
		err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
	}
	return
}

func value2Object(rawValue *reflect.Value) (data interface{}, err error) {
	data, err = reflect2object(rawValue)
	if err != nil {
		return
	}
	return
}

func rows2mapObjects(rows *core.Rows) (resultsSlice []map[string]interface{}, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := rows2mapObject(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func rows2mapObject(rows *core.Rows, fields []string) (resultsMap map[string]interface{}, err error) {
	result := make(map[string]interface{})
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			continue
		}

		if data, err := value2Object(&rawValue); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}

	}
	return result, nil
}

//rows 转换为嵌套的map,父子关系根据 Idfield,Pidfield
func rows2nestedMapObjects(rows *core.Rows, idfield, pidfield string) (resultsSlice []map[string]interface{}, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := rows2mapObject(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

//Struct2Map
func Struct2Map(st interface{}) map[string]interface{} {
	t := reflect.TypeOf(st)
	v := reflect.ValueOf(st)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// map值转化到struct中
func Map2Struct(mp map[string]interface{}, st interface{}) error {
	for k, v := range mp {
		err := setStructField(st, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func setStructField(st interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(st).Elem()         //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in st", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值
	//var err error
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
		//val, err = value2Reflect(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		//if err != nil {
		//	return err
		//}
	}

	structFieldValue.Set(val)
	return nil
}

/*
//类型转换:将值类型转换为反射类型－map转struct使用
func value2Reflect(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "bool" {
		i, err := strconv.ParseBool(value, 64)
		return reflect.ValueOf(i), err
	}
	//else if .......增加其他一些类型的转换
	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}*/

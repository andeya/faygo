/**
* func   : sql参数值的服务端来源定义
* author : 畅雨
* date   : 2016.12.17
* desc   : 可以将常量，变量，函数的值提供给sql参数的服务端默认值，通过reg函数注册到Context(MAP)中,在SQL的参数中通过名字进行调用得到值作为sql参数的值
* history :

 */
package directsql

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/henrylee2cn/faygo"
)

var reIdentifiers = regexp.MustCompile("^[a-zA-Z0-9_]+$")

// A Context type provides constants, variables or functions to a sql parameter's default value.
var sqlcontext map[string]interface{}

func init() {
	sqlcontext = make(map[string]interface{})
	// test func
	RegAny("nowtime", gettime)
}

//
func gettime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//注册新的变量或函数到map
func RegAny(name string, fn interface{}) (err error) {
	//判断名称是否合法
	if !reIdentifiers.MatchString(name) {
		return errors.New(fmt.Sprintf("SQLContext-key '%s' (value: '%+v') is not a valid identifier.", name, fn))
	}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(name + " is not callable.")
		}
	}()
	sqlcontext[name] = fn
	return
}

// call function by name and paramarers
func contextcall(name string, params ...interface{}) (result reflect.Value, err error) {
	if _, ok := sqlcontext[name]; !ok {
		err = errors.New(name + " does not exist.")
		return
	}
	fv := reflect.ValueOf(sqlcontext[name])
	faygo.Debug("Context type : ", name)
	faygo.Debug("Context Kind : ", fv.Kind())
	// is function？
	if fv.Kind() == reflect.Func {
		t := fv.Type()
		faygo.Debug("Context Func: ", t)
		//Check input arguments
		if len(params) != t.NumIn() {
			err = errors.New("parameters of function not adapted")
			return
		}
		// Check output arguments
		if t.NumOut() != 1 {
			err = fmt.Errorf("'%s' must have exactly 1 output argument", fv.String())
			return
		}
		in := make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		//return result
		result = fv.Call(in)[0]
		faygo.Debug("Context func value : ", result)
	} else {
		// value
		result = reflect.ValueOf(fv.Interface())
	}

	return
}

/**
* desc   : 参数校验函数
* author : 畅雨
* date   : 2016.12.13
* desc   :
* history :

 */
package directsql

import (
	"errors"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/uuid"
)

// 类型判断正则表达式定义
const (
	Int   string = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Float string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	Email string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

//字符串是否合法Email地址
func IsEmail(str string) bool {
	return regexp.MustCompile(Email).MatchString(str)
}

//字符串是否整数，空也是合法的.
func IsInt(str interface{}) bool {
	if _, ok := str.(float64); ok {
		return regexp.MustCompile(Int).MatchString(strconv.FormatFloat(str.(float64), 'f', -1, 64))
	}
	return false
}

//字符串是否浮点数
func IsFloat(str interface{}) bool {
	_, ok := str.(float64)
	return ok
}

//字符串是否有效的长度
func IsVaildLength(str string, min, max int) bool {
	//min,max都未定义或都为0则不校验长度
	if min <= 0 && max <= 0 {
		return true
	}
	strLength := utf8.RuneCountInString(str)
	//只验证最小长度，最大长度不验证
	if min > 0 && max <= 0 {
		return strLength >= min
	}
	return strLength >= min && strLength <= max
}

//给定的数值是否在范围内
func IsVaildValue(value, min, max float64) bool {
	//min,max都未定义或都为0则不校验大小
	if min <= 0 && max <= 0 {
		return true
	}
	//只验证最小值，最大值不验证
	if min > 0 && max <= 0 {
		return value >= min
	}
	if min > max {
		min, max = max, min
	}
	return value >= min && value <= max
}

//给定的字符串是否合法的日期时间
func IsVaildDatetime(str string) bool {
	_, err := time.Parse("2006-01-02 15:04", str)
	return err == nil
}

//给定的字符串是否合法的日期(YYYY-MM-DD)
func IsVaildDate(str string) bool {
	_, err := time.Parse("2006-01-02", str)
	return err == nil
}

//检查是否必须的
func CheckRequired(str string) bool {
	return len(str) > 0
}

//检查并处理参数
//paras:sql的cmd中的参数定义slice；mp:客户端提交的参数map；ctx *lessgo.Context当前执行该功能的上下文
//根据待默认值的参数是否需要返回构造返回到客户端的值
func dealwithParameter(paras []*TSqlParameter, mp map[string]interface{}, ctx *faygo.Context) (map[string]interface{}, error) {
	//没有参数处理定义返回
	if len(paras) == 0 {
		//faygo.Debug("Check sql parameters - nil")
		return nil, nil
	}
	//将在服务端生成的默认值需要返回的放入到该结果map中。
	var result map[string]interface{}
	result = make(map[string]interface{})
	//循环处理参数
	for _, para := range paras {
		//默认值处理，存在就不处理使用存在的值，不存在就增加并返回给客户端
		_, exists := mp[para.Name]
		//不是从客户的传入的并且有默认值设置
		if (!exists) && (para.Default != DT_UNDEFINED) {
			//根据默认值类型分别处理
			switch para.Default {
			case DT_UUID: //uuid
				mp[para.Name] = uuid.New().String()
			case DT_NOWDATE: // now date
				mp[para.Name] = time.Now().Format("2006-01-02")
			case DT_NOWDATETIME: //now(date +time)
				mp[para.Name] = time.Now().Format("2006-01-02 15:04:05")
			case DT_NOW_UNIX: //当前日期时间unix值 int64 now datetime
				mp[para.Name] = time.Now().Unix()
			case DT_CUSTOM: //通过RegAny注册的变量或函数
				value, err := contextcall(para.Defaultstr, ctx)
				//faygo.Debug("SQL Default parameter value:", value)
				if err == nil {
					mp[para.Name] = value.Interface()
				} else {
					faygo.Error("Error: sql default parameter error,", err)
				}
			}
			//如果需要返回
			if para.Return {
				result[para.Name] = mp[para.Name]
			}
			//如果不是从客户的传入的并且有默认值设置则后边的验证规则不执行了，如果是从客户的传入的则需要进行进行后边的默认值校验
			continue
		}

		//根据参数名称从提交的参数中获取值循环验证类型、长度、是否为空等信息
		if v, ok := mp[para.Name]; ok {
			//是否必须的
			if _, ok := v.(string); ok && (para.Required) && (len(v.(string)) == 0) {
				return nil, errors.New("错误：参数[" + para.Name + "]不能为空！")
			}
			//faygo.Debug("Check sql parameters - get value")
			//参数类型处理
			switch para.Paratype {
			case PT_STRING:
				//验证长度,是否必须的
				if IsVaildLength(v.(string), para.Minlen, para.Maxlen) {
					continue
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]长度不符合定义范围！")
				}

			case PT_INT:
				//faygo.Debug("Check sql int parameter -", v)
				//验证是否整数
				if IsInt(v) {
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]不是有效的整数！")
				}
				//验证数值范围
				if IsVaildValue(v.(float64), para.MinValue, para.MaxValue) {
					continue
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]的值不符合定义范围！")
				}

			case PT_FLOAT:
				//验证是否浮点数
				if IsFloat(v) {
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]不是有效的浮点数！")
				}
				//验证数值范围
				if IsVaildValue(v.(float64), para.MinValue, para.MaxValue) {
					continue
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]的值不符合定义范围！")
				}

			case PT_DATE:
				//验证日期格式
				if IsVaildDate(v.(string)) {
					continue
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]不是有效的日期！")
				}
			case PT_DATETIME:
				//验证日期时间格式
				if IsVaildDatetime(v.(string)) {
					continue
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]不是有效的日期时间！")
				}
			case PT_EMAIL:
				//验证email格式
				if IsEmail(v.(string)) {
					continue
				} else {
					return nil, errors.New("错误：参数[" + para.Name + "]不是有效的电子信箱！")
				}
			}
			//faygo.Debug("Check sql parameters - " + para.Name + ": " + v.(string))
		} else {
			//sql的cmd参数中存在该参数定义但传入的post参数不存在则返回错误
			return nil, errors.New("错误：参数[" + para.Name + "]未定义！")
		}
	}
	return result, nil
}

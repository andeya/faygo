/**
* desc   : 管理动态SQL,根据配置文件目录从配置文件加载到内存中
* author : 畅雨
* date   : 2016.12.13
* desc   :
* history:
 */
package directsql

import (
	"encoding/xml"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/go-xorm/core"
	"github.com/henrylee2cn/faygo"
	faygoxorm "github.com/henrylee2cn/faygo/ext/db/xorm"
	confpkg "github.com/henrylee2cn/faygo/ini"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//var modelsqls map[string]*TModel

//配置文件配置参数
const MSCONFIGFILE = "./config/directsql.ini"

// 全部业务SQL路由表,不根据目录分层次，直接放在map sqlmodels中，key=带路径不带扩展名的文件名
type TModels struct {
	roots     map[string]string  //需要载入的根目录(可以多个)-短名=实际路径名
	modelsqls map[string]*TModel //全部定义模型对象
	extension string             //模型定义文件的扩展名(默认为.msql)
	lazyload  bool               //true ／ false 配置项  true 则一开始不全部加载，只根据第一次请求才加载然后缓存，false＝一开始就根据配置的roots目录全部加载
	loadLock  sync.RWMutex       //读写锁
	watcher   *fsnotify.Watcher  //监控文件变化的wather
	cached    bool               //是否启用查询数据缓存功能，启用则配置sql中配置属性cached=true 才有效，否则一律不缓存
	cachetime int                //默认缓存的时间，如果使用缓存并且未配置缓存时间则使用该默认时间，单位为分钟，-1为一直有效，-2为一月，-3为一周 -4为一天，单位为分钟
}

//全局所有业务模型对象
var models = &TModels{
	modelsqls: make(map[string]*TModel)}

//sqlmodel 一个配置文件的SQLModel对应的结构
type TModel struct {
	Id   string           //root起用映射、不带扩展名的文件名
	DB   *core.DB         //本模块的db引擎 *xorm.Engine.DB()
	Sqls map[string]*TSql //sqlentity key=sqlentity.id
}

//临时转换用，因为 XML 不支持解析到 map，所以先读入到[]然后再根据[]创建map
type tempTModel struct {
	XMLName  xml.Name `xml:"model"`
	Id       string   `xml:"id,attr"` //不带扩展名的文件名
	Database string   `xml:"database,attr"`
	Sqls     []*TSql  `xml:"sql"`
}

//sql <Select/>等节点对应的结构
type TSql struct {
	XMLName    xml.Name `xml:"sql"`
	Id         string   `xml:"id,attr"` //sqlid
	Sqltypestr string   `xml:"type,attr"`
	Sqltype    TSqltype `xml:"-"`              //SQL类型
	Idfield    string   `xml:"idfield,attr"`   //SQlType为6=嵌套jsoin树时的ID字段
	Pidfield   string   `xml:"pidfield,attr"`  //SQlType为6=嵌套jsoin树时的ParentID字段
	Cmds       []*TCmd  `xml:"cmd"`            // sqlcmd(sqltype为分页查询时的计数SQL放第一个，结果SQL放第二个)
	Cached     bool     `xml:"cached,attr"`    //是否启用查询数据缓存功能
	Cachetime  int      `xml:"cachetime,attr"` //默认缓存的时间，单位为分钟，-1为一直有效，-2为一月，-3为一周 -4为一天，单位为分钟
}

//TCmd  <Select/>等节点的下级节点<sql />对应结构
type TCmd struct {
	XMLName    xml.Name         `xml:"cmd"`
	Pin        string           `xml:"in,attr"`   //输入参数标示
	Rout       string           `xml:"out,attr"`  //输出结果标示
	Sql        string           `xml:",chardata"` //SQL
	Parameters []*TSqlParameter `xml:"parameters>parameter"`
}

//TSql 类型
type TSqltype int

const (
	ST_SELECT         TSqltype = iota //0 普通查询 ---OK!
	ST_PAGINGSELECT                   //1 分页查询 ---OK!
	ST_NESTEDSELECT                   //2 嵌套jsoin树---------未实现
	ST_MULTISELECT                    //3 多结果集查询---OK!
	ST_EXEC                           //4 执行SQL，可以一个事务内批量执行多个cmd
	ST_BATCHEXEC                      //5 根据传入参数在一个事务内多次执行SQL
	ST_BATCHMULTIEXEC                 //6 批量执行复合SQL(多数据集批量插入、更新、删除)---OK!
	ST_IMPORT                         //7 导入数据的SQL：通过xlsx导入数据配置的SQL
	ST_EXPORT                         //8 导出数据的SQL：导出excel格式文件数据
	ST_REPORT                         //9 报表用的SQL：通过xlsx模板创建报表的SQL
)

//TSqlParameter 参数校验定义
type TSqlParameter struct {
	Name        string       `xml:"name,attr"`     //参数名称必须与cmd中的对应
	Paratypestr string       `xml:"type,attr"`     //string/number/email/date/datetime/time  -不定义则不需要验证
	Paratype    TParaType    `-`                   //数值类型
	Required    bool         `xml:"required,attr"` //0=不是必须的   1=必须的不能为空
	Minlen      int          `xml:"minlen,attr"`   //最小长度
	Maxlen      int          `xml:"maxlen,attr"`   //最大长度
	MinValue    float64      `xml:"minvalue,attr"` //最小值
	MaxValue    float64      `xml:"maxvalue,attr"` //最大值
	Defaultstr  string       `xml:"default,attr"`  // 默认值 undefined/uuid/userid/usercode/username/rootgroupid/rootgroupname/groupid/groupname/nowdate/nowtime
	Default     TDefaultType `-`                   //数值类型
	Return      bool         `xml:"return,attr"`   //服务端生成的默认值是否返回到客户端： 0(false)=默认，不返回   1(true)=返回到客户端
}

//参数类型：string/number/email/date/datetime/time
type TParaType int

const (
	PT_STRING   TParaType = iota //0=字符串,默认就是按照字符串处理
	PT_INT                       //1=整数数值
	PT_FLOAT                     //2=浮点数
	PT_DATE                      //3=日期
	PT_DATETIME                  //4=日期时间
	PT_EMAIL                     //5=电子邮件
	PT_BLOB                      //6=二进制
)

//----------------------------------------------------------------------------------------------------
// 默认值类型：uuid/nowdate/now/nowunix
type TDefaultType int

const (
	DT_UNDEFINED   TDefaultType = iota //0=未定义，不处理默认值
	DT_UUID                            //1=uuid
	DT_NOWDATE                         //2=当前日期 now date
	DT_NOWDATETIME                     //3=当前日期时间 now datetime
	DT_NOW_UNIX                        //4=当前时间的unix值 int64 now date
	DT_CUSTOM                          //5=自定义，采用注册自定义变量实现

)

//---------------------------------------------------------------------------------------------------
//读入全部模型sql定义
func init() {
	models.loadTModels()

}

//读入全部模型
func (ms *TModels) loadTModels() {
	ms.loadLock.Lock()
	defer ms.loadLock.Unlock()
	//打开directsql的配置文件
	cfg, err := confpkg.Load(MSCONFIGFILE)
	if err != nil {
		faygo.Error(err.Error())
		return
	}
	//是否缓存与缓存时间
	ms.cached = cfg.Section("").Key("cached").MustBool(false)
	ms.cachetime = cfg.Section("").Key("cachetime").MustInt(30)

	//读取ModelSQL文件的根目录
	roots, err := cfg.GetSection("roots")
	if err != nil {
		faygo.Error(err.Error())
		return
	}

	ms.roots = make(map[string]string)

	for _, v := range roots.Keys() {
		if len(v.String()) > 3 {
			ms.roots[v.Name()] = v.String()
		}
	}

	//读取扩展名，读取不到就用默认的.msql
	ext := cfg.Section("").Key("ext").MustString(".msql")
	ms.extension = ext

	//根据路径遍历加载
	for _, value := range ms.roots {
		faygo.Debug(value)
		err = filepath.Walk(value, ms.walkFunc)
		if err != nil {
			faygo.Error(err.Error())
		}

	}
	//是否监控文件变化
	watch := cfg.Section("").Key("watch").MustBool(false)
	if watch {
		err := ms.StartWatcher()
		if err != nil {
			faygo.Error(err.Error())
		}
	}

}

//将带路径文件名处理成 TModel的 id 示例： bizmodel\demo.msql  --> biz/demo
func (ms *TModels) filenameToModelId(path string) string {
	key := strings.Replace(path, "\\", "/", -1)
	key = strings.TrimSuffix(key, ms.extension) //去掉扩展名
	for root, value := range ms.roots {
		if strings.HasPrefix(key, value) {
			key = strings.Replace(key, value, root, 1) //处理前缀,将定义的根路径替换为名称
			break
		}
	}
	return key
}

//遍历子目录文件处理函数
func (ms *TModels) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if strings.HasSuffix(path, ms.extension) {
		m, err := ms.parseTModel(path)
		if err != nil {
			faygo.Error("Model file: " + path + " --- " + err.Error())
			return nil //单个文件解析出错继续加载其他的文件
		}
		//将本文件对应的TModel放入到TModels
		ms.modelsqls[ms.filenameToModelId(path)] = m
		faygo.Debug("Model file: " + path + " ------> " + ms.filenameToModelId(path) + "  loaded. ")
	}
	return nil
}

//解析单个ModelSQL定义文件
func (ms *TModels) parseTModel(msqlfile string) (*TModel, error) {
	//读取文件
	content, err := ioutil.ReadFile(msqlfile)
	if err != nil {
		return nil, err
	}
	tempresult := tempTModel{}
	err = xml.Unmarshal(content, &tempresult)
	if err != nil {
		return nil, err
	}
	//设置数据库
	dbe, ok := faygoxorm.DB(tempresult.Database)
	if ok == false {
		dbe = faygoxorm.MustDB()
		//faygo.Log.Debug("database:", tempresult.Database)
	}
	//定义一个 TModel将 tempTModel 转换为 TModel
	result := &TModel{Id: tempresult.Id, DB: dbe.DB(), Sqls: make(map[string]*TSql)}
	//处理一遍：设置数据库访问引擎，设置TSql的类型
	for _, se := range tempresult.Sqls {
		//处理SQL类型与查询类语句缓存的配置参数
		switch se.Sqltypestr {
		case "select":
			se.Sqltype = ST_SELECT
			//缓存设置处理
			se.Cached = se.Cached && ms.cached
			if se.Cached && se.Cachetime == 0 {
				se.Cachetime = ms.cachetime
			}
		case "pagingselect":
			se.Sqltype = ST_PAGINGSELECT
			//缓存设置处理
			se.Cached = se.Cached && ms.cached
			if se.Cached && se.Cachetime == 0 {
				se.Cachetime = ms.cachetime
			}
		case "nestedselect":
			se.Sqltype = ST_NESTEDSELECT
			//缓存设置处理
			se.Cached = se.Cached && ms.cached
			if se.Cached && se.Cachetime == 0 {
				se.Cachetime = ms.cachetime
			}
		case "multiselect":
			se.Sqltype = ST_MULTISELECT
			//缓存设置处理
			if !ms.cached {
				se.Cached = false
			}
			//se.cached = se.cached && ms.cached
			if se.Cached && se.Cachetime == 0 {
				se.Cachetime = ms.cachetime
			}
		case "exec", "insert", "update", "delete":
			se.Sqltype = ST_EXEC
		case "batchexec", "batchinsert", "batchupdate", "batchdelete":
			se.Sqltype = ST_BATCHEXEC
		case "batchmultiexec", "batchcomplex":
			se.Sqltype = ST_BATCHMULTIEXEC
		case "import":
			se.Sqltype = ST_IMPORT
		case "export":
			se.Sqltype = ST_EXPORT
		case "report":
			se.Sqltype = ST_REPORT
		default:
			faygo.Error(errors.New("错误：配置文件[ " + msqlfile + " ]中存在无效的sql节点类型[ " + se.Sqltypestr + " ]!"))
		}
		result.Sqls[se.Id] = se
		//faygo.Debug(se)
		//sql下的每个cmd循环处理
		for _, cmd := range se.Cmds {
			//每个cmd下的参数循环处理参数类型与默认值类型
			for _, para := range cmd.Parameters {
				//参数类型
				switch para.Paratypestr { //string/int/float/email/date/datetime/blob
				case "string":
					para.Paratype = PT_STRING //0=字符串,默认就是按照字符串处理
				case "int":
					para.Paratype = PT_INT //1=整数数值
				case "float":
					para.Paratype = PT_FLOAT //2=浮点数
				case "date":
					para.Paratype = PT_DATE //3=日期
				case "datetime":
					para.Paratype = PT_DATETIME //4=日期时间
				case "email":
					para.Paratype = PT_EMAIL //5=电子邮件
				case "blob":
					para.Paratype = PT_BLOB
				}
				//默认值类型
				switch para.Defaultstr {
				case "uuid":
					para.Default = DT_UUID //1=uuid
				case "nowdate":
					para.Default = DT_NOWDATE //12=当前日期 now date
				case "now":
					para.Default = DT_NOWDATETIME //当前时间
				case "nowunix":
					para.Default = DT_NOW_UNIX //=当前日期时间unix值 int64 now datetime
				default:
					if len(strings.TrimSpace(para.Paratypestr)) > 0 {
						para.Default = DT_CUSTOM
					} else {
						para.Default = DT_UNDEFINED
					}
				}
			}
		}
	}
	return result, nil
}

//获取sqlentity SQL的执行实体
func (ms *TModels) findsql(modelid string, sqlid string) *TSql {
	if sm, ok := ms.modelsqls[modelid]; ok {
		if se, ok := sm.Sqls[sqlid]; ok {
			return se
		}
	}
	return nil
}

//获取sqlentity SQL的执行实体与DB执行引擎
func (ms *TModels) findsqlanddb(modelid string, sqlid string) (*TSql, *core.DB) {
	if sm, ok := ms.modelsqls[modelid]; ok {
		if se, ok := sm.Sqls[sqlid]; ok {
			return se, sm.DB
		}
	}
	return nil, nil
}

//获取sqlentity 的类型
func (ms *TModels) getSqlType(modelid string, sqlid string) TSqltype {
	if sm, ok := ms.modelsqls[modelid]; ok {
		if se, ok := sm.Sqls[sqlid]; ok {
			return se.Sqltype
		}
	}
	return -1
}

// 根据路径加文件名(不带文件扩展名)获取其TModel
func (ms *TModels) findmodel(modelid string) *TModel {
	if sm, ok := ms.modelsqls[modelid]; ok {
		return sm
	}
	return nil
}

//文件内容改变重新载入(新增、修改的都触发)
func (ms *TModels) refreshModelFile(msqlfile string) error {
	ms.loadLock.Lock()
	defer ms.loadLock.Unlock()
	//重新解析
	m, err := ms.parseTModel(msqlfile)
	if err != nil {
		faygo.Error(err.Error())
		return err //单个文件解析出错继续加载其他的文件
	}
	//将本文件对应的TModel放入到TModels
	ms.modelsqls[ms.filenameToModelId(msqlfile)] = m
	return nil
}

//文件已经被移除，从内存中删除
func (ms *TModels) removeModelFile(msqlfile string) error {
	ms.loadLock.Lock()
	defer ms.loadLock.Unlock()
	delete(ms.modelsqls, ms.filenameToModelId(msqlfile))
	return nil
}

//文件改名---暂无实现
func (ms *TModels) renameModelFile(msqlfile, newfilename string) error {
	//err := ms.removeModelFile(msqlfile)
	//err = ms.refreshModelFile(newfilename)
	return nil
}

//单元访问文件--------------------------------------------------------------
//获取sqlentity SQL的执行实体与数据库引擎
func findSqlAndDB(modelid string, sqlid string) (*TSql, *core.DB) {
	return models.findsqlanddb(modelid, sqlid)
}

//
func GetSqlType(modelid string, sqlid string) TSqltype {
	return models.getSqlType(modelid, sqlid)
}

//获取sqlentity SQL的执行实体
func findSql(modelid string, sqlid string) *TSql {
	//faygo.Debug("Model Path: " + modelid + " ,SqlId: " + sqlid)
	return models.findsql(modelid, sqlid)
}

//根据TModel文件路径获取 TModel
func findModel(modelid string) *TModel {
	//faygo.Debug("Model Path: " + modelid)
	return models.findmodel(modelid)
}

//重置配置文件全部重新载入,API：/bom/reload  handle调用
func ReloadAll() {
	models = &TModels{
		modelsqls: make(map[string]*TModel)}
	models.loadTModels()
}

//重新载入单个模型文件---未测试！！！
func ReloadModel(msqlfile string) error {
	//已经去掉 "/bom/reload/",需要加上扩展名
	return models.refreshModelFile(msqlfile + models.extension)
}

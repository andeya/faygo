/**
* func   : 查询sql结果缓存处理单元
* author : 畅雨
* date   : 2016.11.25
* desc   : 缓存查询的结果
* history :
           - 2106.11.30 -优化缓存的存储
*/
package directsql

import (
	"github.com/henrylee2cn/faygo"
	"sync"
	"time"
)

//缓存结果对象
type memo struct {
	Timeout time.Time
	Result  []byte //interface{}
	Suffix  string //附加的标识，对于有参数的sql使用该后缀标识参数差异，只缓存第一次查询默认参数的情况，其他的不缓存，缓存的刷新也依赖该
}

//缓存池对象
type MemoPool struct {
	pool  map[string]*memo
	mutex *sync.RWMutex
}

//全局缓存池
var mp *MemoPool

func init() {
	mp = &MemoPool{pool: map[string]*memo{}, mutex: new(sync.RWMutex)}
}

//根据key以及suffix后缀获取缓存的结果 has 表示存在有效的result，result为结果
func GetCache(key string, suffix string) (ok bool, result []byte) {
	mp.mutex.RLock()
	memoized := mp.pool[key]
	mp.mutex.RUnlock()
	//lessgo.Log.Debug("Get Cache:[" + key + " - " + suffix + "]")
	// 有缓存,并且后缀标识相同，且未到期则返回 true 以及缓存值
	if (memoized != nil) && (memoized.Suffix == suffix) && memoized.Timeout.After(time.Now()) {
		//lessgo.Log.Debug("Get Cache Result:["+key+" - "+suffix+"] ", memoized.Result)
		faygo.Debug("Get Cache time:["+key+" - "+suffix+"]", time.Now())
		faygo.Debug("Get Cache timeout:["+key+" - "+suffix+"]", memoized.Timeout)
		return true, memoized.Result
	}
	return false, nil
}

//将key以及suffix后缀的值放入到缓存中，如果存在则替换，并记录失效日期
func SetCache(key string, suffix string, value []byte, timeout int) {
	if value != nil {
		mp.mutex.RLock()
		memoized := mp.pool[key]
		mp.mutex.RUnlock()
		//缓存key存在值但后缀不同则退出
		if (memoized != nil) && (memoized.Suffix != suffix) {
			return
		}
		//faygo.Debug("Set Cache:[" + key + " - " + suffix + "]")
		//缓存不存在或虽然存在但suffix后缀相同则修改之
		var duration time.Duration
		//-1为一直有效，-2为一月，-3为一周 -4为一天，单位为分钟
		switch timeout {
		case -1:
			duration = 365 * time.Duration(24) * time.Hour //一年
		case -2:
			duration = 30 * time.Duration(24) * time.Hour //一月
		case -3:
			duration = 7 * time.Duration(24) * time.Hour //一周
		case -4:
			duration = time.Duration(24) * time.Hour //一天
		default:
			duration = time.Duration(timeout) * time.Minute //分钟
		}
		mp.mutex.Lock()
		mp.pool[key] = &memo{
			Suffix:  suffix,
			Timeout: time.Now().Add(duration),
			Result:  value,
		}
		mp.mutex.Unlock()
		//faygo.Debug("Set Cache Result:["+key+" - "+suffix+"]", mp.pool[key])
	}
}

//清除key的缓存
func RemoveCache(key string) {
	mp.mutex.Lock()
	delete(mp.pool, key)
	mp.mutex.Unlock()
}

//清除全部缓存
func ClearCache() {
	mp.mutex.Lock()
	for key, _ := range mp.pool {
		delete(mp.pool, key)
	}
	mp.mutex.Unlock()
}

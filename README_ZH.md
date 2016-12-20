# Thinkgo    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/thinkgo)    ![Thinkgo goreportcard](https://goreportcard.com/badge/github.com/henrylee2cn/thinkgo)

![Thinkgo Favicon](https://github.com/henrylee2cn/thinkgo/raw/master/doc/thinkgo_96x96.png)

## 概述
Thinkgo以全新的架构实现，它面向Handler接口开发，支持智能参数映射与校验、支持自动化API文档的Go语言web框架。

官方QQ群：Go-Web 编程 42730308    [![Go-Web 编程群](http://pub.idqqimg.com/wpa/images/group.png)](http://jq.qq.com/?_wv=1027&k=fzi4p1)

![thinkgo server](https://github.com/henrylee2cn/thinkgo/raw/master/doc/server.png)

![thinkgo apidoc](https://github.com/henrylee2cn/thinkgo/raw/master/doc/apidoc.png)

![thinkgo index](https://github.com/henrylee2cn/thinkgo/raw/master/doc/index.png)

## 最新版本

### 版本号
v1.0

### 安装要求
Go Version ≥1.6

## 快速使用

### 框架下载

```sh
go get -u -v github.com/henrylee2cn/thinkgo
```

### 简单示例
```
package main

import (
    "github.com/henrylee2cn/thinkgo"
    "time"
)

type Index struct {
    Id        int      `param:"in(path),required,desc(ID),range(0:10)"`
    Title     string   `param:"in(query),nonzero"`
    Paragraph []string `param:"in(query),name(p),len(1:10)" regexp:"(^[\\w]*$)"`
    Cookie    string   `param:"in(cookie),name(thinkgoID)"`
    // Picture         multipart.FileHeader `param:"in(formData),name(pic),maxmb(30)"`
}

func (i *Index) Serve(ctx *thinkgo.Context) error {
    if ctx.CookieParam("thinkgoID") == "" {
        ctx.SetCookie("thinkgoID", time.Now().String())
    }
    return ctx.JSON(200, i)
}

func main() {
  app := thinkgo.New("myapp", "0.1")

  // Register the route in a chain style
  app.GET("/index/:id", new(Index))

  // Register the route in a tree style
  // app.Route(
  //   app.NewGET("/index/:id", new(Index)),
  // )

  // Start the service
  app.Run()
}

/*
http GET:
    http://localhost:8080/index/1?title=test&p=abc&p=xyz
response:
    {
      "Id": 1,
      "Title": "test",
      "Paragraph": [
        "abc",
        "xyz"
      ],
      "Cookie": "2016-11-13 01:14:40.9038005 +0800 CST"
    }
*/
```
[示例库](https://github.com/henrylee2cn/thinkgo/raw/master/samples)

## 框架特性

- 面向Handler接口开发（func or struct），中间件与操作完全等同可任意拼接路由操作链
- 支持用struct Handler在Tag标签定义请求参数信息及其校验信息
- 由struct Handler自动构建API文档（swagger2.0）
- 支持HTTP/HTTP2、HTTPS(tls/letsencrypt)、UNIX多种Server类型
- 支持多实例运行，且配置信息相互独立
- 支持同一实例监听多Server类型、多端口
- 基于著名的httprouter构建路由器，且支持链式与树形两种路由注册风格
- 强大的文件路由功能，支持自定义文件系统，框架提供快捷的DirFS、RenderFS、MarkdownFS等
- 提供近似LRU的文件缓存功能
- 跨平台的彩色日志系统，且同时支持console和file两种输出形式（可以同时使用）
- 提供Session管理功能
- 支持Gzip全局配置
- 提供XSRF跨站请求伪造安全过滤
- 简单整洁的配置文件，且自动补填默认值方便设置




## 配置文件说明

- 应用的各实例均有单独一份配置，其文件名格式 `config/{appname}[_{version}].ini`，配置详情：

```
net_types              = normal|tls              # 多种Server类型列表，支持 normal | tls | letsencrypt | unix
addrs                  = 0.0.0.0:80|0.0.0.0:443  # 多个监听地址列表
tls_certfile           =                         # TLS证书文件路径
tls_keyfile            =                         # TLS密钥文件路径
letsencrypt_file       =                         # SSL免费证书路径
unix_filemode          = 438                     # UNIX Server的文件权限（438即0666）
read_timeout           = 0                       # 读取请求数据超时
write_timeout          = 0                       # 写入响应数据超时
multipart_maxmemory_mb = 32                      # 接收上传文件时允许使用的最大内存

[router]                                         # 路由配置区
redirect_trailing_slash   = true                 # 当前请求的URL含`/`后缀如`/foo/`且相应路由不存在时，如存在`/foo`，则自动跳转至`/foo`
redirect_fixed_path       = true                 # 自动修复URL，如`/FOO` `/..//Foo`均被跳转至`/foo`（依赖redirect_trailing_slash=true）
handle_method_not_allowed = true                 # 若开启，当前请求方法不存在时返回405，否则返回404
handle_options            = true                 # 若开启，自动应答OPTIONS类请求，可在Thinkgo中设置默认Handler

[xsrf]                                           # XSRF跨站请求伪造过滤配置区
enable = false                                   # 是否开启
key    = thinkgoxsrf                             # 加密key
expire = 3600                                    # xsrf防伪token有效时长

[session]                                        # Session配置区（详情参考beego session模块）
enable                 = false                   # 是否开启
provider               = memory                  # 数据存储方式
name                   = thinkgosessionID        # 客户端存储cookie的名字
gc_max_lifetime        = 3600                    # 触发GC的时间
provider_config        =                         # 配置信息，根据不同的引擎设置不同的配置信息
cookie_lifetime        = 0                       # 客户端存储的cookie的时间，默认值是0，即浏览器生命周期
auto_setcookie         = true                    # 是否自动设置关于session的cookie值，一般默认true
domain                 =                         # 可以访问此cookie的域名
enable_sid_in_header   = false                   # 是否将session ID写入Header
name_in_header         = Thinkgosessionid        # 将session ID写入Header时的头名称
enable_sid_in_urlquery = false                   # 是否将session ID写入url的query部分

[apidoc]                                         # API文档
enable      = true                               # 是否启用
path        = /apidoc                            # 访问的URL路径
nolimit     = false                              # 是否不限访问IP
real_ip     = false                              # 使用真实客户端的IP进行过滤
whitelist   = 192.*|202.122.246.170              # 表示允许带有`192.`前缀或等于`202.122.246.170`的IP访问
desc        =                                    # 项目描述
email       =                                    # 联系人邮箱
terms_url   =                                    # 服务条款URL
license     =                                    # 协议类型
license_url =                                    # 协议内容URL
```

- 应用只有一份全局配置，文件名为 `config/__global__.ini`，配置详情：

```
[cache]                                          # 文件内存缓存配置区
enable  = false                                  # 是否开启
size_mb = 32                                     # 允许缓存使用的最大内存（单位MB），为0时系统自动设置为512KB
expire  = 60                                     # 缓存最大时长

[gzip]                                           # gzip压缩配置区
enable         = false                           # 是否开启
min_length     = 20                              # 进行压缩的最小内容长度
compress_level = 1                               # 非文件类响应Body的压缩水平（0-9），注意文件压缩始终为最优压缩比（9）
methods        = GET                             # 允许压缩的请求方法，为空时默认为GET

[log]                                            # 日志配置区
console_enable = true                            # 是否启用控制台日志
console_level  = debug                           # 控制台日志打印水平
file_enable    = true                            # 是否启用文件日志
file_level     = debug                           # 文件日志打印水平
```

## Handler结构体字段标签说明

tag   |   key    | required |     value     |   desc
------|----------|----------|---------------|----------------------------------
param |    in    | 有且只有一个 |     path      | （参数位置）为空时自动补全，如URL `http://www.abc.com/a/{path}`
param |    in    | 有且只有一个 |     query     | （参数位置）如URL `http://www.abc.com/a?b={query}`
param |    in    | 有且只有一个 |     formData  | （参数位置）请求表单，如 `a=123&b={formData}`
param |    in    | 有且只有一个 |     body      | （参数位置）请求Body
param |    in    | 有且只有一个 |     header    | （参数位置）请求头
param |    in    | 有且只有一个 |     cookie    | （参数位置）请求cookie，支持：`http.Cookie`、`fasthttp.Cookie`、`string`、`[]byte`等
param |   name   |      否      |     (如`id`)   | 自定义参数名
param | required |      否      |   required    | 参数是否必须
param |   desc   |      否      |   (e.g. `id`)  | 参数描述
param |   len    |      否      | (e.g. `3:6``3`) | 字符串类型参数的长度范围
param |   range  |      否      | (e.g. `0:10`)  | 数字类型参数的数值范围
param |  nonzero |      否      |    nonzero    | 是否能为零值
param |   maxmb  |      否      |  (e.g. `32`)   | 当前`Content-Type`为`multipart/form-data`时，允许使用的最大内存，当设置了多个时使用较大值
regexp|          |      否      | (e.g. `^\w+$`) | 使用正则验证参数值
err   |          |      否      |(e.g. `密码格式错误`)| 自定义参数绑定或验证的错误信息

**NOTES**:
* 绑定的对象必须为结构体指针类型
* 绑定的结构体字段类型不能为指针类型
* 只有在`param:"type(xxx)"`存在时，`regexp` 和 `param` 标签才有效
* 若`param`标签不存在，将尝试解析匿名字段
* 当结构体标签`in`为`formData`且字段类型为`multipart.FileHeader`时，该参数接收文件类型
* 当结构体标签`in`为`cookie`，字段类型必须为`http.Cookie`
* 标签`in(formData)`和`in(body)`不能同时出现在同一结构体
* 不能存在多个`in(body)`标签

## Handler结构体字段类型说明

base    |   slice    | special
--------|------------|-------------------------------------------------------
string  |  []string  | [][]byte
byte    |  []byte    | [][]uint8
uint8   |  []uint8   | multipart.FileHeader (仅`formData`参数使用)
bool    |  []bool    | http.Cookie (仅`net/http`下的`cookie`参数使用)
int     |  []int     | fasthttp.Cookie (仅`fasthttp`下的`cookie`参数使用)
int8    |  []int8    | struct (`body`参数使用或用于匿名字段扩展参数)
int16   |  []int16   |
int32   |  []int32   |
int64   |  []int64   |
uint8   |  []uint8   |
uint16  |  []uint16  |
uint32  |  []uint32  |
uint64  |  []uint64  |
float32 |  []float32 |
float64 |  []float64 |

## 扩展包
- [各种条码](https://github.com/henrylee2cn/thinkgo/raw/master/ext/barcode):       `github.com/henrylee2cn/thinkgo/ext/barcode`
- [比特单位](https://github.com/henrylee2cn/thinkgo/raw/master/ext/bitconv):       `github.com/henrylee2cn/thinkgo/ext/bitconv`
- [gorm数据库引擎](https://github.com/henrylee2cn/thinkgo/raw/master/ext/db/gorm): `github.com/henrylee2cn/thinkgo/ext/db/gorm`
- [sqlx数据库引擎](https://github.com/henrylee2cn/thinkgo/raw/master/ext/db/sqlx): `github.com/henrylee2cn/thinkgo/ext/db/sqlx`
- [xorm数据库引擎](https://github.com/henrylee2cn/thinkgo/raw/master/ext/db/xorm): `github.com/henrylee2cn/thinkgo/ext/db/xorm`
- [口令算法](https://github.com/henrylee2cn/thinkgo/raw/master/ext/otp):           `github.com/henrylee2cn/thinkgo/ext/otp`
- [UUID](https://github.com/henrylee2cn/thinkgo/raw/master/ext/uuid):              `github.com/henrylee2cn/thinkgo/ext/uuid`
- [Websocket](https://github.com/henrylee2cn/thinkgo/raw/master/ext/websocket):    `github.com/henrylee2cn/thinkgo/ext/websocket`
- [ini配置](https://github.com/henrylee2cn/thinkgo/raw/master/ini):                `github.com/henrylee2cn/thinkgo/ini`
- [定时器](https://github.com/henrylee2cn/thinkgo/raw/master/ext/cron):            `github.com/henrylee2cn/thinkgo/ext/cron`
- [任务工具](https://github.com/henrylee2cn/thinkgo/raw/master/ext/task):            `github.com/henrylee2cn/thinkgo/ext/task`


## 开源协议
Thinkgo 项目采用商业应用友好的 [Apache2.0](https://github.com/henrylee2cn/thinkgo/raw/master/LICENSE) 协议发布。
# Faygo    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/faygo)    ![Faygo goreportcard](https://goreportcard.com/badge/github.com/henrylee2cn/faygo)

![Faygo Favicon](https://github.com/henrylee2cn/faygo/raw/master/doc/faygo_96x96.png)

## 概述

Faygo 是一款快速、简洁的Go Web框架，可用极少的代码开发出高性能的Web应用程序（尤其是API接口）。只需定义 struct Handler，Faygo 就能自动绑定、验证请求参数并生成在线API文档。

官方QQ群：Go-Web 编程 42730308    [![Go-Web 编程群](http://pub.idqqimg.com/wpa/images/group.png)](http://jq.qq.com/?_wv=1027&k=fzi4p1)

[查看《用户手册》](https://github.com/henrylee2cn/faydoc)

![faygo index](https://github.com/henrylee2cn/faygo/raw/master/doc/index.png)

![faygo apidoc](https://github.com/henrylee2cn/faygo/raw/master/doc/apidoc.png)

![faygo server](https://github.com/henrylee2cn/faygo/raw/master/doc/server.png)

## 最新版本

### 版本号

v1.0

### 安装要求

Go Version ≥1.8

## 快速使用

- 方式一 源码下载

```sh
go get -u -v github.com/henrylee2cn/faygo
```

- 方式二 部署工具 （[Go to fay](https://github.com/henrylee2cn/fay)）

```sh
go get -u -v github.com/henrylee2cn/fay
```

```
        fay command [arguments]

The commands are:
        new        创建、编译和运行（监控文件变化）一个新的faygo项目
        run        编译和运行（监控文件变化）任意一个已存在的golang项目

fay new appname [apptpl]
        appname    指定新faygo项目的创建目录
        apptpl     指定一个faygo项目模板（可选）

fay run [appname]
        appname    指定待运行的golang项目路径（可选）
```

## 框架特性

- 一个 `struct Handler` 搞定多件事：
 * 定义 Handler/Middleware
 * 绑定与验证请求参数
 * 生成 Swagger2.0 API 在线文档
 * 数据库 ORM 映射

- Handler与Middleware完全相同，都是实现Handler接口（`func`或`struct`类型），共同构成路由操作链，只是概念层面的说法不同
- 支持多种网络类型：

网络类型                                      | 配置`net_types`值
----------------------------------------------|----------------
HTTP                                          | `http`
HTTPS/HTTP2(TLS)                              | `https`
HTTPS/HTTP2(Let's Encrypt TLS)                | `letsencrypt`
HTTPS/HTTP2(Let's Encrypt TLS on UNIX socket) | `unix_letsencrypt`
HTTP(UNIX socket)                             | `unix_http`
HTTPS/HTTP2(TLS on UNIX socket)               | `unix_https`

- 支持单服务单监听、单服务多监听、多服务多监听等，多个服务的配置信息相互独立
- 基于 `httprouter` 开发高性能路由，支持链式与树形两种注册风格，支持灵活的静态文件路由（如DirFS、RenderFS、MarkdownFS等）
- 支持平滑关闭、平滑升级，提供fay工具进行新建项目、热编译、元编程
- 采用最强大的 `pongo2` 作为HTML渲染引擎
- 提供近似LRU的文件缓存功能，主要用途是静态文件缓存
- 跨平台的彩色日志系统，且同时支持console和file两种输出形式（可以同时使用）
- 提供Session管理功能
- 支持Gzip全局配置
- 提供XSRF跨站请求伪造安全过滤
- 大多数功能尽量使用简洁的ini进行配置来避免不必要的重新编译，并且这些配置文件支持自动补填默认值
- 提供 `gorm`、`xorm`、`sqlx`、`directSQL`、`Websocket`、`ini` 、`http client` 等很多常用扩展包

![faygo struct handler 多重用途合一](https://github.com/henrylee2cn/faygo/raw/master/doc/MultiUsage.png)

## 简单示例

```go
package main

import (
    // "mime/multipart"
    "time"
    "github.com/henrylee2cn/faygo"
)

type Index struct {
    Id        int      `param:"<in:path> <required> <desc:ID> <range: 0:10>"`
    Title     string   `param:"<in:query> <nonzero>"`
    Paragraph []string `param:"<in:query> <name:p> <len: 1:10> <regexp: ^[\\w]*$>"`
    Cookie    string   `param:"<in:cookie> <name:faygoID>"`
    // Picture         *multipart.FileHeader `param:"<in:formData> <name:pic> <maxmb:30>"`
}

func (i *Index) Serve(ctx *faygo.Context) error {
    if ctx.CookieParam("faygoID") == "" {
        ctx.SetCookie("faygoID", time.Now().String())
    }
    return ctx.JSON(200, i)
}

func main() {
    app := faygo.New("myapp", "0.1")

    // Register the route in a chain style
    app.GET("/index/:id", new(Index))

    // Register the route in a tree style
    // app.Route(
    //     app.NewGET("/index/:id", new(Index)),
    // )

    // Start the service
    faygo.Run()
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

[示例库](https://github.com/henrylee2cn/faygo/raw/master/samples)

## 操作和中间件

操作和中间件是相同的，都是实现了Handler接口！

- 函数类型

```go
// 不含API文档描述
func Page() faygo.HandlerFunc {
    return func(ctx *faygo.Context) error {
        return ctx.String(200, "faygo")
    }
}

// 含API文档描述
var Page2 = faygo.WrapDoc(Page(), "测试页2的注意事项", "文本")
```

- 结构体类型

```go
// Param操作通过Tag绑定并验证请求参数
type Param struct {
    Id    int    `param:"<in:path> <required> <desc:ID> <range: 0:10>"`
    Title string `param:"<in:query>"`
}

// Serve实现Handler接口
func (p *Param) Serve(ctx *faygo.Context) error {
    return ctx.JSON(200,
        faygo.Map{
            "Struct Params":    p,
            "Additional Param": ctx.PathParam("additional"),
        }, true)
}

// Doc实现API文档接口（可选）
func (p *Param) Doc() faygo.Doc {
    return faygo.Doc{
        // 向API文档声明接口注意事项
        Note: "param desc",
        // 向API文档声明响应内容格式
        Return: faygo.JSONMsg{
            Code: 1,
            Info: "success",
        },
        // 向API文档增加额外的请求参数声明（可选）
        Params: []faygo.ParamInfo{
            {
                Name:  "additional",
                In:    "path",
                Model: "a",
                Desc:  "defined by the `Doc()` method",
            },
        },
    }
}
```

## 过滤函数

过滤函数必须是HandlerFunc类型！

```go
func Root2Index(ctx *faygo.Context) error {
    // 不允许直接访问`/index`
    if ctx.Path() == "/index" {
        ctx.Stop()
        return nil
    }
    if ctx.Path() == "/" {
        ctx.ModifyPath("/index")
    }
    return nil
}
```

## 路由注册

- 树状

```go
// 新建应用服务，参数：名称、版本
var app1 = faygo.New("myapp1", "1.0")

// 路由
app1.Filter(Root2Index).
    Route(
        app1.NewNamedGET("测试页1", "/page", Page()),
        app1.NewNamedGET("测试页2", "/page2", Page2),
        app1.NewGroup("home",
            app1.NewNamedGET("test param", "/param", &Param{
                // 为绑定的参数设定API文档中缺省值（可选）
                Id:    1,
                Title: "test param",
            }),
        ),
    )
```

- 链状

```go
// 新建应用服务，参数：名称、版本
var app2 = faygo.New("myapp2", "1.0")

// 路由
app2.Filter(Root2Index)
app2.NamedGET("test page", "/page", Page())
app2.NamedGET("test page2", "/page2", Page2)
app2.Group("home")
{
    app2.NamedGET("test param", "/param", &Param{
        // 为绑定的参数设定API文档中缺省值（可选）
        Id:    1,
        Title: "test param",
    })
}
```

## 平滑关闭与重启

- 平滑关闭

```sh
kill [pid]
```

- 平滑重启

```sh
kill -USR2 [pid]
```

## 配置文件说明

- 应用的各服务均有单独一份配置，其文件名格式 `config/{appname}[_{version}].ini`，配置详情：

```
net_types              = http|https              # 多种网络类型列表，支持 http | https | unix_http | unix_https | letsencrypt | unix_letsencrypt
addrs                  = 0.0.0.0:80|0.0.0.0:443  # 多个监听地址列表
tls_certfile           =                         # TLS证书文件路径
tls_keyfile            =                         # TLS密钥文件路径
letsencrypt_dir        =                         # Let's Encrypt TLS证书缓存目录
unix_filemode          = 438                     # UNIX listener的文件权限（438即0666）
read_timeout           = 0                       # 读取请求数据超时
write_timeout          = 0                       # 写入响应数据超时
multipart_maxmemory_mb = 32                      # 接收上传文件时允许使用的最大内存

[router]                                         # 路由配置区
redirect_trailing_slash   = true                 # 当前请求的URL含`/`后缀如`/foo/`且相应路由不存在时，如存在`/foo`，则自动跳转至`/foo`
redirect_fixed_path       = true                 # 自动修复URL，如`/FOO` `/..//Foo`均被跳转至`/foo`（依赖redirect_trailing_slash=true）
handle_method_not_allowed = true                 # 若开启，当前请求方法不存在时返回405，否则返回404
handle_options            = true                 # 若开启，自动应答OPTIONS类请求，可在Faygo中设置默认Handler

[xsrf]                                           # XSRF跨站请求伪造过滤配置区
enable = false                                   # 是否开启
key    = faygoxsrf                             # 加密key
expire = 3600                                    # xsrf防伪token有效时长

[session]                                        # Session配置区（详情参考beego session模块）
enable                 = false                   # 是否开启
provider               = memory                  # 数据存储方式
name                   = faygosessionID        # 客户端存储cookie的名字
provider_config        =                         # 配置信息，根据不同的引擎设置不同的配置信息
cookie_lifetime        = 0                       # 客户端存储的cookie的时间，默认值是0，即浏览器生命周期
gc_lifetime            = 300                     # 触发GC的时间
max_lifetime           = 3600                    # 会话的最大生命周期
auto_setcookie         = true                    # 是否自动设置关于session的cookie值，一般默认true
domain                 =                         # 可以访问此cookie的域名
enable_sid_in_header   = false                   # 是否将session ID写入Header
name_in_header         = Faygosessionid        # 将session ID写入Header时的头名称
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
async_len      = 0                               # 0表示同步打印，大于0表示异步缓存长度
```

## Handler结构体字段标签说明

tag   |   key    | required |     value     |   desc
------|----------|----------|---------------|----------------------------------
param |    in    | 有且只有一个 |     path      | （参数位置）为空时自动补全，如URL `http://www.abc.com/a/{path}`
param |    in    | 有且只有一个 |     query     | （参数位置）如URL `http://www.abc.com/a?b={query}`
param |    in    | 有且只有一个 |     formData  | （参数位置）请求表单，如 `a=123&b={formData}`
param |    in    | 有且只有一个 |     body      | （参数位置）请求Body
param |    in    | 有且只有一个 |     header    | （参数位置）请求头
param |    in    | 有且只有一个 |     cookie    | （参数位置）请求cookie，支持：`*http.Cookie`、`http.Cookie`、`string`、`[]byte`等
param |   name   |      否      |     (如`id`)   | 自定义参数名
param | required |      否      |   required    | 参数是否必须
param |   desc   |      否      |     (如`id`)   | 参数描述
param |   len    |      否      |   (如`3:6``3`)  | 字符串类型参数的长度范围
param |   range  |      否      |   (如`0:10`)   | 数字类型参数的数值范围
param |  nonzero |      否      |    nonzero    | 是否能为零值
param |   maxmb  |      否      |    (如`32`)    | 当前`Content-Type`为`multipart/form-data`时，允许使用的最大内存，当设置了多个时使用较大值
param |  regexp  |      否      |   (如`^\w+$`)  | 使用正则验证参数值
param |   err    |      否      |(如`密码格式错误`)| 自定义参数绑定或验证的错误信息

**NOTES**:
* 绑定的对象必须为结构体指针类型
* 除`*multipart.FileHeader`外，绑定的结构体字段类型不能为指针类型
* 只有在`param:"type(xxx)"`存在时，`regexp` 和 `param` 标签才有效
* 若`param`标签不存在，将尝试解析匿名字段
* 当结构体标签`in`为`formData`且字段类型为`*multipart.FileHeader`、`multipart.FileHeader`、`[]*multipart.FileHeader`或`[]multipart.FileHeader`时，该参数接收文件类型
* 当结构体标签`in`为`cookie`，字段类型必须为`*http.Cookie`或`http.Cookie`
* 标签`in(formData)`和`in(body)`不能同时出现在同一结构体
* 不能存在多个`in(body)`标签

## Handler结构体字段类型说明

base    |   slice    | special
--------|------------|-------------------------------------------------------
string  |  []string  | [][]byte
byte    |  []byte    | [][]uint8
uint8   |  []uint8   | *multipart.FileHeader (仅`formData`参数使用)
bool    |  []bool    | []*multipart.FileHeader (仅`formData`参数使用)
int     |  []int     | *http.Cookie (仅`net/http`下的`cookie`参数使用)
int8    |  []int8    | http.Cookie (仅`net/http`下的`cookie`参数使用)
int16   |  []int16   | struct (`body`参数使用或用于匿名字段扩展参数)
int32   |  []int32   |
int64   |  []int64   |
uint8   |  []uint8   |
uint16  |  []uint16  |
uint32  |  []uint32  |
uint64  |  []uint64  |
float32 |  []float32 |
float64 |  []float64 |

## 扩展包

扩展包           |  导入路径
-----------------|-----------------------------------------------------------------------------------------------------------------
[各种条码](https://github.com/henrylee2cn/faygo/raw/master/ext/barcode)       | `github.com/henrylee2cn/faygo/ext/barcode`
[比特单位](https://github.com/henrylee2cn/faygo/raw/master/ext/bitconv)       | `github.com/henrylee2cn/faygo/ext/bitconv`
[gorm数据库引擎](https://github.com/henrylee2cn/faygo/raw/master/ext/db/gorm) | `github.com/henrylee2cn/faygo/ext/db/gorm`
[sqlx数据库引擎](https://github.com/henrylee2cn/faygo/raw/master/ext/db/sqlx) | `github.com/henrylee2cn/faygo/ext/db/sqlx`
[xorm数据库引擎](https://github.com/henrylee2cn/faygo/raw/master/ext/db/xorm) | `github.com/henrylee2cn/faygo/ext/db/xorm`
[directSQL(配置化SQL引擎)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/directsql) | `github.com/henrylee2cn/faygo/ext/db/directsql`
[口令算法](https://github.com/henrylee2cn/faygo/raw/master/ext/otp)           | `github.com/henrylee2cn/faygo/ext/otp`
[UUID](https://github.com/henrylee2cn/faygo/raw/master/ext/uuid)              | `github.com/henrylee2cn/faygo/ext/uuid`
[Websocket](https://github.com/henrylee2cn/faygo/raw/master/ext/websocket)    | `github.com/henrylee2cn/faygo/ext/websocket`
[ini配置](https://github.com/henrylee2cn/faygo/raw/master/ini)                | `github.com/henrylee2cn/faygo/ini`
[定时器](https://github.com/henrylee2cn/faygo/raw/master/ext/cron)            | `github.com/henrylee2cn/faygo/ext/cron`
[任务工具](https://github.com/henrylee2cn/faygo/raw/master/ext/task)          | `github.com/henrylee2cn/faygo/ext/task`
[HTTP客户端](https://github.com/henrylee2cn/faygo/raw/master/ext/surfer)      | `github.com/henrylee2cn/faygo/ext/surfer`


## 开源协议

Faygo 项目采用商业应用友好的 [Apache2.0](https://github.com/henrylee2cn/faygo/raw/master/LICENSE) 协议发布

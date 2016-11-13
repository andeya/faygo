# Thinkgo    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/thinkgo)

# 概述
Thinkgo目前进行了全面重构，是一款面向中间件开发、智能参数映射与校验、自动化API文档的Go语言web框架。

官方QQ群：Go-Web 编程 42730308    [![Go-Web 编程群](http://pub.idqqimg.com/wpa/images/group.png)](http://jq.qq.com/?_wv=1027&k=fzi4p1)

# 框架下载

```sh
go get github.com/henrylee2cn/thinkgo
```

# 安装要求

Go Version ≥1.6

# 代码示例
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
  // Register the route in a chain style
  // thinkgo.GET("/index/:id", new(Index))

  // Register the route in a tree style
  thinkgo.Route(
    thinkgo.NewGET("/index/:id", new(Index)),
  )

  // Start the service
  thinkgo.Run()
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

# Handler结构体字段标签说明

tag   |   key    | required |     value     |   desc
------|----------|----------|---------------|----------------------------------
param |    in    | only one |     path      | (position of param) if `required` is unsetted, auto set it. e.g. url: "http://www.abc.com/a/{path}"
param |    in    | only one |     query     | (position of param) e.g. url: "http://www.abc.com/a?b={query}"
param |    in    | only one |     formData  | (position of param) e.g. "request body: a=123&b={formData}"
param |    in    | only one |     body      | (position of param) request body can be any content
param |    in    | only one |     header    | (position of param) request header info
param |    in    | only one |     cookie    | (position of param) request cookie info, support: `http.Cookie`, `fasthttp.Cookie`, `string`, `[]byte` and so on
param |   name   |    no    |  (e.g. "id")  | specify request param`s name
param | required |    no    |   required    | request param is required
param |   desc   |    no    |  (e.g. "id")  | request param description
param |   len    |    no    | (e.g. 3:6, 3) | length range of param's value
param |   range  |    no    |  (e.g. 0:10)  | numerical range of param's value
param |  nonzero |    no    |    nonzero    | param`s value can not be zero
param |   maxmb  |    no    |   (e.g. 32)   | when request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
regexp|          |    no    |(e.g. "^\\w+$")| param value can not be null
err   |          |    no    |(e.g. "incorrect password format")| customize the prompt for validation error

**NOTES**:
* the binding object must be a struct pointer
* the binding struct's field can not be a pointer
* `regexp` or `param` tag is only usable when `param:"type(xxx)"` is exist
* if the `param` tag is not exist, anonymous field will be parsed
* when the param's position(`in`) is `formData` and the field's type is `multipart.FileHeader`, the param receives file uploaded
* if param's position(`in`) is `cookie`, field's type must be `http.Cookie`
* param tags `in(formData)` and `in(body)` can not exist at the same time
* there should not be more than one `in(body)` param tag

# Handler结构体字段类型说明

base    |   slice    | special
--------|------------|-------------------------------------------------------
string  |  []string  | [][]byte
byte    |  []byte    | [][]uint8
uint8   |  []uint8   | multipart.FileHeader (only for `formData` param)
bool    |  []bool    | http.Cookie (only for `net/http`'s `cookie` param)
int     |  []int     | fasthttp.Cookie (only for `fasthttp`'s `cookie` param)
int8    |  []int8    | struct (struct type only for `body` param or as an anonymous field to extend params)
int16   |  []int16   |
int32   |  []int32   |
int64   |  []int64   |
uint8   |  []uint8   |
uint16  |  []uint16  |
uint32  |  []uint32  |
uint64  |  []uint64  |
float32 |  []float32 |
float64 |  []float64 |

## 开源协议
Lessgo 项目采用商业应用友好的 [Apache2.0](https://github.com/henrylee2cn/thinkgo/raw/master/LICENSE) 协议发布。

# Apiware    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/faygo/apiware)

Apiware binds the specified parameters of the Golang `net/http` and `fasthttp` requests to the structure and verifies the validity of the parameter values.

It is suggested that you can use the struct as the Handler of the web framework, and use the middleware to quickly bind the request parameters, saving a lot of parameter type conversion and validity verification. At the same time through the struct tag, create swagger json configuration file, easy to create api document services.

Apiware将Go语言`net/http`及`fasthttp`请求的指定参数绑定到结构体，并验证参数值的合法性。
建议您可以使用结构体作为web框架的Handler，并用该中间件快速绑定请求参数，节省了大量参数类型转换与有效性验证的工作。同时还可以通过该结构体标签，创建swagger的json配置文件，轻松创建api文档服务。

# Demo 示例

```
package main

import (
    "encoding/json"
    "github.com/henrylee2cn/faygo/apiware"
    // "mime/multipart"
    "net/http"
    "strings"
)

type TestApiware struct {
    Id           int         `param:"<in:path> <required> <desc:ID> <range: 1:2>"`
    Num          float32     `param:"<in:query> <name:n> <range: 0.1:10.19>"`
    Title        string      `param:"<in:query> <nonzero>"`
    Paragraph    []string    `param:"<in:query> <name:p> <len: 1:10> <regexp: ^[\\w]*$>"`
    Cookie       http.Cookie `param:"<in:cookie> <name:apiwareid>"`
    CookieString string      `param:"<in:cookie> <name:apiwareid>"`
    // Picture   multipart.FileHeader `param:"<in:formData> <name:pic> <maxmb:30>"`
}

var myApiware = apiware.New(pathdecoder, nil, nil)

var pattern = "/test/:id"

func pathdecoder(urlPath, pattern string) apiware.KV {
    idx := map[int]string{}
    for k, v := range strings.Split(pattern, "/") {
        if !strings.HasPrefix(v, ":") {
            continue
        }
        idx[k] = v[1:]
    }
    pathParams := make(map[string]string, len(idx))
    for k, v := range strings.Split(urlPath, "/") {
        name, ok := idx[k]
        if !ok {
            continue
        }
        pathParams[name] = v
    }
    return apiware.Map(pathParams)
}

func testHandler(resp http.ResponseWriter, req *http.Request) {
    // set cookies
    http.SetCookie(resp, &http.Cookie{
        Name:  "apiwareid",
        Value: "http_henrylee2cn",
    })

    // bind params
    params := new(TestApiware)
    err := myApiware.Bind(params, req, pattern)
    b, _ := json.MarshalIndent(params, "", " ")
    if err != nil {
        resp.WriteHeader(http.StatusBadRequest)
        resp.Write(append([]byte(err.Error()+"\n"), b...))
    } else {
        resp.WriteHeader(http.StatusOK)
        resp.Write(b)
    }
}

func main() {
    // Check whether `testHandler` meet the requirements of apiware, and register it
    err := myApiware.Register(new(TestApiware))
    if err != nil {
        panic(err)
    }

    // server
    http.HandleFunc("/test/0", testHandler)
    http.HandleFunc("/test/1", testHandler)
    http.HandleFunc("/test/1.1", testHandler)
    http.HandleFunc("/test/2", testHandler)
    http.HandleFunc("/test/3", testHandler)
    http.ListenAndServe(":8080", nil)
}
```

# Struct&Tag 结构体及其标签

tag   |   key    | required |     value     |   desc
------|----------|----------|---------------|----------------------------------
param |    in    | only one |     path      | (position of param) if `required` is unsetted, auto set it. e.g. url: "http://www.abc.com/a/{path}"
param |    in    | only one |     query     | (position of param) e.g. url: "http://www.abc.com/a?b={query}"
param |    in    | only one |     formData  | (position of param) e.g. "request body: a=123&b={formData}"
param |    in    | only one |     body      | (position of param) request body can be any content
param |    in    | only one |     header    | (position of param) request header info
param |    in    | only one |     cookie    | (position of param) request cookie info, support: `http.Cookie`,`fasthttp.Cookie`,`string`,`[]byte`
param |   name   |    no    |   (e.g.`id`)   | specify request param`s name
param | required |    no    |               | request param is required
param |   desc   |    no    |   (e.g.`id`)   | request param description
param |   len    |    no    | (e.g.`3:6` `3`) | length range of param's value
param |   range  |    no    |  (e.g.`0:10`)  | numerical range of param's value
param |  nonzero |    no    |               | param`s value can not be zero
param |   maxmb  |    no    |   (e.g.`32`)   | when request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
param |  regexp  |    no    | (e.g.`^\\w+$`) | verify the value of the param with a regular expression(param value can not be null)
param |   err    |    no    |(e.g.`incorrect password format`)| the custom error for binding or validating

**NOTES**:
* the binding object must be a struct pointer
* in addition to `*multipart.FileHeader`, the binding struct's field can not be a pointer
* `regexp` or `param` tag is only usable when `param:"type(xxx)"` is exist
* if the `param` tag is not exist, anonymous field will be parsed
* when the param's position(`in`) is `formData` and the field's type is `multipart.FileHeader`, the param receives file uploaded
* if param's position(`in`) is `cookie`, field's type must be `http.Cookie`
* param tags `in(formData)` and `in(body)` can not exist at the same time
* there should not be more than one `in(body)` param tag

# Field Types 结构体字段类型

base    |   slice    | special
--------|------------|-------------------------------------------------------
string  |  []string  | [][]byte
byte    |  []byte    | [][]uint8
uint8   |  []uint8   | *multipart.FileHeader (only for `formData` param)
bool    |  []bool    | []*multipart.FileHeader (only for `formData` param)
int     |  []int     | http.Cookie (only for `net/http`'s `cookie` param)
int8    |  []int8    | fasthttp.Cookie (only for `fasthttp`'s `cookie` param)
int16   |  []int16   | struct (struct type only for `body` param or as an anonymous field to extend params)
int32   |  []int32   |
int64   |  []int64   |
uint8   |  []uint8   |
uint16  |  []uint16  |
uint32  |  []uint32  |
uint64  |  []uint64  |
float32 |  []float32 |
float64 |  []float64 |

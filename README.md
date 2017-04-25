# Faygo    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/faygo)    ![Faygo goreportcard](https://goreportcard.com/badge/github.com/henrylee2cn/faygo)

![Faygo Favicon](https://github.com/henrylee2cn/faygo/raw/master/doc/faygo_96x96.png)

Faygo is a fast and concise Go Web framework that can be used to develop high-performance web app(especially API) with fewer codes. Just define a struct Handler, Faygo will automatically bind/verify the request parameters and generate the online API doc. [Go to \<User Manual\>](https://github.com/henrylee2cn/faydoc)

[简体中文](https://github.com/henrylee2cn/faygo/blob/master/README_ZH.md)

![faygo index](https://github.com/henrylee2cn/faygo/raw/master/doc/index.png)

![faygo apidoc](https://github.com/henrylee2cn/faygo/raw/master/doc/apidoc.png)

![faygo server](https://github.com/henrylee2cn/faygo/raw/master/doc/server.png)


## Latest version

### Version

v1.0

### Requirements

Go Version ≥1.8

## Quick Start

- Way 1: download source

```sh
go get -u -v github.com/henrylee2cn/faygo
```

- Way 2: deployment tools ([Go to fay](https://github.com/henrylee2cn/fay))

```sh
go get -u -v github.com/henrylee2cn/fay
```

```
        fay command [arguments]

The commands are:
        new        create, compile and run (monitor changes) a new faygo project
        run        compile and run (monitor changes) an any existing go project

fay new appname [apptpl]
        appname    specifies the path of the new faygo project
        apptpl     optionally, specifies the faygo project template type

fay run [appname]
        appname    optionally, specifies the path of the new project
```

## Features

- One `struct Handler` can get more things:
 * Define Handler/Middleware
 * Bind and verify request parameters
 * Generate an online document for the Swagger 2.0 API
 * Database ORM mapping

- Handler and Middleware are exactly the same, both implement the Handler interface (`func` or` struct`), which together constitute the handler chain of the router.
- Supports multiple network types:

Network types                                 | Configuration `net_types`
----------------------------------------------|----------------
HTTP                                          | `http`
HTTPS/HTTP2(TLS)                              | `https`
HTTPS/HTTP2(Let's Encrypt TLS)                | `letsencrypt`
HTTPS/HTTP2(Let's Encrypt TLS on UNIX socket) | `unix_letsencrypt`
HTTP(UNIX socket)                             | `unix_http`
HTTPS/HTTP2(TLS on UNIX socket)               | `unix_https`

- Support single-service & single-listener, single-service & multi-listener, multi-service & multi-listener and so on. The configuration of multiple services is independent of each other.
- The high-performance router based on `httprouter` supports both chain and tree registration styles; supports flexible static file router (such as DirFS, RenderFS, MarkdownFS, etc.).
- Support graceful shutdown and rebooting, provide fay tools which has new projects, hot compilation , meta programming function.
- Use the most powerful `pongo2` as the HTML rendering engine.
- Support near-LRU memory caching. (mainly used for static file cache)
- Support cross-platform color log system, and has two output interface(console and file).
- Support session management.
- Support global gzip compression configuration.
- Support XSRF security filtering.
- Most features try to use simple ini configurations to avoid unnecessary recompilation, and these profiles can be automatically assigned default values.
- Provide `gorm`, ` xorm`, `sqlx`, ` directSQL`, `Websocket`, ` ini`, `http client` and many other commonly used expansion packages.

![faygo handler multi-usage](https://github.com/henrylee2cn/faygo/raw/master/doc/MultiUsage.png)

## Simple example

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

[All samples](https://github.com/henrylee2cn/faygo/raw/master/samples)

## Handler and middleware

Handler and middleware are the same, both implemente Handler interface!

- function type

```go
// Page handler doesn't contains API doc description
func Page() faygo.HandlerFunc {
    return func(ctx *faygo.Context) error {
        return ctx.String(200, "faygo")
    }
}

// Page2 handler contains API doc description
var Page2 = faygo.WrapDoc(Page(), "test page2 notes", "test")
```

- struct type

```go
// Param binds and validates the request parameters by Tags
type Param struct {
    Id    int    `param:"<in:path> <required> <desc:ID> <range: 0:10>"`
    Title string `param:"<in:query>"`
}

// Serve implemente Handler interface
func (p *Param) Serve(ctx *faygo.Context) error {
    return ctx.JSON(200,
        faygo.Map{
            "Struct Params":    p,
            "Additional Param": ctx.PathParam("additional"),
        }, true)
}

// Doc implemente API Doc interface (optional)
func (p *Param) Doc() faygo.Doc {
    return faygo.Doc{
        // Add the API notes to the API doc
        Note: "param desc",
        // declare the response content format to the API doc
        Return: faygo.JSONMsg{
            Code: 1,
            Info: "success",
        },
        // additional request parameter declarations to the API doc (optional)
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

## Filter function

The filter function must be HandleFunc type!

```go
func Root2Index(ctx *faygo.Context) error {
    // Direct access to `/index` is not allowed
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

## Route registration

- tree style

```go
// New application object, params: name, version
var app1 = faygo.New("myapp1", "1.0")

// router
app1.Filter(Root2Index).
    Route(
        app1.NewNamedGET("test page", "/page", Page()),
        app1.NewNamedGET("test page2", "/page2", Page2),
        app1.NewGroup("home",
            app1.NewNamedGET("test param", "/param", &Param{
                // sets the default value in the API documentation for the request parameters (optional)
                Id:    1,
                Title: "test param",
            }),
        ),
    )
```

- chain style

```go
// New application object, params: name, version
var app2 = faygo.New("myapp2", "1.0")

// router
app2.Filter(Root2Index)
app2.NamedGET("test page", "/page", Page())
app2.NamedGET("test page2", "/page2", Page2)
app2.Group("home")
{
    app2.NamedGET("test param", "/param", &Param{
        // sets the default value in the API documentation for the request parameters(optional)
        Id:    1,
        Title: "test param",
    })
}
```

## Shutdown and reboot

- shutdown gracefully

```sh
kill [pid]
```

- reboot gracefully

```sh
kill -USR2 [pid]
```

## Configuration

- Each instance of the application has a single configuration (file name format `config/{appname}[_{version}].ini`). Refer to the following:

```
net_types              = http|https              # List of network type: http | https | unix_http | unix_https | letsencrypt | unix_letsencrypt
addrs                  = 0.0.0.0:80|0.0.0.0:443  # List of multiple listening addresses
tls_certfile           =                         # TLS certificate file path
tls_keyfile            =                         # TLS key file path
letsencrypt_dir        =                         # Let's Encrypt TLS certificate cache directory
unix_filemode          = 438                     # File permissions for UNIX listener (438 equivalent to 0666)
read_timeout           = 0                       # Maximum duration for reading the full request (including body)
write_timeout          = 0                       # Maximum duration for writing the full response (including body)
multipart_maxmemory_mb = 32                      # Maximum size of memory that can be used when receiving uploaded files

[router]                                         # Routing configuration section
redirect_trailing_slash   = true                 # Automatic redirection (for example, `/foo/` -> `/foo`)
redirect_fixed_path       = true                 # Tries to fix the current request path, if no handle is registered for it
handle_method_not_allowed = true                 # Returns 405 if the requested method does not exist, otherwise returns 404
handle_options            = true                 # Automatic response OPTIONS request, you can set the default Handler in Faygo

[xsrf]                                           # XSRF security section
enable = false                                   # Whether enabled or not
key    = faygoxsrf                             # Encryption key
expire = 3600                                    # Expire of XSRF token

[session]                                        # Session section
enable                 = false                   # Whether enabled or not
provider               = memory                  # Data storage
name                   = faygosessionID        # The client stores the name of the cookie
provider_config        =                         # According to the different engine settings different configuration information
cookie_lifetime        = 0                       # The default value is 0, which is the lifetime of the browser
gc_lifetime            = 300                     # The interval between triggering the GC
max_lifetime           = 3600                    # The session max lefetime
auto_setcookie         = true                    # Automatically set on the session cookie value, the general default true
domain                 =                         # The domain name that is allowed to access this cookie
enable_sid_in_header   = false                   # Whether to write a session ID to the header
name_in_header         = Faygosessionid        # The name of the header when the session ID is written to the header
enable_sid_in_urlquery = false                   # Whether to write the session ID to the URL Query params

[apidoc]                                         # API documentation section
enable      = true                               # Whether enabled or not
path        = /apidoc                            # The URL path
nolimit     = false                              # If true, access is not restricted
real_ip     = false                              # If true, means verifying the real IP of the visitor
whitelist   = 192.*|202.122.246.170              # `whitelist=192.*|202.122.246.170` means: only IP addresses that are prefixed with `192 'or equal to` 202.122.246.170' are allowed
desc        =                                    # Description of the application
email       =                                    # Technician's Email
terms_url   =                                    # Terms of service
license     =                                    # The license used by the API
license_url =                                    # The URL of the protocol content page
```

- Only one global configuration is applied (`config/__global__.ini`). Refer to the following:

```
[cache]                                          # Cache section
enable  = false                                  # Whether enabled or not
size_mb = 32                                     # Max size by MB for file cache, the cache size will be set to 512KB at minimum.
expire  = 60                                     # Maximum duration for caching

[gzip]                                           # compression section
enable         = false                           # Whether enabled or not
min_length     = 20                              # The minimum length of content to be compressed
compress_level = 1                               # Non-file response Body's compression level is 0-9, but the files' always 9
methods        = GET                             # List of HTTP methods to compress. If not set, only GET requests are compressed.

[log]                                            # Log section
console_enable = true                            # Whether enabled or not console logger
console_level  = debug                           # Console logger level
file_enable    = true                            # Whether enabled or not file logger
file_level     = debug                           # File logger level
async_len      = 0                               # The length of asynchronous buffer, 0 means synchronization
```

## Handler struct tags

tag   |   key    | required |     value     |   desc
------|----------|----------|---------------|----------------------------------
param |    in    | only one |     path      | (position of param) if `required` is unsetted, auto set it. e.g. url: "http://www.abc.com/a/{path}"
param |    in    | only one |     query     | (position of param) e.g. url: "http://www.abc.com/a?b={query}"
param |    in    | only one |     formData  | (position of param) e.g. "request body: a=123&b={formData}"
param |    in    | only one |     body      | (position of param) request body can be any content
param |    in    | only one |     header    | (position of param) request header info
param |    in    | only one |     cookie    | (position of param) request cookie info, support: `*http.Cookie`,`http.Cookie`,`string`,`[]byte`
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
* when the param's position(`in`) is `formData` and the field's type is `*multipart.FileHeader`, `multipart.FileHeader`, `[]*multipart.FileHeader` or `[]multipart.FileHeader`, the param receives file uploaded
* if param's position(`in`) is `cookie`, field's type must be `*http.Cookie` or `http.Cookie`
* param tags `in(formData)` and `in(body)` can not exist at the same time
* there should not be more than one `in(body)` param tag

## Handler struct fields type

base    |   slice    | special
--------|------------|-------------------------------------------------------
string  |  []string  | [][]byte
byte    |  []byte    | [][]uint8
uint8   |  []uint8   | *multipart.FileHeader (only for `formData` param)
bool    |  []bool    | []*multipart.FileHeader (only for `formData` param)
int     |  []int     | *http.Cookie (only for `net/http`'s `cookie` param)
int8    |  []int8    | http.Cookie (only for `net/http`'s `cookie` param)
int16   |  []int16   | struct (struct type only for `body` param or as an anonymous field to extend params)
int32   |  []int32   |
int64   |  []int64   |
uint8   |  []uint8   |
uint16  |  []uint16  |
uint32  |  []uint32  |
uint64  |  []uint64  |
float32 |  []float32 |
float64 |  []float64 |

## Expansion package

package summary  |  import path
-----------------|-----------------------------------------------------------------------------------------------------------------
[barcode](https://github.com/henrylee2cn/faygo/raw/master/ext/barcode)             | `github.com/henrylee2cn/faygo/ext/barcode`
[Bit unit conversion](https://github.com/henrylee2cn/faygo/raw/master/ext/bitconv) | `github.com/henrylee2cn/faygo/ext/bitconv`
[gorm(DB ORM)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/gorm)        | `github.com/henrylee2cn/faygo/ext/db/gorm`
[sqlx(DB ext)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/sqlx)        | `github.com/henrylee2cn/faygo/ext/db/sqlx`
[xorm(DB ORM)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/xorm)        | `github.com/henrylee2cn/faygo/ext/db/xorm`
[directSQL(Configured SQL engine)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/directsql) | `github.com/henrylee2cn/faygo/ext/db/directsql`
[One-time Password](https://github.com/henrylee2cn/faygo/raw/master/ext/otp)       | `github.com/henrylee2cn/faygo/ext/otp`
[UUID](https://github.com/henrylee2cn/faygo/raw/master/ext/uuid)                   | `github.com/henrylee2cn/faygo/ext/uuid`
[Websocket](https://github.com/henrylee2cn/faygo/raw/master/ext/websocket)         | `github.com/henrylee2cn/faygo/ext/websocket`
[ini](https://github.com/henrylee2cn/faygo/raw/master/ini)                         | `github.com/henrylee2cn/faygo/ini`
[cron](https://github.com/henrylee2cn/faygo/raw/master/ext/cron)                   | `github.com/henrylee2cn/faygo/ext/cron`
[task](https://github.com/henrylee2cn/faygo/raw/master/ext/task)                   | `github.com/henrylee2cn/faygo/ext/task`
[http client](https://github.com/henrylee2cn/faygo/raw/master/ext/surfer)          | `github.com/henrylee2cn/faygo/ext/surfer`


## License

Faygo is under Apache v2 License. See the [LICENSE](https://github.com/henrylee2cn/faygo/raw/master/LICENSE) file for the full license text

# Faygo    [![GoDoc](https://godoc.org/github.com/tsuna/gohbase?status.png)](https://godoc.org/github.com/henrylee2cn/faygo)

![Lessgo Favicon](https://github.com/henrylee2cn/faygo/raw/master/doc/faygo_96x96.png)

Faygo is a Golang Web framework that handler is middleware, supports intelligent parameter mapping and validation, and automates API documentation.

[简体中文](https://github.com/henrylee2cn/faygo/blob/master/README_ZH.md)

![faygo server](https://github.com/henrylee2cn/faygo/raw/master/doc/server.png)

![faygo apidoc](https://github.com/henrylee2cn/faygo/raw/master/doc/apidoc.png)

![faygo index](https://github.com/henrylee2cn/faygo/raw/master/doc/index.png)

## Quick Start

### Version requirements

Go Version ≥1.6

### Download and install

```sh
go get -u -v github.com/henrylee2cn/faygo
```

### Simple example
```
package main

import (
    "github.com/henrylee2cn/faygo"
    "time"
)

type Index struct {
    Id        int      `param:"in(path),required,desc(ID),range(0:10)"`
    Title     string   `param:"in(query),nonzero"`
    Paragraph []string `param:"in(query),name(p),len(1:10)" regexp:"(^[\\w]*$)"`
    Cookie    string   `param:"in(cookie),name(faygoID)"`
    // Picture         multipart.FileHeader `param:"in(formData),name(pic),maxmb(30)"`
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
[Full example](https://github.com/henrylee2cn/faygo/raw/master/demo)

## Features

- Handler interface oriented development (func or struct)
- Middleware and handler exactly the same, they together constitute the handler chain
- Supports the use of struct (implemented handler) tag tags to define request parameter information and its validation information
- The API documentation (swagger2.0) is automatically built by the handler
- Supports HTTP/HTTP2, HTTPS (tls/letsencrypt), UNIX and other Server types
- Multi-instance is supported, and these configurations information are independent of each other
- Supports the same instance to monitor multi-server and multi-port
- Based on the popular httprouter build router, and supports chain or tree style to register router
- Supports cross-platform color log system, and has two output interface (console and file)
- Supports session management
- Supports global gzip compression configuration
- Supports XSRF security filtering
- Supports near-LRU memory caching (mainly used for static file cache)
- Nice and easy to use configuration file, automatically write default values

# Configuration

- Each instance of the application has a single configuration (file name format `config/{appname}[_{version}].ini`). Refer to the following:

```
net_types              = normal|tls              # List of network type: normal | tls | letsencrypt | unix
addrs                  = 0.0.0.0:80|0.0.0.0:443  # List of multiple listening addresses
tls_certfile           =                         # TLS certificate file path
tls_keyfile            =                         # TLS key file path
letsencrypt_file       =                         # SSL free certificate path
unix_filemode          = 438                     # File permissions for UNIX Server (438 equivalent to 0666)
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
gc_max_lifetime        = 3600                    # The interval between triggering the GC
provider_config        =                         # According to the different engine settings different configuration information
cookie_lifetime        = 0                       # The default value is 0, which is the lifetime of the browser
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
```

# Handler struct tags

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

# Handler struct fields type

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

# Expansion package
- [barcode](https://github.com/henrylee2cn/faygo/raw/master/ext/barcode):       `github.com/henrylee2cn/faygo/ext/barcode`
- [Bit unit conversion](https://github.com/henrylee2cn/faygo/raw/master/ext/bitconv):       `github.com/henrylee2cn/faygo/ext/bitconv`
- [timer](https://github.com/henrylee2cn/faygo/raw/master/ext/cron):            `github.com/henrylee2cn/faygo/ext/cron`
- [gorm(DB ORM)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/gorm): `github.com/henrylee2cn/faygo/ext/db/gorm`
- [sqlx(DB ext)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/sqlx): `github.com/henrylee2cn/faygo/ext/db/sqlx`
- [xorm(DB ORM)](https://github.com/henrylee2cn/faygo/raw/master/ext/db/xorm): `github.com/henrylee2cn/faygo/ext/db/xorm`
- [One-time Password](https://github.com/henrylee2cn/faygo/raw/master/ext/otp):           `github.com/henrylee2cn/faygo/ext/otp`
- [UUID](https://github.com/henrylee2cn/faygo/raw/master/ext/uuid):              `github.com/henrylee2cn/faygo/ext/uuid`
- [Websocket](https://github.com/henrylee2cn/faygo/raw/master/ext/websocket):    `github.com/henrylee2cn/faygo/ext/websocket`
- [ini](https://github.com/henrylee2cn/faygo/raw/master/ini):                `github.com/henrylee2cn/faygo/ini`


# License
Faygo is under Apache v2 License. See the [LICENSE](https://github.com/henrylee2cn/faygo/raw/master/LICENSE) file for the full license text

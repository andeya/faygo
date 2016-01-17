
#ThinkGo Web Framework  [![GoDoc](https://godoc.org/github.com/henrylee2cn/thinkgo?status.svg)](https://godoc.org/github.com/henrylee2cn/thinkgo)

![ThinkGo Admin](https://github.com/henrylee2cn/thinkgo/raw/master/doc/favicon.png)

ThinkGo 是一款 Go 语言编写的 web 快速开发框架。它基于开源框架 Gin 进行二次开发，旨在实现一种类 ThinkPHP 的高可用、高效率的 web 框架。在此感谢 [Gin](https://github.com/gin-gonic/gin) 和 [httprouter](https://github.com/julienschmidt/httprouter)。它最显著的特点是模块、控制器、操作三段式的标准 MVC 架构，且模块与插件的目录结构完全一致，令开发变得非常简单灵活。

* 官方QQ群：Go-Web 编程 42730308    [![Go-Web 编程群](http://pub.idqqimg.com/wpa/images/group.png)](http://jq.qq.com/?_wv=1027&k=fzi4p1)

![ThinkGo Admin](https://github.com/henrylee2cn/thinkgo/raw/master/doc/server.jpg)

![ThinkGo Admin](https://github.com/henrylee2cn/thinkgo/raw/master/doc/admin.jpg)


##目录结构

```
├─core 框架目录
│ 
├─main.go 主文件
│ 
├─application 应用模块目录
│  ├─common 公共模块目录
│  │  ├─template.go 模板函数
│  │  ├─common.go 公共变量与函数
│  │  ├─controller 公共控制器类目录
│  │  ├─middleware 中间件目录
│  │  └─model 公共数据模型目录
│  │
│  ├─module.go 模块定义文件
│  ├─module 模块目录
│  │  ├─template.go 模板函数
│  │  ├─common.go 公共变量与函数
│  │  ├─controller.go 基础控制器
│  │  ├─controller 控制器目录
│  │  ├─model 模型目录
│  │  └─view 视图文件目录
│  │      └─default 主题文件目录
│  │          ├─__public__ 资源文件目录
│  │          └─xxx 控制器模板目录
│  │
│  └─... 扩展的可装卸功能模块或插件
│
├─deploy 部署文件目录
│
├─conf 配置文件目录
│
└─uploads 上传根目录
```

## 安装

1.下载框架源码
```sh
go get github.com/henrylee2cn/thinkgo
```

2.安装部署工具
```sh
go install
```

3.创建项目（在项目目录下运行cmd）
```sh
$ thinkgo new appname
```

4.以热编译模式运行（在项目目录下运行cmd）
```sh
$ thinkgo run
```

##使用说明

#### main.go

```go
package main

import (
    "github.com/henrylee2cn/thinkgo/core"
    
    _ "appname/application"
    _ "appname/application/common"
    _ "appname/deploy"
)

func main() {
    core.ThinkGo.
        // 以下为可选设置
        // 设置自定义的中间件列表
        // Use(...).
        // 必须调用的启动服务
        Run()
}
```

#### 定义模块

```go
package application

import (
    // "appname/application/common/middleware"
    _ "appname/application/home"
    . "appname/application/home/controller"
    "github.com/henrylee2cn/thinkgo/core"
)

func init() {
    core.ModulePrepare(&core.Module{
        Name:        "home",
        Class:       "模块示例",
        Description: "这是一个模块示例",
    }).SetThemes(
        // 自动设置传入的第1个主题为当前主题
        &core.Theme{
            Name:        "default",
            Description: "default",
            Src:         map[string]string{},
        },
    ).
        // 指定当前主题
        UseTheme("default").
        // 中间件
        //  Use(
        //  middleware.BasicAuth(middleware.Accounts{
        //      "foo":    "bar",
        //      "manu":   "4321",
        //  }),
        // ).
        // 注册路由
        GET("/index", &IndexController{})
}
```

#### 定义中间件

```go
package middleware

import (
    "log"
    "time"
    "github.com/henrylee2cn/thinkgo/core"
)
func Logger() core.HandlerFunc {
    return func(c *core.Context) {
        t := time.Now()

        // Set example variable
        c.Set("example", "12345")

        // before request

        c.Next()

        // after request
        latency := time.Since(t)
        log.Print(latency)

        // access the status we are sending
        status := c.Writer.Status()
        log.Println(status)
    }
}
```

#### 定义控制器

```go
package controller

import (
    "github.com/henrylee2cn/thinkgo/application/admin/common"
)

type IndexController struct {
    common.BaseController
}

func (this *IndexController) Index() {
    // 当路由规则为 `/admin/index/index?addon` 时，这样获取参数
    id := this.Query("addon")
    if id == "" {
    // 当路由规则为 `/admin/index/index/:mail` 时，这样获取参数
        id = this.Param("addon")
    }
    // 传入模板变量
    this.Data["name"] = "henrylee2cn"
    // 渲染模板并写回响应流
    this.HTML()
}
```

##FAQ

更多操作可以参考[Gin](https://github.com/gin-gonic/gin)的一些用法。

##开源协议

ThinkGo 项目采用商业应用友好的 [MIT](https://github.com/henrylee2cn/thinkgo/raw/master/doc/LICENSE) 协议发布。

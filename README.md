
#ThinkGo Web Framework  [![GoDoc](https://godoc.org/github.com/henrylee2cn/thinkgo?status.svg)](https://godoc.org/github.com/henrylee2cn/thinkgo)

![ThinkGo Admin](https://github.com/henrylee2cn/thinkgo/raw/master/doc/favicon.png)

ThinkGo 是一款 Go 语言编写的 web 快速开发框架。它基于开源框架 Gin 进行二次开发，旨在实现一种类 ThinkPHP 的高可用、高效率的 web 框架。在此感谢 [Gin](https://github.com/gin-gonic/gin) 和 [httprouter](https://github.com/julienschmidt/httprouter)。它最显著的特点是模块、控制器、操作三段式的标准 MVC 架构，且模块与插件的目录结构完全一致，令开发变得非常简单灵活。

* 官方QQ群：Go-Web 编程 42730308    [![Go-Web 编程群](http://pub.idqqimg.com/wpa/images/group.png)](http://jq.qq.com/?_wv=1027&k=fzi4p1)

![ThinkGo Admin](https://github.com/henrylee2cn/thinkgo/raw/master/doc/server.jpg)

![ThinkGo Admin](https://github.com/henrylee2cn/thinkgo/raw/master/doc/admin.jpg)


##目录结构

```
├─main.go 主文件
│ 
├─core 框架目录
├─application 应用模块目录
│  ├─example.go 模块定义文件
│  ├─example 模块目录
│  │  ├─conf 配置文件目录
│  │  ├─common 公共文件目录
│  │  ├─controller 控制器目录
│  │  ├─model 模型目录
│  │  └─view 视图文件目录
│  │      └─default 主题文件目录
│  │          ├─__public__ 资源文件目录
│  │          └─xxx 控制器模板目录
│  │
│  ├─example2.go 插件定义文件
│  ├─example2 插件目录
│  │  ├─conf 配置文件目录
│  │  ├─common 公共文件目录
│  │  ├─controller 控制器目录
│  │  ├─model 模型目录
│  │  └─view 视图文件目录
│  │      └─default 主题文件目录
│  │          ├─__public__ 资源文件目录
│  │          └─xxx 控制器模板目录
│  │
│  └─... 扩展的可装卸功能模块或插件
│
├─common 公共文件目录
│  ├─deploy 部署文件目录
│  └─utils 工具集
│
├─model 模型目录
├─conf 配置文件目录
└─uploads 上传根目录
```

## 安装

```sh
go get github.com/henrylee2cn/thinkgo
```
```
解压 application.zip 至项目根目录（示例）
```

##使用说明

#### main.go

```go
package main

import (
    _ "github.com/henrylee2cn/thinkgo/application"
    "github.com/henrylee2cn/thinkgo/core"
)

func main() {
    core.ThinkGoDefault().
    // 以下为可选设置
    TemplateDelims("{{{","}}}").
    TemplateSuffex(".html").
    TemplateFuncs(map[string]interface{}).
    Use(middleware1,middleware2,middleware3,...).
    // 启动服务
    Run(":8080")
}
```

#### 定义模块/插件

```go
package application

import (
    "github.com/henrylee2cn/thinkgo/application/admin/conf"
    . "github.com/henrylee2cn/thinkgo/application/admin/controller"
    . "github.com/henrylee2cn/thinkgo/core"
)

func init() {
    ModulePrepare(&Module{
    // 定义插件时改为: AddonPrepare(&Module{
        Name:        conf.NAME,
        Class:       conf.CLASS,
        Description: conf.DESCRIPTION,
    }).SetThemes(
        // 自动设置传入的第1个主题为当前主题
        &Theme{
            Name:        "default",
            Description: "default",
            Src:         map[string]string{},
        },
        &Theme{
            Name:        "blue",
            Description: "blue",
            Src:         map[string]string{
                "img":"/public/banner.jpg",
            },
        },
    ).
        // 指定当前主题
        UseTheme("blue").
        // 中间件
        // Use(...).
        // 注册路由，且可添加中间件
        GET("/index?addon", &IndexController{}).
        HEAD("/to", &IndexController{}).
        DELETE("/del/:addon", someHandlerFunc, &IndexController{}).
        POST("/add", &IndexController{}, someHandlerFunc).
        PUT("/ReadMail", someHandlerFunc, &ReadController{}, someHandlerFunc).
        PATCH("/ReadMail", someHandlerFunc, &ReadController{}, someHandlerFunc).
        OPTIONS("/ReadMail", someHandlerFunc, &ReadController{}, someHandlerFunc)
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
    // 当url类似 "/admin/index/index?addon=mail" 时，这样获取参数
    id := this.Query("addon")
    if id == "" {
    // 当url类似 "/admin/index/index/mail" 时，这样获取参数
        id = this.Param("addon")
    }
    // 传入模板变量
    this.Data["name"] = "henrylee2cn"
    // 渲染模板并写回响应流
    this.HTML()
}
```

##开源协议

ThinkGo 项目采用商业应用友好的 [MIT](https://github.com/henrylee2cn/thinkgo/raw/master/doc/LICENSE) 协议发布。

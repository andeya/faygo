package router

import (
	"net/http"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/handler"
)

func init() {
	thinkgo.Route(
		thinkgo.Group("home",
			thinkgo.NamedAPI("test", "GETPOST", "test/:id", &handler.Index{
				Paragraph: []string{"abc"},
				Returns: thinkgo.Returns{{
					Code:        200,
					Description: "成功",
				}, {
					Code:         400,
					Description:  "参数解析错误",
					ExampleValue: "error:???",
				}},
			}),
		),
		thinkgo.StaticFS("/public", http.Dir("./static/public")),
		thinkgo.Static("/public2", "./static/public"),
	)
}

func Route(frame *thinkgo.Framework) {
	frame.Route(
		frame.Group("home",
			frame.NamedAPI("test", "GETPOST", "test/:id", &handler.Index{
				Paragraph: []string{"abc"},
				Returns: thinkgo.Returns{{
					Code:        200,
					Description: "成功",
				}, {
					Code:         400,
					Description:  "参数解析错误",
					ExampleValue: "error:???",
				}},
			}),
		),
		frame.StaticFS("/public", http.Dir("./static/public")),
		frame.Static("/public2", "./static/public"),
	)
}

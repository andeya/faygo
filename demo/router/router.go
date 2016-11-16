package router

import (
	"net/http"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/handler"
)

func init() {
	// Register the route in a tree style
	thinkgo.Route(
		thinkgo.NewGroup("home",
			thinkgo.NewNamedAPI("test", "GET", "render", &handler.Render{}),
			thinkgo.NewNamedAPI("test", "GET POST", "test/:id", &handler.Index{
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
		thinkgo.NewStaticFS("/public", http.Dir("./static/public")),
		thinkgo.NewStatic("/public2", "./static/public"),
	)
}

func Route(frame *thinkgo.Framework) {
	// Register the route in a chain style
	home := frame.Group("home")
	{
		home.NamedAPI("test", "GETPOST", "test/:id", &handler.Index{
			Paragraph: []string{"abc"},
			Returns: thinkgo.Returns{{
				Code:        200,
				Description: "成功",
			}, {
				Code:         400,
				Description:  "参数解析错误",
				ExampleValue: "error:???",
			}},
		})
	}

	frame.StaticFS("/public", http.Dir("./static/public"))
	frame.Static("/public2", "./static/public")
}

package router

import (
	"net/http"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/handler"
)

func Route1(frame *thinkgo.Framework) {
	// Register the route in a tree style
	frame.Route(
		frame.NewGroup("home",
			frame.NewNamedAPI("test", "GET", "render", &handler.Render{}),
			frame.NewNamedAPI("test", "GET POST", "test/:id", &handler.Index{
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
		frame.NewNamedPOST("body的JSON绑定", "/body", &handler.Body{}),
		frame.NewStaticFS("/public", http.Dir("./static/public")),
		frame.NewStatic("/public2", "./static/public"),
	)
}

func Route2(frame *thinkgo.Framework) {
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

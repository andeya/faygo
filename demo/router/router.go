package router

import (
	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/handler"
	"github.com/henrylee2cn/thinkgo/demo/middleware"
)

// Register the route in a tree style
func Route1(frame *thinkgo.Framework) {
	frame.
		Filter(middleware.Root2Index).
		Route(
			// frame.NewNamedAPI("index", "*", "/", handler.Index()),
			frame.NewNamedAPI("index", "*", "/index", handler.Index()),
			frame.NewGroup("home",
				frame.NewNamedGET("html", "render", &handler.Render{}),
				frame.NewNamedAPI("params", "GET POST", "param/:id", &handler.Param{
					Paragraph: []string{"abc"},
				}),
			),
			frame.NewNamedGET("websocket", "/ws", handler.WebsocketPage()),
			frame.NewNamedGET("websocket_server", "/ws_server", handler.Websocket),
			frame.NewNamedPOST("binds the body in JSON format", "/body", &handler.Body{}),
			frame.NewStaticFS("/public", thinkgo.DirFS("./static/public")),
			frame.NewStatic("/syso", "../_syso"),
			frame.NewNamedStaticFS("renderfs test", "/renderfs", thinkgo.RenderFS(
				"./static/renderfs",
				".html", // "*"
				thinkgo.Map{"title": "RenderFS page"},
			)),
		)
}

// Register the route in a chain style
func Route2(frame *thinkgo.Framework) {
	frame.Filter(middleware.Root2Index)
	// frame.NamedAPI("index", "*", "/", handler.Index())
	frame.NamedAPI("index", "*", "/index", handler.Index())
	home := frame.Group("home")
	{
		home.NamedGET("html", "render", &handler.Render{})
		home.NamedAPI("params", "GET POST", "param/:id", &handler.Param{
			Paragraph: []string{"abc"},
		})
	}
	frame.NamedPOST("binds the body in JSON format", "/body", &handler.Body{})
	frame.StaticFS("/public", thinkgo.DirFS("./static/public"))
	frame.Static("/syso", "../_syso")

	frame.NamedStaticFS("renderfs test", "/renderfs", thinkgo.RenderFS(
		"./static/renderfs",
		".html", // "*"
		thinkgo.Map{"title": "RenderFS page"},
	))
}

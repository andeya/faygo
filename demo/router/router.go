package router

import (
	"net/http"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/demo/handler"
)

func Route1(frame *thinkgo.Framework) {
	// Register the route in a tree style
	frame.Route(
		frame.NewNamedAPI("index", "*", "/", handler.Index()),
		frame.NewGroup("home",
			frame.NewNamedGET("html", "render", &handler.Render{}),
			frame.NewNamedAPI("params", "GET POST", "param/:id", &handler.Param{
				Paragraph: []string{"abc"},
			}),
		),
		frame.NewNamedGET("websocket", "/ws", handler.WebsocketPage()),
		frame.NewNamedGET("websocket_server", "/ws_server", handler.Websocket),
		frame.NewNamedPOST("binds the body in JSON format", "/body", &handler.Body{}),
		frame.NewStaticFS("/public", http.Dir("./static/public")),
		frame.NewStatic("/syso", "../_syso"),
	)
}

func Route2(frame *thinkgo.Framework) {
	// Register the route in a chain style
	frame.NamedAPI("index", "*", "/", handler.Index())
	home := frame.Group("home")
	{
		home.NamedGET("html", "render", &handler.Render{})
		home.NamedAPI("params", "GET POST", "param/:id", &handler.Param{
			Paragraph: []string{"abc"},
		})
	}
	frame.NamedPOST("binds the body in JSON format", "/body", &handler.Body{})
	frame.StaticFS("/public", http.Dir("./static/public"))
	frame.Static("/syso", "../_syso")
}

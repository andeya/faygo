package router

import (
	"github.com/henrylee2cn/faygo"
	mw "github.com/henrylee2cn/faygo/ext/middleware"
	"github.com/henrylee2cn/faygo/samples/demo/handler"
	"github.com/henrylee2cn/faygo/samples/demo/middleware"
)

// Register the route in a tree style
func Route1(frame *faygo.Framework) {
	frame.
		Filter(middleware.Root2Index).
		Route(
			// frame.NewNamedAPI("index", "*", "/", handler.Index()),
			frame.NewNamedAPI("index", "*", "/index", handler.Index()),
			frame.NewGroup("home",
				frame.NewNamedGET("html", "render", &handler.Render{}),
				frame.NewNamedAPI("params", "GET POST", "param/:id/*additional", &handler.Param{
					Paragraph: []string{"default_paragraph"},
				}),
			),
			frame.NewNamedGET("websocket", "/ws", handler.WebsocketPage()),
			frame.NewNamedGET("websocket_server", "/ws_server", handler.Websocket),
			frame.NewNamedPOST("binds the body in JSON format", "/body", &handler.Body{}),
			frame.NewStaticFS("/public", faygo.DirFS("./static/public")),
			frame.NewStatic("/syso", "../../_syso"),
			frame.NewNamedStaticFS("render fs test", "/renderfs", faygo.RenderFS(
				"./static/renderfs",
				".html", // "*"
				faygo.Map{"title": "RenderFS page"},
			)).Use(mw.AutoHTMLSuffix()),
			frame.NewNamedStaticFS("markdown fs test", "/md", faygo.MarkdownFS(
				"./static/markdown",
			)),
			frame.NewNamedAPI("reverse proxy", "GETOPTIONS", "/search", handler.Search(0)).Use(mw.CrossOrigin),
		)
}

// Register the route in a chain style
func Route2(frame *faygo.Framework) {
	frame.Filter(middleware.Root2Index)
	// frame.NamedAPI("index", "*", "/", handler.Index())
	frame.NamedAPI("index", "*", "/index", handler.Index())
	home := frame.Group("home")
	{
		home.NamedGET("html", "render", &handler.Render{})
		home.NamedAPI("params", "GET POST", "param/:id/*additional", &handler.Param{
			Paragraph: []string{"abc"},
		})
	}
	frame.NamedPOST("binds the body in JSON format", "/body", &handler.Body{})
	frame.StaticFS("/public", faygo.DirFS("./static/public"))
	frame.Static("/syso", "../../_syso")

	frame.NamedAPI("reverse proxy", "GET OPTIONS", "/search", handler.Search(0)).Use(mw.CrossOrigin)

	frame.NamedStaticFS("render fs test", "/renderfs", faygo.RenderFS(
		"./static/renderfs",
		".html", // "*"
		faygo.Map{"title": "RenderFS page"},
	))

	frame.NamedStaticFS("markdown fs test", "/md", faygo.MarkdownFS(
		"./static/markdown",
	))
}

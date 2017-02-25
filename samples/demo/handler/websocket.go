package handler

import (
	"time"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/websocket"
)

func WebsocketPage() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		return ctx.Render(200, faygo.JoinStatic("websocket.html"), nil)
	}
}

var Websocket = faygo.WrapDoc(func(ctx *faygo.Context) error {
	var upgrader = websocket.Upgrader{}
	conn, err := upgrader.FayUpgrade(ctx, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		var req interface{}
		if err := conn.ReadJSON(&req); err != nil {
			ctx.Log().Warning("read:", err)
			return nil
		}
		ctx.Log().Info("req:", req)
		if err := conn.WriteJSON(map[string]string{"server_time": time.Now().String()}); err != nil {
			ctx.Log().Warning("write:", err)
			return nil
		}
	}
},
	"websocket example",
	map[string]string{"server_time": time.Now().String()},
)

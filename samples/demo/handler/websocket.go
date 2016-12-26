package handler

import (
	"time"

	"github.com/henrylee2cn/thinkgo"
	"github.com/henrylee2cn/thinkgo/ext/websocket"
)

func WebsocketPage() thinkgo.HandlerFunc {
	return func(ctx *thinkgo.Context) error {
		return ctx.Render(200, thinkgo.JoinStatic("websocket.html"), nil)
	}
}

var Websocket = thinkgo.DocWrap(
	thinkgo.HandlerFunc(func(ctx *thinkgo.Context) error {
		var upgrader = websocket.Upgrader{}
		conn, err := upgrader.ThinkUpgrade(ctx, nil)
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
	}),
	"websocket example",
	map[string]string{"server_time": time.Now().String()},
)

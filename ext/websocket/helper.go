package websocket

import (
	"github.com/henrylee2cn/thinkgo"
	"net/http"
)

// ThinkUpgrade upgrades the thinkgo server connection to the WebSocket protocol.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// application negotiated subprotocol (Sec-Websocket-Protocol).
//
// If the upgrade fails, then ThinkUpgrade replies to the client with an HTTP error
// response.
func (u *Upgrader) ThinkUpgrade(ctx *thinkgo.Context, responseHeader http.Header) (*Conn, error) {
	return u.Upgrade(ctx.W, ctx.R, responseHeader)
}

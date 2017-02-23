package websocket

import (
	"github.com/henrylee2cn/faygo"
	"net/http"
)

// FayUpgrade upgrades the faygo server connection to the WebSocket protocol.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// application negotiated subprotocol (Sec-Websocket-Protocol).
//
// If the upgrade fails, then FayUpgrade replies to the client with an HTTP error
// response.
func (u *Upgrader) FayUpgrade(ctx *faygo.Context, responseHeader http.Header) (*Conn, error) {
	return u.Upgrade(ctx.W, ctx.R, responseHeader)
}

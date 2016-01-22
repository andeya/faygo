package middleware

import (
	"encoding/base64"
	"net/http"

	"github.com/henrylee2cn/thinkgo/core"
)

type (
	BasicValidateFunc func(string, string) bool
)

const (
	Basic = "Basic"
)

// BasicAuth returns an HTTP basic authentication middleware.
//
// For valid credentials it calls the next handler.
// For invalid credentials, it sends "401 - Unauthorized" response.
func BasicAuth(fn BasicValidateFunc) core.HandlerFunc {
	return func(c *core.Context) error {
		// Skip WebSocket
		if (c.Request().Header.Get(core.Upgrade)) == core.WebSocket {
			return nil
		}

		auth := c.Request().Header.Get(core.Authorization)
		l := len(Basic)

		if len(auth) > l+1 && auth[:l] == Basic {
			b, err := base64.StdEncoding.DecodeString(auth[l+1:])
			if err == nil {
				cred := string(b)
				for i := 0; i < len(cred); i++ {
					if cred[i] == ':' {
						// Verify credentials
						if fn(cred[:i], cred[i+1:]) {
							return nil
						}
					}
				}
			}
		}
		c.Response().Header().Set(core.WWWAuthenticate, Basic+" realm=Restricted")
		return core.NewHTTPError(http.StatusUnauthorized)
	}
}

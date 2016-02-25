package core

import (
	"net"
	"time"

	// "github.com/henrylee2cn/thinkgo/core"
	"github.com/henrylee2cn/thinkgo/core/color"
)

func Logger() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			req := c.Request()
			res := c.Response()
			logger := c.Echo().Logger()

			remoteAddr := req.RemoteAddr
			if ip := req.Header.Get(XRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header.Get(XForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			method := req.Method
			path := req.URL.Path
			if path == "" {
				path = "/"
			}
			size := res.Size()

			n := res.Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Yellow(n)
			case n >= 300:
				code = color.Cyan(n)
			}

			logger.Info("%s %s %s %s %s %d", remoteAddr, method, path, code, stop.Sub(start), size)
			return nil
		}
	}
}

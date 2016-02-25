package core

import (
	"fmt"

	"runtime"
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() MiddlewareFunc {
	// TODO: Provide better stack trace `https://github.com/go-errors/errors` `https://github.com/docker/libcontainer/tree/master/stacktrace`
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			defer func() {
				if err := recover(); err != nil {
					trace := make([]byte, 1<<16)
					n := runtime.Stack(trace, true)
					c.Error(fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s",
						err, n, trace[:n]))
				}
			}()
			return next(c)
		}
	}
}

package middleware

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/henrylee2cn/thinkgo/core"
)

type (
	gzipWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

func (w gzipWriter) Write(b []byte) (int, error) {
	if w.Header().Get(core.ContentType) == "" {
		w.Header().Set(core.ContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w gzipWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *gzipWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

var writerPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip() core.MiddlewareFunc {
	scheme := "gzip"

	return func(h core.HandlerFunc) core.HandlerFunc {
		return func(c *core.Context) error {
			c.Response().Header().Add(core.Vary, core.AcceptEncoding)
			if strings.Contains(c.Request().Header.Get(core.AcceptEncoding), scheme) {
				w := writerPool.Get().(*gzip.Writer)
				w.Reset(c.Response().Writer())
				defer func() {
					w.Close()
					writerPool.Put(w)
				}()
				gw := gzipWriter{Writer: w, ResponseWriter: c.Response().Writer()}
				c.Response().Header().Set(core.ContentEncoding, scheme)
				c.Response().SetWriter(gw)
			}
			if err := h(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}

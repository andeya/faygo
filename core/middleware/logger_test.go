package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrylee2cn/thinkgo/core"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	// Note: Just for the test coverage, not a real test.
	e := core.New()
	req, _ := http.NewRequest(core.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := core.NewContext(req, core.NewResponse(rec, e), e)

	// Status 2xx
	h := func(c *core.Context) error {
		return c.String(http.StatusOK, "test")
	}
	Logger()(h)(c)

	// Status 3xx
	rec = httptest.NewRecorder()
	c = core.NewContext(req, core.NewResponse(rec, e), e)
	h = func(c *core.Context) error {
		return c.String(http.StatusTemporaryRedirect, "test")
	}
	Logger()(h)(c)

	// Status 4xx
	rec = httptest.NewRecorder()
	c = core.NewContext(req, core.NewResponse(rec, e), e)
	h = func(c *core.Context) error {
		return c.String(http.StatusNotFound, "test")
	}
	Logger()(h)(c)

	// Status 5xx with empty path
	req, _ = http.NewRequest(core.GET, "", nil)
	rec = httptest.NewRecorder()
	c = core.NewContext(req, core.NewResponse(rec, e), e)
	h = func(c *core.Context) error {
		return errors.New("error")
	}
	Logger()(h)(c)
}

func TestLoggerIPAddress(t *testing.T) {
	e := core.New()
	req, _ := http.NewRequest(core.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := core.NewContext(req, core.NewResponse(rec, e), e)
	buf := new(bytes.Buffer)
	e.Logger().SetOutput(buf)
	ip := "127.0.0.1"
	h := func(c *core.Context) error {
		return c.String(http.StatusOK, "test")
	}

	mw := Logger()

	// With X-Real-IP
	req.Header.Add(core.XRealIP, ip)
	mw(h)(c)
	assert.Contains(t, buf.String(), ip)

	// With X-Forwarded-For
	buf.Reset()
	req.Header.Del(core.XRealIP)
	req.Header.Add(core.XForwardedFor, ip)
	mw(h)(c)
	assert.Contains(t, buf.String(), ip)

	// with req.RemoteAddr
	buf.Reset()
	mw(h)(c)
	assert.Contains(t, buf.String(), ip)
}

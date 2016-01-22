package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrylee2cn/thinkgo/core"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	e := core.New()
	e.SetDebug(true)
	req, _ := http.NewRequest(core.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := core.NewContext(req, core.NewResponse(rec, e), e)
	h := func(c *core.Context) error {
		panic("test")
	}
	Recover()(h)(c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "panic recover")
}

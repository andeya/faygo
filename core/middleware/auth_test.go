package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrylee2cn/thinkgo/core"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	e := core.New()
	req, _ := http.NewRequest(core.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := core.NewContext(req, core.NewResponse(rec, e), e)
	fn := func(u, p string) bool {
		if u == "joe" && p == "secret" {
			return true
		}
		return false
	}
	ba := BasicAuth(fn)

	// Valid credentials
	auth := Basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))
	req.Header.Set(core.Authorization, auth)
	assert.NoError(t, ba(c))

	//---------------------
	// Invalid credentials
	//---------------------

	// Incorrect password
	auth = Basic + " " + base64.StdEncoding.EncodeToString([]byte("joe:password"))
	req.Header.Set(core.Authorization, auth)
	he := ba(c).(*core.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, he.Code())
	assert.Equal(t, Basic+" realm=Restricted", rec.Header().Get(core.WWWAuthenticate))

	// Empty Authorization header
	req.Header.Set(core.Authorization, "")
	he = ba(c).(*core.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, he.Code())
	assert.Equal(t, Basic+" realm=Restricted", rec.Header().Get(core.WWWAuthenticate))

	// Invalid Authorization header
	auth = base64.StdEncoding.EncodeToString([]byte("invalid"))
	req.Header.Set(core.Authorization, auth)
	he = ba(c).(*core.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, he.Code())
	assert.Equal(t, Basic+" realm=Restricted", rec.Header().Get(core.WWWAuthenticate))

	// WebSocket
	c.Request().Header.Set(core.Upgrade, core.WebSocket)
	assert.NoError(t, ba(c))
}

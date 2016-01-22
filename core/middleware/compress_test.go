package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/henrylee2cn/thinkgo/core"
	"github.com/stretchr/testify/assert"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

func TestGzip(t *testing.T) {
	e := core.New()
	req, _ := http.NewRequest(core.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := core.NewContext(req, core.NewResponse(rec, e), e)
	h := func(c *core.Context) error {
		c.Response().Write([]byte("test")) // For Content-Type sniffing
		return nil
	}

	// Skip if no Accept-Encoding header
	Gzip()(h)(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())

	req, _ = http.NewRequest(core.GET, "/", nil)
	req.Header.Set(core.AcceptEncoding, "gzip")
	rec = httptest.NewRecorder()
	c = core.NewContext(req, core.NewResponse(rec, e), e)

	// Gzip
	Gzip()(h)(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "gzip", rec.Header().Get(core.ContentEncoding))
	assert.Contains(t, rec.Header().Get(core.ContentType), core.TextPlain)
	r, err := gzip.NewReader(rec.Body)
	defer r.Close()
	if assert.NoError(t, err) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		assert.Equal(t, "test", buf.String())
	}
}

func TestGzipFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	buf := new(bytes.Buffer)
	w := gzip.NewWriter(buf)
	gw := gzipWriter{Writer: w, ResponseWriter: rec}

	n0 := buf.Len()
	if n0 != 0 {
		t.Fatalf("buffer size = %d before writes; want 0", n0)
	}

	if err := gw.Flush(); err != nil {
		t.Fatal(err)
	}

	n1 := buf.Len()
	if n1 == 0 {
		t.Fatal("no data after first flush")
	}

	gw.Write([]byte("x"))

	n2 := buf.Len()
	if n1 != n2 {
		t.Fatalf("after writing a single byte, size changed from %d to %d; want no change", n1, n2)
	}

	if err := gw.Flush(); err != nil {
		t.Fatal(err)
	}

	n3 := buf.Len()
	if n2 == n3 {
		t.Fatal("Flush didn't flush any data")
	}
}

func TestGzipCloseNotify(t *testing.T) {
	rec := newCloseNotifyingRecorder()
	buf := new(bytes.Buffer)
	w := gzip.NewWriter(buf)
	gw := gzipWriter{Writer: w, ResponseWriter: rec}
	closed := false
	notifier := gw.CloseNotify()
	rec.close()

	select {
	case <-notifier:
		closed = true
	case <-time.After(time.Second):
	}

	assert.Equal(t, closed, true)
}

func BenchmarkGzip(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	h := func(c *core.Context) error {
		c.Response().Write([]byte("test")) // For Content-Type sniffing
		return nil
	}
	req, _ := http.NewRequest(core.GET, "/", nil)
	req.Header.Set(core.AcceptEncoding, "gzip")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		e := core.New()
		rec := httptest.NewRecorder()
		c := core.NewContext(req, core.NewResponse(rec, e), e)
		Gzip()(h)(c)
	}

}

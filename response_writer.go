// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package faygo

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"runtime"
)

// Response wraps an http.ResponseWriter and implements its interface to be used
// by an HTTP handler to construct an HTTP response.
// See [http.ResponseWriter](https://golang.org/pkg/net/http/#ResponseWriter)
type Response struct {
	context   *Context
	writer    http.ResponseWriter
	status    int
	size      int64
	committed bool
}

var _ http.ResponseWriter = new(Response)

func (resp *Response) reset(w http.ResponseWriter) {
	resp.writer = w
	resp.status = 0
	resp.size = 0
	resp.committed = false
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (resp *Response) Header() http.Header {
	return resp.writer.Header()
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (resp *Response) WriteHeader(status int) {
	if resp.committed {
		resp.multiCommitted()
		return
	}
	resp.status = status
	resp.context.beforeWriteHeader()
	resp.writer.WriteHeader(status)
	resp.committed = true
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (resp *Response) Write(b []byte) (int, error) {
	if !resp.committed {
		resp.WriteHeader(200)
	}
	n, err := resp.writer.Write(b)
	resp.size += int64(n)
	return n, err
}

// AddCookie adds a Set-Cookie header.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (resp *Response) AddCookie(cookie *http.Cookie) {
	resp.Header().Add(HeaderSetCookie, cookie.String())
}

// SetCookie sets a Set-Cookie header.
func (resp *Response) SetCookie(cookie *http.Cookie) {
	resp.Header().Set(HeaderSetCookie, cookie.String())
}

// DelCookie sets Set-Cookie header.
func (resp *Response) DelCookie() {
	resp.Header().Del(HeaderSetCookie)
}

// ReadFrom is here to optimize copying from an *os.File regular file
// to a *net.TCPConn with sendfile.
func (resp *Response) ReadFrom(src io.Reader) (int64, error) {
	if rf, ok := resp.writer.(io.ReaderFrom); ok {
		n, err := rf.ReadFrom(src)
		resp.size += int64(n)
		return n, err
	}
	var buf = make([]byte, 32*1024)
	var n int64
	var err error
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := resp.writer.Write(buf[0:nr])
			if nw > 0 {
				n += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	resp.size += n
	return n, err
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
func (resp *Response) Flush() {
	if f, ok := resp.writer.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
func (resp *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := resp.writer.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("webserver doesn't support Hijack")
}

// CloseNotify implements the http.CloseNotifier interface to allow detecting
// when the underlying connection has gone away.
// This mechanism can be used to cancel long operations on the server if the
// client has disconnected before the response is ready.
func (resp *Response) CloseNotify() <-chan bool {
	if cn, ok := resp.writer.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}

// Size returns the current size, in bytes, of the response.
func (resp *Response) Size() int64 {
	return resp.size
}

// Committed returns whether the response has been submitted or not.
func (resp *Response) Committed() bool {
	return resp.committed
}

// Status returns the HTTP status code of the response.
func (resp *Response) Status() int {
	return resp.status
}

func (resp *Response) multiCommitted() {
	if resp.status == 200 {
		line := []byte("\n")
		e := []byte("\ngoroutine ")
		stack := make([]byte, 2<<10) //2KB
		runtime.Stack(stack, true)
		start := bytes.Index(stack, line) + 1
		stack = stack[start:]
		end := bytes.LastIndex(stack, line)
		if end != -1 {
			stack = stack[:end]
		}
		end = bytes.Index(stack, e)
		if end != -1 {
			stack = stack[:end]
		}
		stack = bytes.TrimRight(stack, "\n")
		resp.context.Log().Warningf("multiple response.WriteHeader calls\n[TRACE]\n%s\n", stack)
	}
}

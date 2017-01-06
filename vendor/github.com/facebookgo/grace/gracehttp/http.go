// Package gracehttp provides easy to use graceful restart
// functionality for HTTP server.
// modified by henrylee 2016.10.29
package gracehttp

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/facebookgo/grace/gracenet"
	"github.com/facebookgo/httpdown"
)

var (
	verbose    = flag.Bool("gracehttp.log", true, "Enable logging.")
	didInherit = os.Getenv("LISTEN_FDS") != ""
	ppid       = os.Getppid()
)

// An app contains one or more servers and associated configuration.
type app struct {
	servers       []*http.Server
	http          *httpdown.HTTP
	net           *gracenet.Net
	listeners     []net.Listener
	sds           []httpdown.Server
	errors        chan error
	terminateFunc func() error
}

func newApp(servers []*http.Server) *app {
	return &app{
		servers:   servers,
		http:      &httpdown.HTTP{},
		net:       &gracenet.Net{},
		listeners: make([]net.Listener, 0, len(servers)),
		sds:       make([]httpdown.Server, 0, len(servers)),

		// 2x num servers for possible Close or Stop errors + 1 for possible
		// StartProcess error + 1 for possible terminateFunc error.
		errors: make(chan error, 2+(len(servers)*2)),
	}
}

func (a *app) listen() error {
	for _, s := range a.servers {
		// TODO: default addresses
		l, err := a.net.Listen("tcp", s.Addr)
		if err != nil {
			return err
		}
		if s.TLSConfig != nil {
			l = tls.NewListener(l, s.TLSConfig)
		}
		a.listeners = append(a.listeners, l)
	}
	return nil
}

func (a *app) serve() {
	for i, s := range a.servers {
		a.sds = append(a.sds, a.http.Serve(s, a.listeners[i]))
	}
}

func (a *app) wait() {
	var wg sync.WaitGroup
	wg.Add(len(a.sds) * 2) // Wait & Stop
	go a.signalHandler(&wg)
	for _, s := range a.sds {
		go func(s httpdown.Server) {
			defer wg.Done()
			if err := s.Wait(); err != nil {
				a.errors <- err
			}
		}(s)
	}
	wg.Wait()
}

func (a *app) term(wg *sync.WaitGroup) {
	for _, s := range a.sds {
		go func(s httpdown.Server) {
			defer wg.Done()
			if err := s.Stop(); err != nil {
				a.errors <- err
			}
		}(s)
	}
}

// Serve will serve the given http.Servers and will monitor for signals
// allowing for graceful termination (SIGTERM) or restart (SIGUSR2).
func Serve(servers ...*http.Server) error {
	return ServeWithTerminateFunc(nil, servers...)
}

// Used for pretty printing addresses.
func pprintAddr(listeners []net.Listener) []byte {
	var out bytes.Buffer
	for i, l := range listeners {
		if i != 0 {
			fmt.Fprint(&out, ", ")
		}
		fmt.Fprint(&out, l.Addr())
	}
	return out.Bytes()
}

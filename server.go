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
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/henrylee2cn/faygo/errors"
	"github.com/henrylee2cn/faygo/gracenet"
	"github.com/henrylee2cn/faygo/logging"
	"golang.org/x/crypto/acme/autocert"
)

// network types
const (
	// listenAndServe listens on the TCP network address and then
	// calls Serve to handle requests on incoming connections.
	// Accepted connections are configured to enable TCP keep-alives.
	// If srv.Addr is blank, ":http" is used.
	NETTYPE_HTTP = "http"
	// NETTYPE_HTTPS listens on the TCP network address and
	// then calls Serve to handle requests on incoming TLS connections.
	// Accepted connections are configured to enable TCP keep-alives.
	//
	// Filenames containing a certificate and matching private key for the
	// server must be provided if neither the Server's TLSConfig.Certificates
	// nor TLSConfig.GetCertificate are populated. If the certificate is
	// signed by a certificate authority, the certFile should be the
	// concatenation of the server's certificate, any intermediates, and
	// the CA's certificate.
	//
	// If server.Addr is blank, ":https" is used.
	NETTYPE_HTTPS = "https"
	// NETTYPE_LETSENCRYPT listens on a new Automatic TLS using letsencrypt.org service.
	// if you want to disable cache directory then simple give config `letsencrypt_dir` a value of empty string "".
	//
	// If server.Addr is blank, ":https" is used.
	NETTYPE_LETSENCRYPT = "letsencrypt"
	// NETTYPE_UNIX_LETSENCRYPT listens on a new Automatic TLS using letsencrypt.org Unix service.
	// if you want to disable cache directory then simple give config `letsencrypt_dir` a value of empty string "".
	//
	// If server.Addr is blank, ":https" is used.
	NETTYPE_UNIX_LETSENCRYPT = "unix_letsencrypt"
	// NETTYPE_UNIX_HTTP announces on the Unix domain socket addr and listens a Unix service.
	//
	// If server.Addr is blank, ":http" is used.
	NETTYPE_UNIX_HTTP = "unix_http"
	// NETTYPE_UNIX_HTTPS announces on the Unix domain socket addr and listens a secure Unix service.
	//
	// If server.Addr is blank, ":https" is used.
	NETTYPE_UNIX_HTTPS = "unix_https"

	__netTypes__ = "http | https | unix_http | unix_https | letsencrypt | unix_letsencrypt"
)

// Server web server object
type Server struct {
	nameWithVersion string
	netType         string
	net             string
	tlsCertFile     string
	tlsKeyFile      string
	letsencryptDir  string
	unixFileMode    os.FileMode
	*http.Server
	log *logging.Logger
}

func (server *Server) run() {
	server.initAddr()
	server.setNet()
	ln := server.listen()

	typ := strings.ToUpper(server.netType)
	switch server.netType {
	case NETTYPE_HTTPS, NETTYPE_UNIX_HTTPS, NETTYPE_LETSENCRYPT, NETTYPE_UNIX_LETSENCRYPT:
		typ += "/HTTP2"
	}
	server.log.Criticalf("\x1b[46m[SYS]\x1b[0m listen and serve %s on %v", typ, server.Addr)

	err := server.Server.Serve(ln)
	if realServeError(err) != nil {
		server.log.Fatalf("%v\n", err)
	}
}

func (server *Server) initAddr() {
	switch server.netType {
	case NETTYPE_HTTP, NETTYPE_UNIX_HTTP:
		if server.Addr == "" {
			server.Addr = ":http"
		}
	case NETTYPE_HTTPS, NETTYPE_UNIX_HTTPS, NETTYPE_LETSENCRYPT, NETTYPE_UNIX_LETSENCRYPT:
		if server.Addr == "" {
			server.Addr = ":https"
		}
	}
}

func (server *Server) setNet() {
	switch server.netType {
	case NETTYPE_HTTP, NETTYPE_HTTPS, NETTYPE_LETSENCRYPT:
		server.net = "tcp"
	case NETTYPE_UNIX_HTTP, NETTYPE_UNIX_HTTPS, NETTYPE_UNIX_LETSENCRYPT:
		server.net = "unix"
	default:
		server.log.Fatalf("Please set a valid config item net_type, refer to the following:\n%s\n", __netTypes__)
	}
}

var (
	errRemoveUnix = errors.New("[NET:UNIX] Unexpected error when trying to remove unix socket file. Addr: %s | Trace: %s")
	errChmod      = errors.New("[NET:UNIX] Cannot chmod %#o for %q: %s")
)

var grace = new(gracenet.Net)

func (server *Server) listen() net.Listener {
	switch server.netType {
	case NETTYPE_HTTPS, NETTYPE_UNIX_HTTPS:
		var cert tls.Certificate
		cert, err := tls.LoadX509KeyPair(server.tlsCertFile, server.tlsKeyFile)
		if err != nil {
			server.log.Fatalf("%v\n", err)
			return nil
		}
		server.TLSConfig = &tls.Config{
			Certificates:             []tls.Certificate{cert},
			NextProtos:               []string{"http/1.1", "h2"},
			PreferServerCipherSuites: true,
		}

	case NETTYPE_LETSENCRYPT, NETTYPE_UNIX_LETSENCRYPT:
		m := autocert.Manager{
			Prompt: autocert.AcceptTOS,
		}

		if server.letsencryptDir == "" {
			// then the user passed empty by own will, then I guess user doesnt' want any cache directory
		} else {
			m.Cache = autocert.DirCache(server.letsencryptDir)
		}
		server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
	}

	switch server.netType {
	case NETTYPE_UNIX_HTTPS, NETTYPE_UNIX_LETSENCRYPT:
		if errOs := os.Remove(server.Addr); errOs != nil && !os.IsNotExist(errOs) {
			server.log.Fatalf("%v\n", errRemoveUnix.Format(server.Addr, errOs.Error()))
			return nil
		}
		defer func() {
			err := os.Chmod(server.Addr, server.unixFileMode)
			if err != nil {
				server.log.Fatalf("%v\n", errChmod.Format(server.unixFileMode, server.Addr, err.Error()))
			}
		}()
	}

	ln, err := grace.Listen(server.net, server.Addr)
	if err != nil {
		server.log.Fatalf("%v\n", err)
		return nil
	}
	ln = tcpKeepAliveListener{ln.(*net.TCPListener)}
	if server.TLSConfig != nil {
		ln = tls.NewListener(ln, server.TLSConfig)
	}

	return ln
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func realServeError(err error) error {
	if err != nil && err == http.ErrServerClosed {
		return nil
	}
	return err
}

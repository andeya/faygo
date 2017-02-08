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

package thinkgo

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/henrylee2cn/thinkgo/errors"
	"github.com/henrylee2cn/thinkgo/gracenet"
	"github.com/henrylee2cn/thinkgo/logging"
	"github.com/rsc/letsencrypt"
	// "github.com/facebookgo/grace/gracehttp"
)

// Server web server object
type Server struct {
	nameWithVersion string
	netType         string
	net             string
	tlsCertFile     string
	tlsKeyFile      string
	letsencryptFile string
	unixFileMode    os.FileMode
	*http.Server
	log *logging.Logger
}

func (server *Server) run() {
	var err error
	switch server.netType {
	case NETTYPE_NORMAL:
		server.net = "tcp"
		err = server.listenAndServe()
	case NETTYPE_TLS:
		server.net = "tcp"
		err = server.listenAndServeTLS()
	case NETTYPE_LETSENCRYPT:
		server.net = "tcp"
		err = server.listenAndServeLETSENCRYPT()
	case NETTYPE_UNIX:
		server.net = "unix"
		err = server.listenAndServeUNIX()
	default:
		server.log.Fatal("Please set a valid config item net_type, refer to the following:\nnormal | tls | letsencrypt | unix\n")
	}
	if err != nil {
		server.log.Fatalf("%v\n", err)
	}
}

// listenAndServe listens on the TCP network address and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// If srv.Addr is blank, ":http" is used, listenAndServe always returns a non-nil error.
func (server *Server) listenAndServe() error {
	if server.Addr == "" {
		server.Addr = ":http"
	}
	ln, err := server.listen()
	if err != nil {
		return errors.New("[server.listenAndServe()] " + err.Error())
	}

	server.log.Criticalf("\x1b[46m[SYS]\x1b[0m listen and serve HTTP/HTTP2 on %v", server.Addr)

	err = server.Server.Serve(ln)
	if realServeError(err) != nil {
		return errors.New("[server.listenAndServe()] " + err.Error())
	}
	return nil
}

// listenAndServeTLS listens on the TCP network address and
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
// If server.Addr is blank, ":https" is used, listenAndServeTLS always returns a non-nil error.
func (server *Server) listenAndServeTLS() error {
	if server.Addr == "" {
		server.Addr = ":https"
	}

	var cert tls.Certificate
	cert, err := tls.LoadX509KeyPair(server.tlsCertFile, server.tlsKeyFile)
	if err != nil {
		return err
	}

	server.TLSConfig = &tls.Config{
		Certificates:             []tls.Certificate{cert},
		NextProtos:               []string{"http/1.1", "h2"},
		PreferServerCipherSuites: true,
	}

	ln, err := server.listen()
	if err != nil {
		return errors.New("[server.listenAndServeTLS()] " + err.Error())
	}

	server.log.Criticalf("\x1b[46m[SYS]\x1b[0m listen and serve HTTPS(TLS)/HTTP2 on %v", server.Addr)

	err = server.Server.Serve(ln)
	if realServeError(err) != nil {
		return errors.New("[server.listenAndServeTLS()] " + err.Error())
	}
	return nil
}

// listenAndServeLETSENCRYPT listens on a new Automatic TLS using letsencrypt.org service.
// if you want to disable cache file then simple give server.letsencryptFile a value of empty string ""
func (server *Server) listenAndServeLETSENCRYPT() error {
	if server.Addr == "" {
		server.Addr = ":https"
	}

	var m letsencrypt.Manager
	if server.letsencryptFile != "" {
		if err := m.CacheFile(server.letsencryptFile); err != nil {
			return errors.New("[server.listenAndServeLETSENCRYPT()] " + err.Error())
		}
	}

	server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	ln, err := server.listen()
	if err != nil {
		return errors.New("[server.listenAndServeLETSENCRYPT()] " + err.Error())
	}

	server.log.Criticalf("\x1b[46m[SYS]\x1b[0m listen and serve HTTPS(SSL)/HTTP2 on %v", server.Addr)

	err = server.Serve(ln)
	if realServeError(err) != nil {
		return errors.New("[server.listenAndServeLETSENCRYPT()] " + err.Error())
	}
	return nil
}

var (
	errPortAlreadyUsed = errors.New("[server.listenAndServeUNIX()] Port is already used")
	errRemoveUnix      = errors.New("[server.listenAndServeUNIX()] Unexpected error when trying to remove unix socket file. Addr: %s | Trace: %s")
	errChmod           = errors.New("[server.listenAndServeUNIX()] Cannot chmod %#o for %q: %s")
	// errCertKeyMissing  = errors.New("[server.listenAndServeUNIX()] You should provide certFile and keyFile for TLS/SSL")
	// errParseTLS        = errors.New("Couldn't load TLS, certFile=%q, keyFile=%q. Trace: %s")
)

// listenAndServeUNIX announces on the Unix domain socket laddr and listens a Unix service.
func (server *Server) listenAndServeUNIX() error {
	if errOs := os.Remove(server.Addr); errOs != nil && !os.IsNotExist(errOs) {
		return errRemoveUnix.Format(server.Addr, errOs.Error())
	}

	if err := os.Chmod(server.Addr, server.unixFileMode); err != nil {
		return errChmod.Format(server.unixFileMode, server.Addr, err.Error())
	}

	ln, err := server.listen()
	if err != nil {
		return errPortAlreadyUsed.AppendErr(err)
	}

	server.log.Criticalf("\x1b[46m[SYS]\x1b[0m listen and serve HTTP(UNIX)/HTTP2 on %v", server.Addr)

	err = server.Serve(ln)
	if realServeError(err) != nil {
		return errors.New("[server.listenAndServeUNIX()] " + err.Error())
	}
	return nil
}

var grace = new(gracenet.Net)

func (server *Server) listen() (net.Listener, error) {
	ln, err := grace.Listen(server.net, server.Addr)
	if err != nil {
		return nil, err
	}
	ln = tcpKeepAliveListener{ln.(*net.TCPListener)}
	if server.TLSConfig != nil {
		return tls.NewListener(ln, server.TLSConfig), nil
	}
	return ln, nil
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

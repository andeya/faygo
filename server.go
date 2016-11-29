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

	"github.com/henrylee2cn/thinkgo/logging"
	"github.com/henrylee2cn/thinkgo/utils/errors"
	"github.com/rsc/letsencrypt"
	// "github.com/facebookgo/grace/gracehttp"
)

type Server struct {
	nameWithVersion string
	netType         string
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
		err = server.listenAndServe()
	case NETTYPE_TLS:
		err = server.listenAndServeTLS()
	case NETTYPE_LETSENCRYPT:
		err = server.listenAndServeLETSENCRYPT()
	case NETTYPE_UNIX:
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
	server.log.Criticalf("[%s] listen and serve HTTP/HTTP2 on %v", server.nameWithVersion, server.Addr)
	return server.ListenAndServe()
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
	server.log.Criticalf("[%s] listen and serve HTTPS(TLS)/HTTP2 on %v", server.nameWithVersion, server.Addr)
	return server.ListenAndServeTLS(server.tlsCertFile, server.tlsKeyFile)
}

// listenAndServeLETSENCRYPT listens on a new Automatic TLS using letsencrypt.org service.
// if you want to disable cache file then simple give server.letsencryptFile a value of empty string ""
func (server *Server) listenAndServeLETSENCRYPT() error {
	if server.Addr == "" {
		server.Addr = ":https"
	}

	ln, err := net.Listen("tcp4", server.Addr)
	if err != nil {
		return err
	}

	var m letsencrypt.Manager
	if server.letsencryptFile != "" {
		if err = m.CacheFile(server.letsencryptFile); err != nil {
			return err
		}
	}

	tlsConfig := &tls.Config{GetCertificate: m.GetCertificate}
	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, tlsConfig)
	server.log.Criticalf("[%s] listen and serve HTTPS(SSL)/HTTP2 on %v", server.nameWithVersion, server.Addr)

	return server.Serve(tlsListener)
}

var (
	errPortAlreadyUsed = errors.New("Port is already used")
	errRemoveUnix      = errors.New("Unexpected error when trying to remove unix socket file. Addr: %s | Trace: %s")
	errChmod           = errors.New("Cannot chmod %#o for %q: %s")
	errCertKeyMissing  = errors.New("You should provide certFile and keyFile for TLS/SSL")
	errParseTLS        = errors.New("Couldn't load TLS, certFile=%q, keyFile=%q. Trace: %s")
)

// listenAndServeUNIX announces on the Unix domain socket laddr and listens a Unix service.
func (server *Server) listenAndServeUNIX() error {
	if errOs := os.Remove(server.Addr); errOs != nil && !os.IsNotExist(errOs) {
		return errRemoveUnix.Format(server.Addr, errOs.Error())
	}

	ln, err := net.Listen("unix", server.Addr)
	if err != nil {
		return errPortAlreadyUsed.AppendErr(err)
	}

	if err = os.Chmod(server.Addr, server.unixFileMode); err != nil {
		return errChmod.Format(server.unixFileMode, server.Addr, err.Error())
	}
	server.log.Criticalf("[%s] listen and serve HTTP(UNIX)/HTTP2 on %v", server.nameWithVersion, server.Addr)

	return server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
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

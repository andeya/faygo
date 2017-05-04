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
	"fmt"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

type (
	// GlobalConfig is global configuration
	GlobalConfig struct {
		Cache   CacheConfig `ini:"cache"`
		Gzip    GzipConfig  `ini:"gzip"`
		Log     LogConfig   `ini:"log"`
		warnMsg string      `int:"-"`
	}
	// Config is the configuration information for each web instance
	Config struct {
		// RunMode         string      `ini:"run_mode"`         // run mode: dev | prod
		NetTypes       []string    `ini:"net_types" delim:"|"` // List of network type: http | https | unix_http | unix_https | letsencrypt | unix_letsencrypt
		Addrs          []string    `ini:"addrs" delim:"|"`     // List of multiple listening addresses
		TLSCertFile    string      `ini:"tls_certfile"`        // TLS certificate file path
		TLSKeyFile     string      `ini:"tls_keyfile"`         // TLS key file path
		LetsencryptDir string      `ini:"letsencrypt_dir"`     // Let's Encrypt TLS certificate cache directory
		UNIXFileMode   os.FileMode `ini:"unix_filemode"`       // File permissions for UNIX Server (438 equivalent to 0666)
		// Maximum duration for reading the full request (including body).
		//
		// This also limits the maximum duration for idle keep-alive
		// connections.
		//
		// By default request read timeout is unlimited.
		ReadTimeout time.Duration `ini:"read_timeout"`
		// Maximum duration for writing the full response (including body).
		//
		// By default response write timeout is unlimited.
		WriteTimeout time.Duration `ini:"write_timeout"`
		// Maximum size of memory that can be used when receiving uploaded files.
		MultipartMaxMemoryMB int64         `ini:"multipart_maxmemory_mb"`
		multipartMaxMemory   int64         `ini:"-"`
		Router               RouterConfig  `ini:"router"`
		XSRF                 XSRFConfig    `ini:"xsrf"`
		Session              SessionConfig `ini:"session"`
		APIdoc               APIdocConfig  `ini:"apidoc"`
	}
	// RouterConfig is the configuration about router
	RouterConfig struct {
		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// For example if /foo/ is requested but a route only exists for /foo, the
		// client is redirected to /foo with http status code 301 for GET requests
		// and 307 for all other request methods.
		RedirectTrailingSlash bool `ini:"redirect_trailing_slash"`
		// If enabled, the router tries to fix the current request path, if no
		// handle is registered for it.
		// First superfluous path elements like ../ or // are removed.
		// Afterwards the router does a case-insensitive lookup of the cleaned path.
		// If a handle can be found for this route, the router makes a redirection
		// to the corrected path with status code 301 for GET requests and 307 for
		// all other request methods.
		// For example /FOO and /..//Foo could be redirected to /foo.
		// RedirectTrailingSlash is independent of this option.
		RedirectFixedPath bool `ini:"redirect_fixed_path"`
		// If enabled, the router checks if another method is allowed for the
		// current route, if the current request can not be routed.
		// If this is the case, the request is answered with 'Method Not Allowed'
		// and HTTP status code 405.
		// If no other Method is allowed, the request is delegated to the NotFound
		// handler.
		HandleMethodNotAllowed bool `ini:"handle_method_not_allowed"`
		// If enabled, the router automatically replies to OPTIONS requests.
		// Custom OPTIONS handlers take priority over automatic replies.
		HandleOPTIONS bool `ini:"handle_options"`
	}
	// GzipConfig is the configuration about gzip
	GzipConfig struct {
		// if EnableGzip, compress response content.
		Enable bool `ini:"enable"`
		//Content will only be compressed if content length is either unknown or greater than gzipMinLength.
		//Default size==20B same as nginx
		MinLength int `ini:"min_length"`
		//The compression level used for deflate compression. (0-9).
		//Non-file response Body's compression level is 0-9, but the files' always 9
		CompressLevel int `ini:"compress_level"`
		//List of HTTP methods to compress. If not set, only GET requests are compressed.
		Methods []string `ini:"methods" delim:"|"`
		// StaticExtensionsToGzip []string
	}
	// CacheConfig is the configuration about cache
	CacheConfig struct {
		// Whether to enable caching static files
		Enable bool `ini:"enable"`
		// Max size by MB for file cache.
		// The cache size will be set to 512KB at minimum.
		// If the size is set relatively large, you should call
		// `debug.SetGCPercent()`, set it to a much smaller value
		// to limit the memory consumption and GC pause time.
		SizeMB int64 `ini:"size_mb"`
		// expire in xxx seconds for file cache.
		// Expire <= 0 (second) means no expire, but it can be evicted when cache is full.
		Expire int `ini:"expire"`
	}
	// XSRFConfig is the configuration about XSRF filter
	XSRFConfig struct {
		Enable bool   `ini:"enable"`
		Key    string `ini:"key"`
		Expire int    `ini:"expire"`
	}
	// SessionConfig is the configuration about session
	SessionConfig struct {
		Enable                bool   `ini:"enable"`                 // Whether enabled or not
		Provider              string `ini:"provider"`               // Data storage
		Name                  string `ini:"name"`                   // The client stores the name of the cookie
		ProviderConfig        string `ini:"provider_config"`        // According to the different engine settings different configuration information
		CookieLifetime        int    `ini:"cookie_lifetime"`        // The default value is 0, which is the lifetime of the browser
		GCLifetime            int64  `ini:"gc_lifetime"`            // The interval between triggering the GC
		MaxLifetime           int64  `ini:"max_lifetime"`           // The session max lefetime
		AutoSetCookie         bool   `ini:"auto_setcookie"`         // Automatically set on the session cookie value, the general default true
		Domain                string `ini:"domain"`                 // The domain name that is allowed to access this cookie
		EnableSidInHttpHeader bool   `ini:"enable_sid_in_header"`   // Whether to write a session ID to the header
		NameInHttpHeader      string `ini:"name_in_header"`         // The name of the header when the session ID is written to the header
		EnableSidInUrlQuery   bool   `ini:"enable_sid_in_urlquery"` // Whether to write the session ID to the URL Query params
	}
	// LogConfig is the configuration about log
	LogConfig struct {
		ConsoleEnable bool   `ini:"console_enable"`
		ConsoleLevel  string `ini:"console_level"` // critical | error | warning | notice | info | debug
		FileEnable    bool   `ini:"file_enable"`
		FileLevel     string `ini:"file_level"` // critical | error | warning | notice | info | debug
		AsyncLen      int    `ini:"async_len"`  // the length of asynchronous buffer, 0 means synchronization
	}
	// APIdocConfig is the configuration about API doc
	APIdocConfig struct {
		Enable     bool     `ini:"enable"`              // whether to enable API doc
		Path       string   `ini:"path"`                // API doc url path
		NoLimit    bool     `ini:"nolimit"`             // if true, access is not restricted
		RealIP     bool     `ini:"real_ip"`             // if true, means verifying the real IP of the visitor
		Whitelist  []string `ini:"whitelist" delim:"|"` // `whitelist=192.*|202.122.246.170` means: only IP addresses that are prefixed with `192 'or equal to` 202.122.246.170' are allowed
		Desc       string   `ini:"desc"`                // description of the application
		Email      string   `ini:"email"`               // technician's Email
		TermsURL   string   `ini:"terms_url"`           // terms of service
		License    string   `ini:"license"`             // the license used by the API
		LicenseURL string   `ini:"license_url"`         // the URL of the protocol content page
	}
)

// some default config
const (
	// RUNMODE_DEV                 = "dev"
	// RUNMODE_PROD                = "prod"
	MB                          = 1 << 20 // 1MB
	defaultMultipartMaxMemory   = 32 * MB // 32 MB
	defaultMultipartMaxMemoryMB = 32
	defaultPort                 = 8080

	// The path for the configuration files
	CONFIG_DIR = "./config/"
	// global config file name
	GLOBAL_CONFIG_FILE = "__global___.ini"
)

var (
	appCount uint32
)

// global config
var globalConfig = func() GlobalConfig {
	var background = GlobalConfig{
		Cache: CacheConfig{
			Enable: false,
			SizeMB: 32,
			Expire: 60,
		},
		Gzip: GzipConfig{
			Enable:        false,
			MinLength:     20,
			CompressLevel: 1,
			Methods:       []string{"GET"},
		},
		Log: LogConfig{
			ConsoleEnable: true,
			ConsoleLevel:  "debug",
			FileEnable:    false,
			FileLevel:     "debug",
		},
	}
	filename := CONFIG_DIR + GLOBAL_CONFIG_FILE

	err := SyncINI(
		&background,
		func(existed bool, saveOnce func() error) error {
			if !(background.Log.ConsoleEnable || background.Log.FileEnable) {
				background.Log.ConsoleEnable = true
				background.warnMsg = "config: log::enable_console and log::enable_file can not be disabled at the same time, so automatically open console log."
			}
			return nil
		},
		filename,
	)

	if err != nil {
		panic(err)
	}

	return background
}()

func newConfig(filename string) Config {
	var addr string

	addr = fmt.Sprintf("0.0.0.0:%d", defaultPort+atomic.LoadUint32(&appCount))
	atomic.AddUint32(&appCount, 1)
	var background = Config{
		// RunMode:              RUNMODE_DEV,
		NetTypes:             []string{NETTYPE_HTTP},
		Addrs:                []string{addr},
		UNIXFileMode:         0666,
		MultipartMaxMemoryMB: defaultMultipartMaxMemoryMB,
		Router: RouterConfig{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
		},
		XSRF: XSRFConfig{
			Enable: false,
			Key:    "faygoxsrf",
			Expire: 3600,
		},
		Session: SessionConfig{
			Enable:                false,
			Provider:              "memory",
			Name:                  "faygosessionID",
			CookieLifetime:        0, //set cookie default is the browser life
			GCLifetime:            300,
			MaxLifetime:           3600,
			ProviderConfig:        "",
			AutoSetCookie:         true,
			Domain:                "",
			EnableSidInHttpHeader: false, //	enable store/get the sessionId into/from http headers
			NameInHttpHeader:      "Faygosessionid",
			EnableSidInUrlQuery:   false, //	enable get the sessionId from Url Query params
		},
		APIdoc: APIdocConfig{
			Enable:  true,
			Path:    "/apidoc/",
			NoLimit: false,
			RealIP:  false,
			Whitelist: []string{
				"127.*",
				"192.168.*",
			},
		},
	}

	err := SyncINI(
		&background,
		func(existed bool, saveOnce func() error) error {
			// switch background.RunMode {
			// case RUNMODE_DEV, RUNMODE_PROD:
			// default:
			// 	panic("Please set a valid config item run_mode, refer to the following:\ndev | prod")
			// }
			if len(background.NetTypes) != len(background.Addrs) {
				panic("The number of configuration items `net_types` and `addrs` must be equal")
			}
			if len(background.NetTypes) == 0 {
				panic("The number of configuration items `net_types` and `addrs` must be greater than zero")
			}
			for _, t := range background.NetTypes {
				switch t {
				case NETTYPE_HTTP, NETTYPE_UNIX_HTTP, NETTYPE_HTTPS, NETTYPE_UNIX_HTTPS, NETTYPE_LETSENCRYPT, NETTYPE_UNIX_LETSENCRYPT:
				default:
					panic("Please set a valid config item `net_types`, refer to the following:" + __netTypes__ + "\n")
				}
			}
			background.multipartMaxMemory = background.MultipartMaxMemoryMB * MB
			background.APIdoc.Comb()
			return nil
		},
		filename,
	)
	if err != nil {
		panic(err)
	}

	return background
}

// Comb combs APIdoc config
func (conf *APIdocConfig) Comb() {
	ipPrefixMap := map[string]bool{}
	for _, ipPrefix := range conf.Whitelist {
		if len(ipPrefix) > 0 {
			ipPrefixMap[ipPrefix] = true
		}
	}
	conf.Whitelist = conf.Whitelist[:0]
	for ipPrefix := range ipPrefixMap {
		conf.Whitelist = append(conf.Whitelist, ipPrefix)
	}
	sort.Strings(conf.Whitelist)
	conf.Path = "/" + strings.Trim(conf.Path, "/") + "/"
}

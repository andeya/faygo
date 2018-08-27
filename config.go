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
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	// GlobalConfig is global config
	GlobalConfig struct {
		Cache   CacheConfig `ini:"cache" comment:"Cache section"`
		Gzip    GzipConfig  `ini:"gzip" comment:"Gzip section"`
		Log     LogConfig   `ini:"log" comment:"Log section"`
		warnMsg string      `int:"-"`
	}
	// Config is the config information for each web instance
	Config struct {
		// RunMode         string      `ini:"run_mode" comment:"run mode: dev|prod"`
		NetTypes          []string    `ini:"net_types" delim:"|" comment:"List of network type: http|https|unix_http|unix_https|letsencrypt|unix_letsencrypt"`
		Addrs             []string    `ini:"addrs" delim:"|" comment:"List of multiple listening addresses"`
		TLSCertFile       string      `ini:"tls_certfile" comment:"TLS certificate file path"`
		TLSKeyFile        string      `ini:"tls_keyfile" comment:"TLS key file path"`
		LetsencryptDir    string      `ini:"letsencrypt_dir" comment:"Let's Encrypt TLS certificate cache directory"`
		UNIXFileMode      string      `ini:"unix_filemode" comment:"File permissions for UNIX listener, requires octal number"`
		unixFileMode      os.FileMode `ini:"-"`
		HttpRedirectHttps bool        `ini:"http_redirect_https" comment:"Redirect from 'http://hostname:port1' to 'https://hostname:port2'"`
		// Maximum duration for reading the full request (including body).
		//
		// This also limits the maximum duration for idle keep-alive
		// connections.
		//
		// By default request read timeout is unlimited.
		ReadTimeout time.Duration `ini:"read_timeout" comment:"Maximum duration for reading the full request (including body); ns|µs|ms|s|m|h"`
		// Maximum duration for writing the full response (including body).
		//
		// By default response write timeout is unlimited.
		WriteTimeout          time.Duration `ini:"write_timeout" comment:"Maximum duration for writing the full response (including body); ns|µs|ms|s|m|h"`
		MultipartMaxMemoryMB  int64         `ini:"multipart_maxmemory_mb" comment:"Maximum size of memory that can be used when receiving uploaded files"`
		multipartMaxMemory    int64         `ini:"-"`
		Router                RouterConfig  `ini:"router" comment:"Routing config section"`
		XSRF                  XSRFConfig    `ini:"xsrf" comment:"XSRF security section"`
		Session               SessionConfig `ini:"session" comment:"Session section"`
		SlowResponseThreshold time.Duration `ini:"slow_response_threshold" comment:"When response time > slow_response_threshold, log level = 'WARNING'; 0 means not limited; ns|µs|ms|s|m|h"`
		slowResponseThreshold time.Duration `ini:"-"`
		PrintBody             bool          `ini:"print_body" comment:"Form requests are printed in JSON format, but other types are printed as-is"`
		APIdoc                APIdocConfig  `ini:"apidoc" comment:"API documentation section"`
	}
	// RouterConfig is the config about router
	RouterConfig struct {
		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// For example if /foo/ is requested but a route only exists for /foo, the
		// client is redirected to /foo with http status code 301 for GET requests
		// and 307 for all other request methods.
		RedirectTrailingSlash bool `ini:"redirect_trailing_slash" comment:"Automatic redirection (for example, '/foo/' -> '/foo')"`
		// If enabled, the router tries to fix the current request path, if no
		// handle is registered for it.
		// First superfluous path elements like ../ or // are removed.
		// Afterwards the router does a case-insensitive lookup of the cleaned path.
		// If a handle can be found for this route, the router makes a redirection
		// to the corrected path with status code 301 for GET requests and 307 for
		// all other request methods.
		// For example /FOO and /..//Foo could be redirected to /foo.
		// RedirectTrailingSlash is independent of this option.
		RedirectFixedPath bool `ini:"redirect_fixed_path" comment:"Tries to fix the current request path, if no handle is registered for it"`
		// If enabled, the router checks if another method is allowed for the
		// current route, if the current request can not be routed.
		// If this is the case, the request is answered with 'Method Not Allowed'
		// and HTTP status code 405.
		// If no other Method is allowed, the request is delegated to the NotFound
		// handler.
		HandleMethodNotAllowed bool `ini:"handle_method_not_allowed" comment:"Returns 405 if the requested method does not exist, otherwise returns 404"`
		// If enabled, the router automatically replies to OPTIONS requests.
		// Custom OPTIONS handlers take priority over automatic replies.
		HandleOPTIONS   bool `ini:"handle_options" comment:"Automatic response OPTIONS request, you can set the default Handler in faygo"`
		NoDefaultParams bool `ini:"no_default_params" comment:"If true, don't assign default request parameter values based on initial parameter values of the routing handler"`
		DefaultUpload   bool `ini:"default_upload" comment:"Automatically register the default router: /upload/*filepath"`
		DefaultStatic   bool `ini:"default_static" comment:"Automatically register the default router: /static/*filepath"`
	}
	// GzipConfig is the config about gzip
	GzipConfig struct {
		// if EnableGzip, compress response content.
		Enable bool `ini:"enable" comment:"Whether enabled or not"`
		//Content will only be compressed if content length is either unknown or greater than gzipMinLength.
		//Default size==20B same as nginx
		MinLength int `ini:"min_length" comment:"The minimum length of content to be compressed"`
		//The compression level used for deflate compression. (0-9).
		//Non-file response Body's compression level is 0-9, but the files' always 9
		CompressLevel int `ini:"compress_level" comment:"Non-file response Body's compression level is 0-9, but the files' always 9"`
		//List of HTTP methods to compress. If not set, only GET requests are compressed.
		Methods []string `ini:"methods" delim:"|" comment:"List of HTTP methods to compress. If not set, only GET requests are compressed."`
		// StaticExtensionsToGzip []string
	}
	// CacheConfig is the config about cache
	CacheConfig struct {
		// Whether to enable caching static files
		Enable bool `ini:"enable" comment:"Whether enabled or not"`
		// Max size by MB for file cache.
		// The cache size will be set to 512KB at minimum.
		// If the size is set relatively large, you should call
		// `debug.SetGCPercent()`, set it to a much smaller value
		// to limit the memory consumption and GC pause time.
		SizeMB int64 `ini:"size_mb" comment:"Max size by MB for file cache, the cache size will be set to 512KB at minimum."`
		// expire in xxx seconds for file cache.
		// ExpireSecond <= 0 (second) means no expire, but it can be evicted when cache is full.
		ExpireSecond int `ini:"expire_second" comment:"Maximum duration for caching"`
	}
	// XSRFConfig is the config about XSRF filter
	XSRFConfig struct {
		Enable       bool   `ini:"enable" comment:"Whether enabled or not"`
		Key          string `ini:"key" comment:"Encryption key"`
		ExpireSecond int    `ini:"expire_second" comment:"Expire of XSRF token"`
	}
	// SessionConfig is the config about session
	SessionConfig struct {
		Enable                bool   `ini:"enable" comment:"Whether enabled or not"`
		Provider              string `ini:"provider" comment:"Data storage"`
		Name                  string `ini:"name" comment:"The client stores the name of the cookie"`
		ProviderConfig        string `ini:"provider_config" comment:"According to the different engine settings different config information"`
		CookieLifeSecond      int    `ini:"cookie_life_second" comment:"The default value is 0, which is the lifetime of the browser"`
		GcLifeSecond          int64  `ini:"gc_life_second" comment:"The interval between triggering the GC"`
		MaxLifeSecond         int64  `ini:"max_life_second" comment:"The session max lefetime"`
		AutoSetCookie         bool   `ini:"auto_setcookie" comment:"Automatically set on the session cookie value, the general default true"`
		Domain                string `ini:"domain" comment:"The domain name that is allowed to access this cookie"`
		EnableSidInHttpHeader bool   `ini:"enable_sid_in_header" comment:"Whether to write a session ID to the header"`
		NameInHttpHeader      string `ini:"name_in_header" comment:"The name of the header when the session ID is written to the header"`
		EnableSidInUrlQuery   bool   `ini:"enable_sid_in_urlquery" comment:"Whether to write the session ID to the URL Query params"`
	}
	// LogConfig is the config about log
	LogConfig struct {
		ConsoleEnable bool   `ini:"console_enable" comment:"Whether enabled or not console logger"`
		ConsoleLevel  string `ini:"console_level" comment:"Console logger level: critical|error|warning|notice|info|debug"`
		FileEnable    bool   `ini:"file_enable" comment:"Whether enabled or not file logger"`
		FileLevel     string `ini:"file_level" comment:"File logger level: critical|error|warning|notice|info|debug"`
		AsyncLen      int    `ini:"async_len" comment:"The length of asynchronous buffer, 0 means synchronization"`
	}
	// APIdocConfig is the config about API doc
	APIdocConfig struct {
		Enable     bool     `ini:"enable" comment:"Whether enabled or not"`
		Path       string   `ini:"path" comment:"The URL path"`
		NoLimit    bool     `ini:"nolimit" comment:"If true, access is not restricted"`
		RealIP     bool     `ini:"real_ip" comment:"if true, means verifying the real IP of the visitor"`
		Whitelist  []string `ini:"whitelist" delim:"|" comment:"'whitelist=192.*|202.122.246.170' means: only IP addresses that are prefixed with '192.' or equal to '202.122.246.170' are allowed"`
		Desc       string   `ini:"desc" comment:"Description of the application"`
		Email      string   `ini:"email" comment:"Technician's Email"`
		TermsURL   string   `ini:"terms_url" comment:"Terms of service"`
		License    string   `ini:"license" comment:"The license used by the API"`
		LicenseURL string   `ini:"license_url" comment:"The URL of the protocol content page"`
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
)

var (
	// configDir the config files directory
	configDir = "./config/"
	// globalConfigFile global config file name
	globalConfigFile = "__global___.ini"
)

// ConfigDir returns the config files directory
func ConfigDir() string {
	return configDir
}

// global config
var globalConfig = func() GlobalConfig {
	// get config dir
	flag.StringVar(&configDir, "cfg_dir", configDir, "Configuration files directory")
	flag.Parse()

	var background = &GlobalConfig{
		Cache: CacheConfig{
			Enable:       false,
			SizeMB:       32,
			ExpireSecond: 60,
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
	filename := configDir + globalConfigFile

	err := SyncINI(
		background,
		func(onceUpdateFunc func() error) error {
			if !(background.Log.ConsoleEnable || background.Log.FileEnable) {
				background.Log.ConsoleEnable = true
				background.warnMsg = "config: log::enable_console and log::enable_file can not be disabled at the same time, so automatically open console log."
			}
			return onceUpdateFunc()
		},
		filename,
	)

	if err != nil {
		panic(err)
	}

	return *background
}()

// NewDefaultConfig creates a new default framework config.
func NewDefaultConfig() *Config {
	return &Config{
		// RunMode:              RUNMODE_DEV,
		NetTypes:             []string{NETTYPE_HTTP},
		Addrs:                []string{fmt.Sprintf("0.0.0.0:%d", defaultPort+len(AllFrames()))},
		UNIXFileMode:         "0666",
		MultipartMaxMemoryMB: defaultMultipartMaxMemoryMB,
		Router: RouterConfig{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
			DefaultUpload:          true,
			DefaultStatic:          true,
		},
		XSRF: XSRFConfig{
			Enable:       false,
			Key:          "faygoxsrf",
			ExpireSecond: 3600,
		},
		Session: SessionConfig{
			Enable:                false,
			Provider:              "memory",
			Name:                  "faygosessionID",
			CookieLifeSecond:      0, //set cookie default is the browser life
			GcLifeSecond:          300,
			MaxLifeSecond:         3600,
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
}

func (c *Config) check() {
	// switch c.RunMode {
	// case RUNMODE_DEV, RUNMODE_PROD:
	// default:
	// 	panic("Please set a valid config item run_mode, refer to the following:\ndev|prod")
	// }
	if len(c.NetTypes) != len(c.Addrs) {
		panic("The number of config items `net_types` and `addrs` must be equal")
	}
	if len(c.NetTypes) == 0 {
		panic("The number of config items `net_types` and `addrs` must be greater than zero")
	}
	for _, t := range c.NetTypes {
		switch t {
		case NETTYPE_HTTP, NETTYPE_UNIX_HTTP, NETTYPE_HTTPS, NETTYPE_UNIX_HTTPS, NETTYPE_LETSENCRYPT, NETTYPE_UNIX_LETSENCRYPT:
		default:
			panic("Please set a valid config item `net_types`, refer to the following:" + __netTypes__)
		}
	}
	fileMode, err := strconv.ParseUint(c.UNIXFileMode, 8, 32)
	if err != nil {
		panic("The config item `unix_filemode` is not a valid octal number:" + c.UNIXFileMode)
	}
	c.unixFileMode = os.FileMode(fileMode)
	c.UNIXFileMode = fmt.Sprintf("%#o", fileMode)
	c.multipartMaxMemory = c.MultipartMaxMemoryMB * MB
	if c.SlowResponseThreshold <= 0 {
		c.slowResponseThreshold = time.Duration(math.MaxInt64)
	} else {
		c.slowResponseThreshold = c.SlowResponseThreshold
	}
	c.APIdoc.Comb()
}

func newConfigFromFileAndCheck(filename string) *Config {
	var background = NewDefaultConfig()
	err := SyncINI(
		background,
		func(onceUpdateFunc func() error) error {
			background.check()
			return onceUpdateFunc()
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

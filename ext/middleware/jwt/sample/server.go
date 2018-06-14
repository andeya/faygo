package main

import (
	"time"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/middleware/jwt"
)

func helloHandler(c *faygo.Context) error {
	claims := jwt.ExtractClaims(c)
	return c.JSON(200, map[string]interface{}{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

func main() {
	r := faygo.New("jwt-test")
	// the jwt middleware
	authMiddleware := &jwt.FaygoJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Minute,
		MaxRefresh: time.Minute * 3,
		//登录响应
		LoginResponse: func(c *faygo.Context, code int, token string, expire time.Time) error {
			faygo.Debug("LoginResponse")
			return c.JSON(code, faygo.Map{
				"code":   code,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
				"custom": "custom info ",
			})
		},
		/*PayloadFunc: func(data interface{}) map[string]interface{} {
			return faygo.Map{
				"custom2": "custom info2 ",
			}
		},*/
		//认证
		Authenticator: func(userId string, password string, c *faygo.Context) (interface{}, bool) {
			faygo.Debug("Authenticator认证")
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return userId, true
			}

			return userId, false
		},

		//授权
		Authorizator: func(userId interface{}, c *faygo.Context) bool {
			faygo.Debug("Authorizator 授权")
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *faygo.Context, code int, message string) {
			faygo.Debug("Unauthorized 未授权")
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:token",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer" 持票人
		TokenHeadName: "",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
	//调用初始设置函数，必须的。
	err := authMiddleware.MiddlewareInit()
	if err != nil {
		faygo.Error(err)
	}

	r.POST("/login", faygo.HandlerFunc(authMiddleware.LoginHandler))
	auth := r.Group("/auth")
	//auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", faygo.HandlerFunc(helloHandler)).Use(authMiddleware.MiddlewareFunc())
		auth.GET("/refreshtoken", faygo.HandlerFunc(authMiddleware.RefreshHandler))
	}

	r.Run()
}

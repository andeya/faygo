package jwt_test

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

func Example() {
	r := faygo.New("jwt-test")
	// the jwt middleware
	authMiddleware := &jwt.FaygoJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *faygo.Context) (string, bool) {
			if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *faygo.Context) bool {
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *faygo.Context, code int, message string) {
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
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	r.POST("/login", faygo.HandlerFunc(authMiddleware.LoginHandler))

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", faygo.HandlerFunc(helloHandler))
		auth.GET("/refresh_token", faygo.HandlerFunc(authMiddleware.RefreshHandler))
	}

	r.Run()
}

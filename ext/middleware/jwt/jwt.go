package jwt

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"

	"github.com/henrylee2cn/faygo"
)

// FaygoJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
//
// Reference https://github.com/appleboy/gin-jwt
type FaygoJWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on userID and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	Authenticator func(userID string, password string, c *faygo.Context) (string, bool)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userID string, c *faygo.Context) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(userID string) map[string]interface{}

	// User can define own Unauthorized func.
	Unauthorized func(*faygo.Context, int, string)

	// Set the identity handler function
	IdentityHandler func(jwt.MapClaims) string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time
}

// Login form structure.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// MiddlewareInit initialize jwt configs.
func (mw *FaygoJWTMiddleware) MiddlewareInit() error {

	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(userID string, c *faygo.Context) bool {
			return true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *faygo.Context, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(claims jwt.MapClaims) string {
			return claims["id"].(string)
		}
	}

	if mw.Realm == "" {
		return errors.New("realm is required")
	}

	if mw.Key == nil {
		return errors.New("secret key is required")
	}

	return nil
}

// MiddlewareFunc makes FaygoJWTMiddleware implement the Middleware interface.
func (mw *FaygoJWTMiddleware) MiddlewareFunc() faygo.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *faygo.Context) error {
			mw.unauthorized(c, http.StatusInternalServerError, err.Error())
			return nil
		}
	} else {
		return func(c *faygo.Context) error {
			mw.middlewareImpl(c)
			return nil
		}
	}
}

func (mw *FaygoJWTMiddleware) middlewareImpl(c *faygo.Context) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	id := mw.IdentityHandler(claims)
	c.SetData("JWT_PAYLOAD", claims)
	c.SetData("userID", id)

	if !mw.Authorizator(id, c) {
		mw.unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
		return
	}

	c.Next()
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *FaygoJWTMiddleware) LoginHandler(c *faygo.Context) error {

	// Initial middleware default setting.
	mw.MiddlewareInit()

	var loginVals Login

	if c.BindJSON(&loginVals) != nil {
		return mw.unauthorized(c, http.StatusBadRequest, "Missing Username or Password")
	}

	if mw.Authenticator == nil {
		return mw.unauthorized(c, http.StatusInternalServerError, "Missing define authenticator func")
	}

	userID, ok := mw.Authenticator(loginVals.Username, loginVals.Password, c)

	if !ok {
		return mw.unauthorized(c, http.StatusUnauthorized, "Incorrect Username / Password")
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(loginVals.Username) {
			claims[key] = value
		}
	}

	if userID == "" {
		userID = loginVals.Username
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["id"] = userID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	tokenString, err := token.SignedString(mw.Key)

	if err != nil {
		return mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")

	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the FaygoJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *FaygoJWTMiddleware) RefreshHandler(c *faygo.Context) error {
	token, _ := mw.parseToken(c)
	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		return mw.unauthorized(c, http.StatusUnauthorized, "Token is expired.")
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	newClaims["id"] = claims["id"]
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = origIat

	tokenString, err := newToken.SignedString(mw.Key)

	if err != nil {
		return mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *faygo.Context) jwt.MapClaims {

	if exists := c.HasData("JWT_PAYLOAD"); !exists {
		emptyClaims := make(jwt.MapClaims)
		return emptyClaims
	}

	jwtClaims := c.Data("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}

// TokenGenerator handler that clients can use to get a jwt token.
func (mw *FaygoJWTMiddleware) TokenGenerator(userID string) string {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			claims[key] = value
		}
	}

	claims["id"] = userID
	claims["exp"] = mw.TimeFunc().Add(mw.Timeout).Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	tokenString, _ := token.SignedString(mw.Key)

	return tokenString
}

func (mw *FaygoJWTMiddleware) jwtFromHeader(c *faygo.Context, key string) (string, error) {
	authHeader := c.HeaderParam(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

func (mw *FaygoJWTMiddleware) jwtFromQuery(c *faygo.Context, key string) (string, error) {
	token := c.QueryParam(key)

	if token == "" {
		return "", errors.New("Query token empty")
	}

	return token, nil
}

func (mw *FaygoJWTMiddleware) jwtFromCookie(c *faygo.Context, key string) (string, error) {
	cookie := c.CookieParam(key)

	if cookie == "" {
		return "", errors.New("Cookie token empty")
	}

	return cookie, nil
}

func (mw *FaygoJWTMiddleware) parseToken(c *faygo.Context) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split(mw.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = mw.jwtFromHeader(c, parts[1])
	case "query":
		token, err = mw.jwtFromQuery(c, parts[1])
	case "cookie":
		token, err = mw.jwtFromCookie(c, parts[1])
	}

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		return mw.Key, nil
	})
}

func (mw *FaygoJWTMiddleware) unauthorized(c *faygo.Context, code int, message string) error {

	if mw.Realm == "" {
		mw.Realm = "faygo jwt"
	}

	c.SetHeader("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Stop()

	mw.Unauthorized(c, code, message)

	return nil
}

package jwt

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/henrylee2cn/faygo"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// FaygoJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
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
	// Option return user data, if so, user data will be stored in Claim Array.
	Authenticator func(userID string, password string, c *faygo.Context) (interface{}, bool)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(data interface{}, c *faygo.Context) bool

	// Callback function that will be called during Login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(data interface{}) map[string]interface{}

	// User can define own Unauthorized func.
	Unauthorized func(*faygo.Context, int, string)

	// User can define own LoginResponse func.
	LoginResponse func(*faygo.Context, int, string, time.Time) error

	// User can define own RefreshResponse func.
	RefreshResponse func(*faygo.Context, int, string, time.Time) error

	// Set the identity handler function
	IdentityHandler func(jwt.MapClaims) interface{}

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

	// HTTP Status messages for when something in the JWT middleware fails.
	// Check error (e) to determine the appropriate error message.
	HTTPStatusMessageFunc func(e error, c *faygo.Context) string

	// Private key file for asymmetric algorithms
	PrivKeyFile string

	// Public key file for asymmetric algorithms
	PubKeyFile string

	// Private key
	privKey *rsa.PrivateKey

	// Public key
	pubKey *rsa.PublicKey
}

var (
	// ErrMissingRealm indicates Realm name is required
	ErrMissingRealm = errors.New("realm is missing")

	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("secret key is required")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrMissingAuthenticatorFunc indicates Authenticator is required
	ErrMissingAuthenticatorFunc = errors.New("FaygoJWTMiddleware.Authenticator func is undefined")

	// ErrMissingLoginValues indicates a user tried to authenticate without username or password
	ErrMissingLoginValues = errors.New("missing Usercode or Password")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect Usercode or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("token is expired")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("public key invalid")
)

// Login form structure.
type Login struct {
	Usercode string `form:"usercode" json:"usercode" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (mw *FaygoJWTMiddleware) readKeys() error {
	err := mw.privateKey()
	if err != nil {
		return err
	}
	err = mw.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (mw *FaygoJWTMiddleware) privateKey() error {
	keyData, err := ioutil.ReadFile(mw.PrivKeyFile)
	if err != nil {
		return ErrNoPrivKeyFile
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPrivKey
	}
	mw.privKey = key
	return nil
}

func (mw *FaygoJWTMiddleware) publicKey() error {
	keyData, err := ioutil.ReadFile(mw.PubKeyFile)
	if err != nil {
		return ErrNoPubKeyFile
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPubKey
	}
	mw.pubKey = key
	return nil
}

func (mw *FaygoJWTMiddleware) usingPublicKeyAlgo() bool {
	switch mw.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
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
		mw.Authorizator = func(data interface{}, c *faygo.Context) bool {
			return true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *faygo.Context, code int, message string) {
			c.JSON(code, faygo.Map{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.LoginResponse == nil {
		mw.LoginResponse = func(c *faygo.Context, code int, token string, expire time.Time) error {
			return c.JSON(http.StatusOK, faygo.Map{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.RefreshResponse == nil {
		mw.RefreshResponse = func(c *faygo.Context, code int, token string, expire time.Time) error {
			return c.JSON(http.StatusOK, faygo.Map{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(claims jwt.MapClaims) interface{} {
			return claims["id"]
		}
	}

	if mw.HTTPStatusMessageFunc == nil {
		mw.HTTPStatusMessageFunc = func(e error, c *faygo.Context) string {
			return e.Error()
		}
	}

	if mw.Realm == "" {
		return ErrMissingRealm
	}

	if mw.usingPublicKeyAlgo() {
		return mw.readKeys()
	}

	if mw.Key == nil {
		return ErrMissingSecretKey
	}
	return nil
}

// MiddlewareFunc makes FaygoJWTMiddleware implement the Middleware interface.
func (mw *FaygoJWTMiddleware) MiddlewareFunc() faygo.HandlerFunc {
	//初始设置改为手动调用
	/*if err := mw.MiddlewareInit(); err != nil {
		return func(c *faygo.Context) error {
			mw.unauthorized(c, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(err, nil))
			return nil
		}
	}*/

	return func(c *faygo.Context) error {
		mw.middlewareImpl(c)
		return nil
	}
}

func (mw *FaygoJWTMiddleware) middlewareImpl(c *faygo.Context) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
		return
	}
	claims := token.Claims.(jwt.MapClaims)

	id := mw.IdentityHandler(claims)
	c.SetData("JWT_PAYLOAD", claims)
	c.SetData("userID", id)

	if !mw.Authorizator(id, c) {
		mw.unauthorized(c, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, c))
		return
	}

	c.Next()
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *FaygoJWTMiddleware) LoginHandler(c *faygo.Context) error {

	// Initial middleware default setting. 初始设置改为手动调用
	/*if err := mw.MiddlewareInit(); err != nil {
		return mw.unauthorized(c, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(err, c))

	}*/

	var loginVals Login

	if c.BindJSON(&loginVals) != nil {
		return mw.unauthorized(c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrMissingLoginValues, c))
	}

	if mw.Authenticator == nil {
		return mw.unauthorized(c, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(ErrMissingAuthenticatorFunc, c))

	}

	data, ok := mw.Authenticator(loginVals.Usercode, loginVals.Password, c)

	if !ok {
		return mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedAuthentication, c))

	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}
	if claims["id"] == nil {
		claims["id"] = loginVals.Usercode
	}
	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)

	if err != nil {
		return mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, c))

	}

	return mw.LoginResponse(c, http.StatusOK, tokenString, expire)
}

func (mw *FaygoJWTMiddleware) signedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	if mw.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(mw.privKey)
	} else {
		tokenString, err = token.SignedString(mw.Key)
	}
	return tokenString, err
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the FaygoJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *FaygoJWTMiddleware) RefreshHandler(c *faygo.Context) error {
	token, err := mw.parseToken(c)
	//如果有错误并不是过期则触发错误并返回
	if (err != nil) && (err.Error() != "Token is expired") {
		faygo.Error(err)
		return mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
	}
	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))
	//如果超过 刷新期则触发过期错误并返回
	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		return mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, c))

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
	newClaims["orig_iat"] = mw.TimeFunc().Unix() //原来是 origIat貌似不正确
	tokenString, err := mw.signedString(newToken)

	if err != nil {
		return mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, c))
	}

	return mw.RefreshResponse(c, http.StatusOK, tokenString, expire)
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *faygo.Context) jwt.MapClaims {
	if exists := c.HasData("JWT_PAYLOAD"); !exists {
		return make(jwt.MapClaims)
	}

	jwtClaims := c.Data("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)

	/*claims, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)*/
}

// TokenGenerator method that clients can use to get a jwt token.
func (mw *FaygoJWTMiddleware) TokenGenerator(userID string) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().UTC().Add(mw.Timeout)
	claims["id"] = userID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expire, nil
}

func (mw *FaygoJWTMiddleware) jwtFromHeader(c *faygo.Context, key string) (string, error) {
	authHeader := c.HeaderParam(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (mw *FaygoJWTMiddleware) jwtFromQuery(c *faygo.Context, key string) (string, error) {
	token := c.QueryParam(key)

	if token == "" {
		return "", ErrEmptyQueryToken
	}

	return token, nil
}

func (mw *FaygoJWTMiddleware) jwtFromCookie(c *faygo.Context, key string) (string, error) {
	cookie := c.CookieParam(key)

	if cookie == "" {
		return "", ErrEmptyCookieToken
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
			return nil, ErrInvalidSigningAlgorithm
		}
		if mw.usingPublicKeyAlgo() {
			return mw.pubKey, nil
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

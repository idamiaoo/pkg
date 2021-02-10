package auth

import (
	"time"
)

var signKey = "test"

func NewAuthMiddleware() *GinJWTMiddleware {
	// the jwt middleware
	authMiddleware := &GinJWTMiddleware{
		Realm:      "wisroom zone",
		Key:        []byte(signKey),
		Timeout:    time.Hour * 24 * 30,
		MaxRefresh: time.Hour * 24 * 30,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenLookup: []string{"header:Authorization", "cookie:token", "query:token"},

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:   time.Now,
		SendCookie: true,
	}
	authMiddleware.MiddlewareInit()
	return authMiddleware
}

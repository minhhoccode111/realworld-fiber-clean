//nolint:revive
package common

// ctxKey is a unexported type to prevent key collisions
type ctxKey string

const (
	// CtxUserIDKey is used to store and retrieve the user ID from the Fiber context locals.
	CtxUserIDKey ctxKey = "userID"
	// CtxUserRoleKey is used to store and retrieve the user role from the Fiber context locals.
	CtxUserRoleKey ctxKey = "userRole"
	// CtxIsAuthKey is used to store and retrieve the authentication status from the Fiber context locals.
	CtxIsAuthKey ctxKey = "isAuth"

	// CookieJWTName is the name of the JWT cookie.
	CookieJWTName string = "realworld-jwt"

	// AuthorizationScheme is the expected authorization scheme for JWT tokens.
	// e.g. "Token", "Bearer"
	AuthorizationScheme string = "Token"
)

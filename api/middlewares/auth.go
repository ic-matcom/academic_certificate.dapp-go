package middlewares

import (
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
	"dapp/schema/dto"
)

// NewAuthCheckerMiddleware Bearer Authentication token verification middleware
func NewAuthCheckerMiddleware(sigKey []byte) context.Handler {
	checker := jwt.NewVerifier(jwt.HS256, sigKey)
	checker.WithDefaultBlocklist() // Enable server-side token block feature (even before its expiration time):
	// checker.WithDecryption()

	return checker.Verify(func() interface{} {
		// We can add login here

		return new(dto.AccessTokenData)
	})
}

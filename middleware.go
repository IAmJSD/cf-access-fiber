package access

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// CloudflareAccessUserInfo is used to define the user information.
type CloudflareAccessUserInfo struct {
	// Email is the e-mail address of the authenticated Cloudflare Access user.
	Email string `json:"email"`

	// UserIdentifier is the unique identifier for the authenticated Cloudflare Access user.
	UserIdentifier string `json:"sub"`

	// ExpiresAt is the timestamp when the token expires.
	ExpiresAt int `json:"exp"`

	// IssuedAt is the timestamp when the token was issued.
	IssuedAt int `json:"iat"`
}

// Validate is the middleware which is used to validate a request is valid.
func Validate(TeamDomain, ApplicationAUD string, UnauthorizedHandler fiber.Handler) fiber.Handler {
	// Handle if the domain ends with a slash or doesn't start with https://.
	if strings.HasSuffix(TeamDomain, "/") {
		TeamDomain = TeamDomain[:1]
	}
	if !strings.HasPrefix(TeamDomain, "https://") {
		TeamDomain = "https://" + TeamDomain
	}

	// Defines the OIDC verifier.
	keySet := oidc.NewRemoteKeySet(context.TODO(), TeamDomain+"/cdn-cgi/access/certs")
	verifier := oidc.NewVerifier(TeamDomain, keySet, &oidc.Config{
		ClientID: ApplicationAUD,
	})

	// Return the middleware.
	return func(ctx *fiber.Ctx) error {
		// Check the header exists.
		accessHeader := string(ctx.Request().Header.Peek("Cf-Access-Jwt-Assertion"))
		if accessHeader != "" {
			if token, err := verifier.Verify(ctx.Context(), accessHeader); err == nil {
				// Get the user information.
				var x CloudflareAccessUserInfo
				if err = token.Claims(&x); err != nil {
					return err
				}

				// Store the user data and call the function.
				ctx.Locals("user", &x)
				return ctx.Next()
			}
		}

		// Return the unauthorized handler if we fall out of this.
		return UnauthorizedHandler(ctx)
	}
}

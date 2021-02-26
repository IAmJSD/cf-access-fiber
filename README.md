# cf-access-fiber

A simple bit of Fiber middleware to use Cloudflare Access:
```go
// Setup the team domain/application AUD (this can be a env variable).
const (
	applicationAUD = "4714c1358e65fe4b408ad6d432a5f878f08194bdb4752441fd56faefa9b2b6f2"
	teamDomain     = "jakegealer.cloudflareaccess.com"
)

// Called when the request is unauthorized.
func unauthorized(ctx *fiber.Ctx) error {
    ctx.Status(400)
    _, err := ctx.WriteString("No content for you!")
    return err
}

// Initialises the middleware.
var middleware = access.Validate(
    teamDomain, applicationAUD, unauthorized)

// ...then

// Insert the middleware. You could do this per route too.
app.Use(middleware)

// ...then (if you need information about the user in the request)

user := ctx.Locals("user").(*access.CloudflareAccessUserInfo)
```

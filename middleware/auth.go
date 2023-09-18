package middleware

import (
	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/biskitsx/Refresh-Token-Go/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware interface {
	VerifyUser(c *fiber.Ctx) error
}
type authMiddleware struct {
	container  container.Container
	jwtService service.JwtService
}

func NewAuthMiddleware(container container.Container) AuthMiddleware {
	return &authMiddleware{
		container:  container,
		jwtService: service.NewJwtService(container),
	}
}

func (a authMiddleware) VerifyUser(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")

	// If access_token is not found, you might want to return an error or handle it differently
	if accessToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Access token is missing")
	}

	// Verify the access token
	decodedAccessToken, err := a.jwtService.VerifyToken(accessToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Access token is invalid")
	}

	// Extract the payload from the token
	payload, ok := decodedAccessToken.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}
	// Extract the user ID and username from the payload
	userPayload, ok := payload["user"].(map[string]interface{})
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user payload in token")
	}

	userIDFloat, ok := userPayload["user_id"].(float64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
	}

	username, ok := userPayload["username"].(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username in token")
	}

	userID := uint(userIDFloat)

	c.Locals("user_id", userID)
	c.Locals("username", username)

	return c.Next()
}

package middleware

import (
	"fmt"

	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/biskitsx/Refresh-Token-Go/service"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware interface {
	DeserializeUser(c *fiber.Ctx) error
	RequireUser(c *fiber.Ctx) error
}
type authMiddleware struct {
	container      container.Container
	jwtService     service.JwtService
	sessionService service.SessionService
}

func NewAuthMiddleware(container container.Container) AuthMiddleware {
	return &authMiddleware{
		container:      container,
		jwtService:     service.NewJwtService(container),
		sessionService: service.NewSessionService(container),
	}
}

func (a *authMiddleware) DeserializeUser(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token")

	if accessToken == "" {
		return c.Next()
	}
	decodedAccessToken, err := a.jwtService.VerifyToken(accessToken)
	if err != nil {
		if err.Error() == "Token is expired" {
			refreshToken := c.Cookies("refresh_token")
			decodedRefreshToken, err := a.jwtService.VerifyToken(refreshToken)
			if err != nil {
				fmt.Println(err.Error())
				return c.Next()
			}

			sessionId, err := a.jwtService.ExtractRefreshToken(decodedRefreshToken)
			if err != nil {
				fmt.Println(err.Error())
				return c.Next()
			}

			session := a.sessionService.GetSession(sessionId)
			newAccessToken, err := a.jwtService.GenerateTokenBySession(session, 5000)

			if err != nil {
				fmt.Println(err.Error())
				return c.Next()
			}

			newAccessCookie := &fiber.Cookie{
				Name:   "access_token",
				Value:  newAccessToken,
				MaxAge: 600,
			}
			c.Cookie(newAccessCookie)

			c.Locals("user_id", session.UserID)
			c.Locals("username", session.Username)
			fmt.Println("rotate refresh token")
		}
		return c.Next()

	}

	userId, username, err := a.jwtService.ExtractAccessToken(decodedAccessToken)
	if err != nil {
		return c.Next()
	}

	c.Locals("user_id", userId)
	c.Locals("username", username)

	return c.Next()
}

func (a *authMiddleware) RequireUser(c *fiber.Ctx) error {
	userId := c.Locals("user_id")
	if userId == nil {
		return fiber.NewError(401, "unauthenticate")
	}
	return c.Next()
}

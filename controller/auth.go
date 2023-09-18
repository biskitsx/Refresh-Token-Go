package controller

import (
	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/biskitsx/Refresh-Token-Go/model"
	"github.com/biskitsx/Refresh-Token-Go/service"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type authController struct {
	container      container.Container
	authService    service.AuthService
	sessionService service.SessionService
	jwtService     service.JwtService
}

func NewAuthController(container container.Container) AuthController {
	return &authController{
		container:      container,
		authService:    service.NewAuthService(container),
		sessionService: service.NewSessionService(container),
		jwtService:     service.NewJwtService(container),
	}
}

type UserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccessTokenDto struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type RefreshTokenDto struct {
	SessionID uint `json:"session_id"`
}

func (controller authController) Login(c *fiber.Ctx) error {
	userLogin := UserDto{}

	if err := c.BodyParser(&userLogin); err != nil {
		return fiber.NewError(401, err.Error())
	}

	oldUser, err := controller.authService.CheckUsername(userLogin.Username)
	if err != nil {
		return fiber.NewError(401, "invalid username or password")
	}

	if oldUser.Password != userLogin.Password {
		return fiber.NewError(401, "invalid password")
	}

	// Create Session
	session := controller.sessionService.CreateSession(oldUser.ID)

	accessToken, err := controller.jwtService.GenerateToken(&AccessTokenDto{UserID: oldUser.ID, Username: oldUser.Username}, "5s")
	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	refreshToken, err := controller.jwtService.GenerateToken(&RefreshTokenDto{SessionID: session.ID}, "1h")
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	accessCookie := &fiber.Cookie{
		Name:   "access_token",
		Value:  accessToken,
		MaxAge: 10,
	}
	refreshCookie := &fiber.Cookie{
		Name:   "refresh_token",
		Value:  refreshToken,
		MaxAge: 600,
	}

	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.JSON(fiber.Map{
		"msg": "login successfully",
	})
}

func (controller authController) Register(c *fiber.Ctx) error {
	dto := UserDto{}

	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(401, err.Error())
	}

	_, err := controller.authService.CheckUsername(dto.Username)
	if err == nil {
		return fiber.NewError(401, "This username already registered")
	}

	user := model.User{
		Username: dto.Username,
		Password: dto.Password,
	}

	db := controller.container.GetDatabase()
	db.Create(&user)
	return c.JSON(user)
}

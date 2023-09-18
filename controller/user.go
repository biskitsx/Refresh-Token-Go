package controller

import (
	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/biskitsx/Refresh-Token-Go/model"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	GetUser(c *fiber.Ctx) error
	GetMe(c *fiber.Ctx) error
}
type userController struct {
	container container.Container
}

func NewUserController(c container.Container) UserController {
	return &userController{
		container: c,
	}
}

func (controller userController) GetUser(c *fiber.Ctx) error {
	users := []model.User{}
	db := controller.container.GetDatabase()
	db.Find(&users)
	return c.JSON(users)
}

func (controller userController) GetMe(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"msg": "you are logged in",
	})
}

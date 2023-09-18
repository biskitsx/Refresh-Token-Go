package main

import (
	"github.com/biskitsx/Refresh-Token-Go/config"
	"github.com/biskitsx/Refresh-Token-Go/container"
	"github.com/biskitsx/Refresh-Token-Go/controller"
	"github.com/biskitsx/Refresh-Token-Go/database"
	"github.com/biskitsx/Refresh-Token-Go/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	var (
		app = fiber.New(fiber.Config{
			ErrorHandler: config.CustomErrorHandler,
		})
		db             = database.ConnectDb()
		container      = container.NewContainer(db)
		authController = controller.NewAuthController(container)
		userController = controller.NewUserController(container)
		authMiddleware = middleware.NewAuthMiddleware(container)
	)

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "Access-Control-Allow-Origin, Content-Type, Origin, Accept",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	app.Use(authMiddleware.DeserializeUser)
	// routes
	app.Get("/user", userController.GetUser)

	app.Post("/register", authController.Register)
	app.Post("/login", authController.Login)

	app.Get("/me", authMiddleware.RequireUser, userController.GetMe)

	app.Listen(":8080")
}

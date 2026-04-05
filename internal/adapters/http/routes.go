package http

import (
	"github.com/Nishishei01/Go_Hexagonal/internal/adapters/http/middleware"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	"github.com/gofiber/fiber/v3"
)

func Routes(app *fiber.App, userHandler AuthHandler, authService services.AuthService) {

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
	app.Post("/refresh-token", userHandler.RefreshToken)

	app.Use(middleware.JWTMiddleware(authService))
	app.Get("/test", func(c fiber.Ctx) {
		c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Hello World!!"})
	})
}

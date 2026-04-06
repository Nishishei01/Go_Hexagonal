package http

import (
	"github.com/Nishishei01/Go_Hexagonal/internal/adapters/http/middleware"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	"github.com/gofiber/fiber/v3"
)

func Routes(app *fiber.App, authHandler AuthHandler, authService services.AuthService, userHandler UserHandler, userService services.UserService, postHandler PostHandler, roleHandler RoleHandler) {

	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh-token", authHandler.RefreshToken)

	user := app.Group("/user")
	user.Use(middleware.JWTMiddleware(authService))
	// user.Get("/test", func(c fiber.Ctx) {
	// 	c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Hello World!!"})
	// })
	user.Get("/", userHandler.GetAll)
	user.Get("/:id", userHandler.GetByID)
	user.Put("/:id", userHandler.Update)
	user.Delete("/:id", userHandler.Delete)

	post := app.Group("/post")
	post.Use(middleware.JWTMiddleware(authService))
	post.Post("/", postHandler.Create)
	post.Get("/", postHandler.GetAll)
	post.Get("/:id", postHandler.GetByID)
	post.Put("/:id", postHandler.Update)
	post.Delete("/:id", postHandler.Delete)

	role := app.Group("/role")
	role.Use(middleware.JWTMiddleware(authService))
	role.Post("/", roleHandler.Create)
	role.Get("/", roleHandler.GetAll)
	role.Get("/:id", roleHandler.GetByID)
	role.Put("/:id", roleHandler.Update)
	role.Delete("/:id", roleHandler.Delete)
}

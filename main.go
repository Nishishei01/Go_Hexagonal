package main

import (
	"os"

	"github.com/Nishishei01/Go_Hexagonal/internal/config"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	config.LoadEnv()

	dsn := config.DbUrl()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate()

	app.Get("/test", func(c fiber.Ctx) error {
		return c.SendString(os.Getenv("DB_USER"))
	})

	app.Listen(":8080")
}

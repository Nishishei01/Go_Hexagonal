package main

import (
	gormAdapter "github.com/Nishishei01/Go_Hexagonal/internal/adapters/gorm"
	httpAdapter "github.com/Nishishei01/Go_Hexagonal/internal/adapters/http"
	"github.com/Nishishei01/Go_Hexagonal/internal/config"
	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	config.LoadEnv()

	dsn := config.DbUrl()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(
		&domains.User{},
	)

	userRepo := gormAdapter.NewAuthGormRepository(db)
	userService := services.NewAuthService(userRepo)
	userHandler := httpAdapter.NewAuthHandler(userService)

	httpAdapter.Routes(app, userHandler, userService)

	app.Listen(":8080")
}

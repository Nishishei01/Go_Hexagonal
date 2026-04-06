package main

import (
	"context"

	gormAdapter "github.com/Nishishei01/Go_Hexagonal/internal/adapters/gorm"
	httpAdapter "github.com/Nishishei01/Go_Hexagonal/internal/adapters/http"
	redisAdapter "github.com/Nishishei01/Go_Hexagonal/internal/adapters/redis"
	"github.com/Nishishei01/Go_Hexagonal/internal/config"
	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
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
		&domains.Role{},
		&domains.Post{},
	)

	redisUrl := config.RedisUrl()
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic("Failed to parse Redis URL: " + err.Error())
	}
	redisClient := redis.NewClient(opt)
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	cacheRepo := redisAdapter.NewRedisCacheAdapter(redisClient)

	authRepo := gormAdapter.NewAuthGormRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := httpAdapter.NewAuthHandler(authService)

	userRepo := gormAdapter.NewUserGormRepository(db)
	userService := services.NewUserService(userRepo, cacheRepo)
	userHandler := httpAdapter.NewUserHandler(userService)

	postRepo := gormAdapter.NewPostGormRepository(db)
	postService := services.NewPostService(postRepo, cacheRepo)
	postHandler := httpAdapter.NewPostHandler(postService)

	roleRepo := gormAdapter.NewRoleGormRepository(db)
	roleService := services.NewRoleService(roleRepo, cacheRepo)
	roleHandler := httpAdapter.NewRoleHandler(roleService)

	httpAdapter.Routes(app, authHandler, authService, userHandler, userService, postHandler, roleHandler)

	app.Listen(":8080")
}

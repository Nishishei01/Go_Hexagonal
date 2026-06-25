package main

import (
	"context"
	"log"
	"net"

	gormAdapter "github.com/Nishishei01/Go_Hexagonal/internal/adapters/gorm"
	"github.com/Nishishei01/Go_Hexagonal/internal/config"
	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	"google.golang.org/grpc"

	grpcAdapter "github.com/Nishishei01/Go_Hexagonal/internal/adapters/grpc"
	pb "github.com/Nishishei01/Go_Hexagonal/proto/auth"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

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

	// cacheRepo := redisAdapter.NewRedisCacheAdapter(redisClient)

	authRepo := gormAdapter.NewAuthGormRepository(db)
	authService := services.NewAuthService(authRepo)
	// authHandler := httpAdapter.NewAuthHandler(authService)
	authGrpcHandler := grpcAdapter.NewAuthGrpcHandler(authService)

	// userRepo := gormAdapter.NewUserGormRepository(db)
	// userService := services.NewUserService(userRepo, cacheRepo)
	// userHandler := httpAdapter.NewUserHandler(userService)

	// postRepo := gormAdapter.NewPostGormRepository(db)
	// postService := services.NewPostService(postRepo, cacheRepo)
	// postHandler := httpAdapter.NewPostHandler(postService)

	// roleRepo := gormAdapter.NewRoleGormRepository(db)
	// roleService := services.NewRoleService(roleRepo, cacheRepo)
	// roleHandler := httpAdapter.NewRoleHandler(roleService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authGrpcHandler)

	log.Println("gRPC Server listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

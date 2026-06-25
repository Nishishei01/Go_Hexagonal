package grpc

import (
	"context"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	pb "github.com/Nishishei01/Go_Hexagonal/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcHandler struct {
	pb.UnimplementedAuthServiceServer
	authGrpcService services.AuthService
}

func NewAuthGrpcHandler(authService services.AuthService) *AuthGrpcHandler {
	return &AuthGrpcHandler{authGrpcService: authService}
}

func (a *AuthGrpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	domainReq := &domains.RegisterRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}

	err := a.authGrpcService.Register(domainReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register:; %v", err)
	}

	res := &pb.RegisterResponse{
		Message: "User created successfully!",
	}

	return res, nil
}

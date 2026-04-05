package domains

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=25"`
	Password  string `json:"password" validate:"required,min=3"`
	Email     string `json:"email" validate:"required,email,min=3,max=25"`
	FirstName string `json:"firstName" validate:"required,min=3,max=25"`
	LastName  string `json:"lastName" validate:"required,min=3,max=25"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=25"`
	Password string `json:"password" validate:"required,min=3"`
}

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

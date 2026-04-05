package services

import (
	"errors"
	"os"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *domains.RegisterRequest) error
	Login(loginUser *domains.LoginRequest) (string, string, error)
	ValidateToken(tokenString string) (*domains.JWTClaims, error)
	RefreshToken(refreshTokenString string) (string, string, error)
}

type AuthServiceImpl struct {
	authRepo ports.AuthRepository
}

func NewAuthService(authRepo ports.AuthRepository) AuthService {
	return &AuthServiceImpl{authRepo: authRepo}
}

func (a *AuthServiceImpl) Register(authUser *domains.RegisterRequest) error {
	_, err := a.authRepo.GetUserByUsername(authUser.Username)
	if err == nil {
		return errors.New("This username has already exists!")
	}

	_, err = a.authRepo.GetUserByEmail(authUser.Email)
	if err == nil {
		return errors.New("This email has already exists!")
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	authUser.Password = string(hasedPassword)

	if err := a.authRepo.CreateUser(authUser); err != nil {
		return err
	}
	return nil
}

func (a *AuthServiceImpl) Login(authUser *domains.LoginRequest) (string, string, error) {

	user, err := a.authRepo.GetUserByUsername(authUser.Username)
	if err != nil {
		return "", "", errors.New("Username or Password invalid!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password))
	if err != nil {
		return "", "", errors.New("Password invalid!")
	}

	return a.generateTokens(user.ID, user.Username)
}

func (a *AuthServiceImpl) ValidateToken(tokenString string) (*domains.JWTClaims, error) {

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &domains.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domains.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid Token!")
}

func (a *AuthServiceImpl) RefreshToken(refreshTokenString string) (string, string, error) {
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	token, err := jwt.ParseWithClaims(refreshTokenString, &domains.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(*domains.JWTClaims); ok && token.Valid {
		return a.generateTokens(claims.UserID, claims.Username)
	}

	return "", "", errors.New("Invalid Refresh Token!")
}

func (a *AuthServiceImpl) generateTokens(userID uint, username string) (string, string, error) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	accessClaims := domains.JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}

	refreshClaims := domains.JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

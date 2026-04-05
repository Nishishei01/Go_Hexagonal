package http

import (
	"fmt"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return AuthHandler{authService: authService}
}

func (a *AuthHandler) Register(c fiber.Ctx) error {
	var register domains.RegisterRequest
	if err := c.Bind().Body(&register); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request!"})
	}

	if err := a.authService.Register(&register); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully!"})
}

func (a *AuthHandler) Login(c fiber.Ctx) error {
	var login domains.LoginRequest
	if err := c.Bind().Body(&login); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request!"})
	}

	accessToken, refreshToken, err := a.authService.Login(&login)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	a.setRefreshTokenCookie(c, refreshToken)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
		"message":      "User successfully logged in!",
	})
}

func (a *AuthHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing refresh token!"})
	}

	accessToken, newRefreshToken, err := a.authService.RefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired refresh token!"})
	}

	a.setRefreshTokenCookie(c, newRefreshToken)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
		"message":      "Token refreshed successfully!",
	})
}

func (a *AuthHandler) setRefreshTokenCookie(c fiber.Ctx, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})
}

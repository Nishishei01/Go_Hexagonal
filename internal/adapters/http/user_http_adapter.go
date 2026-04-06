package http

import (
	"fmt"
	"strconv"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userservice services.UserService) UserHandler {
	return UserHandler{userService: userservice}
}

func (u *UserHandler) GetAll(c fiber.Ctx) error {
	users, err := u.userService.GetAll()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func (u *UserHandler) GetByID(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	user, err := u.userService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not have this id"})
		}
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (u *UserHandler) Update(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	var userRequest domains.UserRequest
	if err := c.Bind().Body(&userRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request!"})
	}

	if err := u.userService.Update(uint(id), &userRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully!"})
}

func (u *UserHandler) Delete(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	if err := u.userService.Delete(uint(id)); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully!"})
}

package http

import (
	"fmt"
	"strconv"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	"github.com/gofiber/fiber/v3"
)

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) RoleHandler {
	return RoleHandler{roleService: roleService}
}

func (h *RoleHandler) Create(c fiber.Ctx) error {
	var roleRequest domains.RoleRequest
	if err := c.Bind().Body(&roleRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request!"})
	}

	if err := h.roleService.Create(&roleRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Role created successfully!"})
}

func (h *RoleHandler) GetAll(c fiber.Ctx) error {
	roles, err := h.roleService.GetAll()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}
	return c.Status(fiber.StatusOK).JSON(roles)
}

func (h *RoleHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Role not found"})
		}
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}
	return c.Status(fiber.StatusOK).JSON(role)
}

func (h *RoleHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	var roleRequest domains.RoleRequest
	if err := c.Bind().Body(&roleRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request!"})
	}

	if err := h.roleService.Update(uint(id), &roleRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Role updated successfully!"})
}

func (h *RoleHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	if err := h.roleService.Delete(uint(id)); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Role deleted successfully!"})
}

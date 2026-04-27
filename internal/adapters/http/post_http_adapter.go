package http

import (
	"fmt"
	"strconv"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

type PostHandler struct {
	postService services.PostService
}

func NewPostHandler(postService services.PostService) PostHandler {
	return PostHandler{postService: postService}
}

func (h *PostHandler) Create(c fiber.Ctx) error {
	var postRequest domains.PostRequest
	if err := c.Bind().Body(&postRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format!"})
	}

	if err := validate.Struct(&postRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed!"})
	}

	if err := h.postService.Create(&postRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Post created successfully!"})
}

func (h *PostHandler) GetAll(c fiber.Ctx) error {
	posts, err := h.postService.GetAll()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}

func (h *PostHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error!"})
	}
	return c.Status(fiber.StatusOK).JSON(post)
}

func (h *PostHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	var postRequest domains.PostRequest
	if err := c.Bind().Body(&postRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format!"})
	}

	if err := validate.Struct(&postRequest); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed!"})
	}

	if err := h.postService.Update(uint(id), &postRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Post updated successfully!"})
}

func (h *PostHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 0 {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID!"})
	}

	if err := h.postService.Delete(uint(id)); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Post deleted successfully!"})
}

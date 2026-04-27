package http

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPostService struct {
	mock.Mock
}

func (m *mockPostService) Create(postRequest *domains.PostRequest) error {
	args := m.Called(postRequest)
	return args.Error(0)
}

func (m *mockPostService) GetAll() ([]*domains.Post, error) {
	args := m.Called()
	return args.Get(0).([]*domains.Post), args.Error(1)
}

func (m *mockPostService) GetByID(id uint) (*domains.Post, error) {
	args := m.Called(id)
	var post *domains.Post
	if args.Get(0) != nil {
		post = args.Get(0).(*domains.Post)
	}
	return post, args.Error(1)
}

func (m *mockPostService) Update(id uint, postRequest *domains.PostRequest) error {
	args := m.Called(id, postRequest)
	return args.Error(0)
}

func (m *mockPostService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestPostHandler_Create(t *testing.T) {
	mockService := new(mockPostService)
	handler := NewPostHandler(mockService)

	app := fiber.New()
	app.Post("/post", handler.Create)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Create", mock.Anything).Return(nil)

		req := httptest.NewRequest("POST", "/post", bytes.NewBufferString(`{"title": "Test", "content": "Content", "user_id": 1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body format (user_id as string)", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := httptest.NewRequest("POST", "/post", bytes.NewBufferString(`{"title": "Test", "content": "Content", "user_id": "abc"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing required fields (validation failed)", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := httptest.NewRequest("POST", "/post", bytes.NewBufferString(`{"title": "Test", "content": "Content"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := httptest.NewRequest("POST", "/post", bytes.NewBufferString(`invalid json`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Create", mock.Anything).Return(errors.New("service error"))

		req := httptest.NewRequest("POST", "/post", bytes.NewBufferString(`{"title": "Test", "content": "Content", "user_id": 1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestPostHandler_GetAll(t *testing.T) {
	mockService := new(mockPostService)
	handler := NewPostHandler(mockService)

	app := fiber.New()
	app.Get("/post", handler.GetAll)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetAll").Return([]*domains.Post{{ID: 1, Title: "Test"}}, nil)

		req := httptest.NewRequest("GET", "/post", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetAll").Return([]*domains.Post{}, errors.New("service error"))

		req := httptest.NewRequest("GET", "/post", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestPostHandler_GetByID(t *testing.T) {
	mockService := new(mockPostService)
	handler := NewPostHandler(mockService)

	app := fiber.New()
	app.Get("/post/:id", handler.GetByID)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetByID", uint(1)).Return(&domains.Post{ID: 1, Title: "Test"}, nil)

		req := httptest.NewRequest("GET", "/post/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := httptest.NewRequest("GET", "/post/abc", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetByID", uint(1)).Return(nil, errors.New("record not found"))

		req := httptest.NewRequest("GET", "/post/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestPostHandler_Update(t *testing.T) {
	mockService := new(mockPostService)
	handler := NewPostHandler(mockService)

	app := fiber.New()
	app.Put("/post/:id", handler.Update)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Update", uint(1), mock.Anything).Return(nil)

		req := httptest.NewRequest("PUT", "/post/1", bytes.NewBufferString(`{"title": "Updated", "content": "Updated content", "user_id": 1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestPostHandler_Delete(t *testing.T) {
	mockService := new(mockPostService)
	handler := NewPostHandler(mockService)

	app := fiber.New()
	app.Delete("/post/:id", handler.Delete)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Delete", uint(1)).Return(nil)

		req := httptest.NewRequest("DELETE", "/post/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

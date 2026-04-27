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

type mockRoleService struct {
	mock.Mock
}

func (m *mockRoleService) Create(roleRequest *domains.RoleRequest) error {
	args := m.Called(roleRequest)
	return args.Error(0)
}

func (m *mockRoleService) GetAll() ([]*domains.Role, error) {
	args := m.Called()
	return args.Get(0).([]*domains.Role), args.Error(1)
}

func (m *mockRoleService) GetByID(id uint) (*domains.Role, error) {
	args := m.Called(id)
	var role *domains.Role
	if args.Get(0) != nil {
		role = args.Get(0).(*domains.Role)
	}
	return role, args.Error(1)
}

func (m *mockRoleService) Update(id uint, roleRequest *domains.RoleRequest) error {
	args := m.Called(id, roleRequest)
	return args.Error(0)
}

func (m *mockRoleService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestRoleHandler_Create(t *testing.T) {
	mockService := new(mockRoleService)
	handler := NewRoleHandler(mockService)

	app := fiber.New()
	app.Post("/role", handler.Create)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Create", mock.Anything).Return(nil)

		req := httptest.NewRequest("POST", "/role", bytes.NewBufferString(`{"role": "Test"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body format (role as number)", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := httptest.NewRequest("POST", "/role", bytes.NewBufferString(`{"role": 1}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	})

	t.Run("missing required fields (validation failed)", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := httptest.NewRequest("POST", "/role", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		req := httptest.NewRequest("POST", "/role", bytes.NewBufferString(`invalid json`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("server error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Create", mock.Anything).Return(errors.New("service error"))

		req := httptest.NewRequest("POST", "/role", bytes.NewBufferString(`{"role": "Test"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

}

func TestRoleHandler_GetAll(t *testing.T) {
	mockService := new(mockRoleService)
	handler := NewRoleHandler(mockService)

	app := fiber.New()
	app.Get("/role", handler.GetAll)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetAll").Return([]*domains.Role{{ID: 1, Role: "Test"}}, nil)

		req := httptest.NewRequest("GET", "/role", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetAll").Return([]*domains.Role{}, errors.New("service error"))

		req := httptest.NewRequest("GET", "/role", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestRoleHandler_GetByID(t *testing.T) {
	mockService := new(mockRoleService)
	handler := NewRoleHandler(mockService)

	app := fiber.New()
	app.Get("/role/:id", handler.GetByID)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetByID", uint(1)).Return(&domains.Role{ID: 1, Role: "Test"}, nil)

		req := httptest.NewRequest("GET", "/role/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		req := httptest.NewRequest("GET", "/role/abc", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetByID", uint(1)).Return(nil, errors.New("record not found"))

		req := httptest.NewRequest("GET", "/role/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

}

func TestRoleHandler_Update(t *testing.T) {
	mockService := new(mockRoleService)
	handler := NewRoleHandler(mockService)

	app := fiber.New()
	app.Put("/role/:id", handler.Update)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Update", uint(1), mock.Anything).Return(nil)

		req := httptest.NewRequest("PUT", "/role/1", bytes.NewBufferString(`{"role": "Updated"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestRoleHandler_Delete(t *testing.T) {
	mockService := new(mockRoleService)
	handler := NewRoleHandler(mockService)

	app := fiber.New()
	app.Delete("/role/:id", handler.Delete)

	t.Run("success", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Delete", uint(1)).Return(nil)

		req := httptest.NewRequest("DELETE", "/role/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

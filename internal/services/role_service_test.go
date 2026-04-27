package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/stretchr/testify/assert"
)

type mockRoleRepo struct {
	CreateRoleFunc  func(roleRequest *domains.RoleRequest) error
	GetAllRoleFunc  func() ([]*domains.Role, error)
	GetRoleByIDFunc func(id uint) (*domains.Role, error)
	UpdateRoleFunc  func(id uint, roleRequest *domains.RoleRequest) error
	DeleteRoleFunc  func(id uint) error
}

func (m *mockRoleRepo) CreateRole(roleRequest *domains.RoleRequest) error {
	if m.CreateRoleFunc != nil {
		return m.CreateRoleFunc(roleRequest)
	}
	return nil
}

func (m *mockRoleRepo) GetAllRole() ([]*domains.Role, error) {
	if m.GetAllRoleFunc != nil {
		return m.GetAllRoleFunc()
	}
	return nil, nil
}

func (m *mockRoleRepo) GetRoleByID(id uint) (*domains.Role, error) {
	if m.GetRoleByIDFunc != nil {
		return m.GetRoleByIDFunc(id)
	}
	return nil, nil
}

func (m *mockRoleRepo) UpdateRole(id uint, roleRequest *domains.RoleRequest) error {
	if m.UpdateRoleFunc != nil {
		return m.UpdateRoleFunc(id, roleRequest)
	}
	return nil
}

func (m *mockRoleRepo) DeleteRole(id uint) error {
	if m.DeleteRoleFunc != nil {
		return m.DeleteRoleFunc(id)
	}
	return nil
}

type mockRoleCacheRepo struct {
	SetFunc    func(key string, value any, expiration time.Duration) error
	GetFunc    func(key string, dest any) error
	DeleteFunc func(key string) error
}

func (m *mockRoleCacheRepo) Set(key string, value any, expiration time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(key, value, expiration)
	}
	return nil
}

func (m *mockRoleCacheRepo) Get(key string, dest any) error {
	if m.GetFunc != nil {
		return m.GetFunc(key, dest)
	}
	return nil
}

func (m *mockRoleCacheRepo) Delete(key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(key)
	}
	return nil
}

func TestCreateRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockRoleRepo{
			CreateRoleFunc: func(roleRequest *domains.RoleRequest) error {
				return nil
			},
		}

		cacheRepo := &mockRoleCacheRepo{
			DeleteFunc: func(key string) error {
				return nil
			},
		}

		service := NewRoleService(repo, cacheRepo)

		err := service.Create(&domains.RoleRequest{Role: "Unite Test"})
		assert.NoError(t, err)
	})

	t.Run("failure from repository", func(t *testing.T) {
		expectedErr := errors.New("database error")
		repo := &mockRoleRepo{
			CreateRoleFunc: func(roleRequest *domains.RoleRequest) error {
				return expectedErr
			},
		}

		cacheRepo := &mockRoleCacheRepo{}

		service := NewRoleService(repo, cacheRepo)

		err := service.Create(&domains.RoleRequest{Role: "Unite Test"})
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

}

func TestGetAllRole(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		expectedRoles := []*domains.Role{
			{
				ID:   1,
				Role: "Cached Role",
			},
		}

		cacheRepo := &mockRoleCacheRepo{
			GetFunc: func(key string, dest any) error {
				if p, ok := dest.(*[]*domains.Role); ok {
					*p = expectedRoles
				}
				return nil
			},
		}
		repo := &mockRoleRepo{}

		service := NewRoleService(repo, cacheRepo)
		roles, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
	})

	t.Run("cache miss, fetch from repo success", func(t *testing.T) {
		expectedRoles := []*domains.Role{
			{
				ID:   1,
				Role: "Repo Role",
			},
		}

		repo := &mockRoleRepo{
			GetAllRoleFunc: func() ([]*domains.Role, error) {
				return expectedRoles, nil
			},
		}

		cacheRepo := &mockRoleCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("Cache miss")
			},
			SetFunc: func(key string, value any, expiration time.Duration) error {
				return nil
			},
		}

		service := NewRoleService(repo, cacheRepo)
		roles, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
	})

	t.Run("cache miss, fetch from repo failure", func(t *testing.T) {
		expectedErr := errors.New("db error")
		cacheRepo := &mockRoleCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("cache miss")
			},
		}

		repo := &mockRoleRepo{
			GetAllRoleFunc: func() ([]*domains.Role, error) {
				return nil, expectedErr
			},
		}

		service := NewRoleService(repo, cacheRepo)
		roles, err := service.GetAll()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, roles)
	})

}

func TestGetRoleByID(t *testing.T) {
	t.Run("cache hit", func(t *testing.T) {
		expectedRole := &domains.Role{ID: 1, Role: "Unit Test"}
		cacheRepo := &mockRoleCacheRepo{
			GetFunc: func(key string, dest any) error {
				if p, ok := dest.(**domains.Role); ok {
					*p = expectedRole
				}
				return nil
			},
		}
		repo := &mockRoleRepo{}

		service := NewRoleService(repo, cacheRepo)
		role, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, role)
	})

	t.Run("cache miss, fetch from repo success", func(t *testing.T) {
		expectedRole := &domains.Role{ID: 1, Role: "Unit Test"}

		cacheRepo := &mockRoleCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("cache miss")
			},
			SetFunc: func(key string, value any, expiration time.Duration) error {
				return nil
			},
		}

		repo := &mockRoleRepo{
			GetRoleByIDFunc: func(id uint) (*domains.Role, error) {
				return expectedRole, nil
			},
		}

		service := NewRoleService(repo, cacheRepo)
		role, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, role)
	})

	t.Run("cache miss, fetch from repo failure", func(t *testing.T) {
		expectedErr := errors.New("db error")

		cacheRepo := &mockRoleCacheRepo{
			GetFunc: func(key string, dest any) error {
				return errors.New("cache miss")
			},
		}

		repo := &mockRoleRepo{
			GetRoleByIDFunc: func(id uint) (*domains.Role, error) {
				return nil, expectedErr
			},
		}

		service := NewRoleService(repo, cacheRepo)
		role, err := service.GetByID(1)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, role)
	})
}

func TestUpdateRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockRoleRepo{
			UpdateRoleFunc: func(id uint, roleRequest *domains.RoleRequest) error {
				return nil
			},
		}
		cachRepo := &mockRoleCacheRepo{
			DeleteFunc: func(key string) error {
				return nil
			},
		}
		service := NewRoleService(repo, cachRepo)

		err := service.Update(1, &domains.RoleRequest{Role: "Updated"})
		assert.NoError(t, err)
	})
}

func TestDeleteRole(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockRoleRepo{
			DeleteRoleFunc: func(id uint) error {
				return nil
			},
		}
		cacheRepo := &mockRoleCacheRepo{
			DeleteFunc: func(key string) error {
				return nil
			},
		}
		service := NewRoleService(repo, cacheRepo)

		err := service.Delete(1)
		assert.NoError(t, err)
	})
}

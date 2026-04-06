package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
)

type RoleService interface {
	Create(roleRequest *domains.RoleRequest) error
	GetAll() ([]*domains.Role, error)
	GetByID(id uint) (*domains.Role, error)
	Update(id uint, roleRequest *domains.RoleRequest) error
	Delete(id uint) error
}

type RoleServiceImpl struct {
	roleRepo  ports.RoleRepository
	cacheRepo ports.CacheRepository
}

func NewRoleService(roleRepo ports.RoleRepository, cacheRepo ports.CacheRepository) RoleService {
	return &RoleServiceImpl{
		roleRepo:  roleRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *RoleServiceImpl) Create(roleRequest *domains.RoleRequest) error {
	err := s.roleRepo.CreateRole(roleRequest)
	if err == nil {
		s.cacheRepo.Delete("roles:all")
	}
	return err
}

func (s *RoleServiceImpl) GetAll() ([]*domains.Role, error) {
	cacheKey := "roles:all"
	var roles []*domains.Role

	err := s.cacheRepo.Get(cacheKey, &roles)
	if err == nil {
		log.Println("Hit cache for", cacheKey)
		return roles, nil
	}

	roles, err = s.roleRepo.GetAllRole()
	if err != nil {
		return nil, err
	}

	err = s.cacheRepo.Set(cacheKey, roles, 5*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return roles, nil
}

func (s *RoleServiceImpl) GetByID(id uint) (*domains.Role, error) {
	cacheKey := fmt.Sprintf("role:%d", id)
	var role *domains.Role

	err := s.cacheRepo.Get(cacheKey, &role)
	if err == nil {
		log.Println("Hit cache for", cacheKey)
		return role, nil
	}

	role, err = s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}

	err = s.cacheRepo.Set(cacheKey, role, 5*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return role, nil
}

func (s *RoleServiceImpl) Update(id uint, roleRequest *domains.RoleRequest) error {
	err := s.roleRepo.UpdateRole(id, roleRequest)
	if err == nil {
		s.cacheRepo.Delete(fmt.Sprintf("role:%d", id))
		s.cacheRepo.Delete("roles:all")
	}
	return err
}

func (s *RoleServiceImpl) Delete(id uint) error {
	err := s.roleRepo.DeleteRole(id)
	if err == nil {
		s.cacheRepo.Delete(fmt.Sprintf("role:%d", id))
		s.cacheRepo.Delete("roles:all")
	}
	return err
}

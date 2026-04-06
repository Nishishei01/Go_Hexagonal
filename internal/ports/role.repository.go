package ports

import "github.com/Nishishei01/Go_Hexagonal/internal/domains"

type RoleRepository interface {
	CreateRole(roleRequest *domains.RoleRequest) error
	GetAllRole() ([]*domains.Role, error)
	GetRoleByID(id uint) (*domains.Role, error)
	UpdateRole(id uint, roleRequest *domains.RoleRequest) error
	DeleteRole(id uint) error
}

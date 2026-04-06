package gorm

import (
	"log"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
	"gorm.io/gorm"
)

type RoleGormRepository struct {
	db *gorm.DB
}

func NewRoleGormRepository(db *gorm.DB) ports.RoleRepository {
	return &RoleGormRepository{db: db}
}

func (r *RoleGormRepository) CreateRole(roleRequest *domains.RoleRequest) error {
	role := &domains.Role{
		Role: roleRequest.Role,
	}

	if result := r.db.Create(role); result.Error != nil {
		log.Printf("Error to create role: %v", result.Error)
		return result.Error
	}

	log.Println("Create role successfully!")
	return nil
}

func (r *RoleGormRepository) GetAllRole() ([]*domains.Role, error) {
	var roles []*domains.Role

	if result := r.db.Preload("Users").Order("create_at desc").Find(&roles); result.Error != nil {
		log.Printf("Error to get all role: %v", result.Error)
		return nil, result.Error
	}

	for _, role := range roles {
		if role.Users == nil {
			role.Users = []domains.User{}
		}
	}

	log.Println("Get role successfully!")
	return roles, nil
}

func (r *RoleGormRepository) GetRoleByID(id uint) (*domains.Role, error) {
	var role *domains.Role

	if result := r.db.Preload("Users").Where("id = ?", id).First(&role); result.Error != nil {
		log.Printf("Error to get role by id: %v", result.Error)
		return nil, result.Error
	}

	if role.Users == nil {
		role.Users = []domains.User{}
	}

	log.Println("Get role successfully!")
	return role, nil
}

func (r *RoleGormRepository) UpdateRole(id uint, roleRequest *domains.RoleRequest) error {
	result := r.db.Model(&domains.Role{}).Where("id = ?", id).Updates(roleRequest)
	if result.Error != nil {
		log.Printf("Error to update role: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	log.Println("Update role successfully!")
	return nil
}

func (r *RoleGormRepository) DeleteRole(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domains.Role{ID: id}).Association("Users").Clear(); err != nil {
			return err
		}

		result := tx.Where("id = ?", id).Delete(&domains.Role{})
		if result.Error != nil {
			log.Printf("Error to delete role: %v", result.Error)
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		log.Println("Delete role successfully!")
		return nil
	})
}

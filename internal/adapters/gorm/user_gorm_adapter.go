package gorm

import (
	"log"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) ports.UserRepository {
	return &UserGormRepository{db: db}
}

func (u *UserGormRepository) GetAllUser() ([]*domains.User, error) {
	var users []*domains.User

	if result := u.db.Preload("Posts").Preload("Roles").Order("create_at desc").Find(&users); result.Error != nil {
		log.Printf("Error to get all user: %v", result.Error)
		return nil, result.Error
	}

	for _, user := range users {
		if user.Posts == nil {
			user.Posts = []domains.Post{}
		}
		if user.Roles == nil {
			user.Roles = []domains.Role{}
		}
	}

	log.Println("Get user successfully!")
	return users, nil
}

func (u *UserGormRepository) GetUserByID(id uint) (*domains.User, error) {
	var user *domains.User

	if result := u.db.Preload("Posts").Preload("Roles").Where("id = ?", id).First(&user); result.Error != nil {
		log.Printf("Error to get user by id: %v", result.Error)
		return nil, result.Error
	}

	log.Println("Get user successfully!")

	if user.Posts == nil {
		user.Posts = []domains.Post{}
	}
	if user.Roles == nil {
		user.Roles = []domains.Role{}
	}

	return user, nil
}

func (u *UserGormRepository) UpdateUser(id uint, userRequest *domains.UserRequest) error {

	return u.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domains.User{}).Where("id = ?", id).Updates(userRequest)
		if result.Error != nil {
			log.Printf("Error to update user: %v", result.Error)
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if len(userRequest.RoleIDs) > 0 {
			var roles []domains.Role
			if err := tx.Where("id IN ?", userRequest.RoleIDs).Find(&roles).Error; err != nil {
				return err
			}

			if err := tx.Model(&domains.User{ID: id}).Association("Roles").Replace(roles); err != nil {
				return err
			}
		} else if userRequest.RoleIDs != nil {
			if err := tx.Model(&domains.User{ID: id}).Association("Roles").Clear(); err != nil {
				return err
			}
		}

		log.Println("Update user successfully!")
		return nil
	})
}

func (u *UserGormRepository) DeleteUser(id uint) error {
	result := u.db.Where("id = ?", id).Delete(&domains.User{})
	if result.Error != nil {
		log.Printf("Error to update uesr: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	log.Println("Update user successfully!")
	return nil
}

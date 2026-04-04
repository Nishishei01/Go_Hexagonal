package gorm

import (
	"log"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
	"gorm.io/gorm"
)

type AuthGormRepository struct {
	db *gorm.DB
}

func NewAuthGormRepository(db *gorm.DB) ports.AuthRepository {
	return &AuthGormRepository{db: db}
}

func (d *AuthGormRepository) CreateUser(reqUser *domains.RegisterRequest) error {

	user := domains.User{
		Username:  reqUser.Username,
		Password:  reqUser.Password,
		Email:     reqUser.Email,
		FirstName: reqUser.FirstName,
		LastName:  reqUser.LastName,
	}

	if results := d.db.Create(&user); results.Error != nil {
		log.Panicf("Error to register: %v", results.Error)
		return results.Error
	}
	log.Println("Register successfully!")
	return nil
}

func (d *AuthGormRepository) GetUserByEmail(email string) (*domains.User, error) {
	var user domains.User
	if results := d.db.Where("email = ?", email).First(&user); results.Error != nil {
		if results.Error != gorm.ErrRecordNotFound {
			log.Printf("Error to get user by email: %v", results.Error)
		}
		return nil, results.Error
	}
	return &user, nil
}

func (d *AuthGormRepository) GetUserByUsername(username string) (*domains.User, error) {
	var user domains.User
	if results := d.db.Where("username = ?", username).First(&user); results.Error != nil {
		if results.Error != gorm.ErrRecordNotFound {
			log.Printf("Error to get user by username: %v", results.Error)
		}
		return nil, results.Error
	}
	return &user, nil
}

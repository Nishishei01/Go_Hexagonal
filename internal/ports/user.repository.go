package ports

import "github.com/Nishishei01/Go_Hexagonal/internal/domains"

type UserRepository interface {
	GetAllUser() ([]*domains.User, error)
	GetUserByID(id uint) (*domains.User, error)
	UpdateUser(id uint, userRequest *domains.UserRequest) error
	DeleteUser(id uint) error
}

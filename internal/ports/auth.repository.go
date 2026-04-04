package ports

import "github.com/Nishishei01/Go_Hexagonal/internal/domains"

type AuthRepository interface {
	CreateUser(user *domains.RegisterRequest) error
	GetUserByEmail(email string) (*domains.User, error)
	GetUserByUsername(username string) (*domains.User, error)
}

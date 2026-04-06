package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
)

type UserService interface {
	GetAll() ([]*domains.User, error)
	GetByID(id uint) (*domains.User, error)
	Update(id uint, userRequest *domains.UserRequest) error
	Delete(id uint) error
}

type UserServiceImpl struct {
	userRepo  ports.UserRepository
	cacheRepo ports.CacheRepository
}

func NewUserService(userRepo ports.UserRepository, cacheRepo ports.CacheRepository) UserService {
	return &UserServiceImpl{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

func (u *UserServiceImpl) GetAll() ([]*domains.User, error) {
	cacheKey := "users:all"
	var users []*domains.User

	err := u.cacheRepo.Get(cacheKey, &users)
	if err == nil {
		log.Println("Hit cache for", cacheKey)
		return users, nil
	}

	users, err = u.userRepo.GetAllUser()
	if err != nil {
		return nil, err
	}

	err = u.cacheRepo.Set(cacheKey, users, 5*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return users, nil
}

func (u *UserServiceImpl) GetByID(id uint) (*domains.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	var user domains.User

	err := u.cacheRepo.Get(cacheKey, &user)
	if err == nil {
		log.Println("Hit cache for", cacheKey)
		return &user, nil
	}

	userPtr, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	err = u.cacheRepo.Set(cacheKey, userPtr, 5*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return userPtr, nil
}

func (u *UserServiceImpl) Update(id uint, userRequest *domains.UserRequest) error {
	err := u.userRepo.UpdateUser(id, userRequest)
	if err == nil {

		u.cacheRepo.Delete(fmt.Sprintf("user:%d", id))
		u.cacheRepo.Delete("users:all")
	}
	return err
}

func (u *UserServiceImpl) Delete(id uint) error {
	err := u.userRepo.DeleteUser(id)
	if err == nil {

		u.cacheRepo.Delete(fmt.Sprintf("user:%d", id))
		u.cacheRepo.Delete("users:all")
	}
	return err
}

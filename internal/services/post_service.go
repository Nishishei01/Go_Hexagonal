package services

import (
	"fmt"
	"log"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
)

type PostService interface {
	Create(postRequest *domains.PostRequest) error
	GetAll() ([]*domains.Post, error)
	GetByID(id uint) (*domains.Post, error)
	Update(id uint, postRequest *domains.PostRequest) error
	Delete(id uint) error
}

type PostServiceImpl struct {
	postRepo  ports.PostRepository
	cacheRepo ports.CacheRepository
}

func NewPostService(postRepo ports.PostRepository, cacheRepo ports.CacheRepository) PostService {
	return &PostServiceImpl{
		postRepo:  postRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *PostServiceImpl) Create(postRequest *domains.PostRequest) error {
	err := s.postRepo.CreatePost(postRequest)
	if err == nil {
		s.cacheRepo.Delete("posts:all")
	}
	return err
}

func (s *PostServiceImpl) GetAll() ([]*domains.Post, error) {
	cacheKey := "posts:all"
	var posts []*domains.Post

	err := s.cacheRepo.Get(cacheKey, &posts)
	if err == nil {
		log.Println("Hit cache for", cacheKey)
		return posts, nil
	}

	posts, err = s.postRepo.GetAllPost()
	if err != nil {
		return nil, err
	}

	err = s.cacheRepo.Set(cacheKey, posts, 5*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return posts, nil
}

func (s *PostServiceImpl) GetByID(id uint) (*domains.Post, error) {
	cacheKey := fmt.Sprintf("post:%d", id)
	var post *domains.Post

	err := s.cacheRepo.Get(cacheKey, &post)
	if err == nil {
		log.Println("Hit cache for", cacheKey)
		return post, nil
	}

	post, err = s.postRepo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	err = s.cacheRepo.Set(cacheKey, post, 5*time.Minute)
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	return post, nil
}

func (s *PostServiceImpl) Update(id uint, postRequest *domains.PostRequest) error {
	err := s.postRepo.UpdatePost(id, postRequest)
	if err == nil {
		s.cacheRepo.Delete(fmt.Sprintf("post:%d", id))
		s.cacheRepo.Delete("posts:all")
	}
	return err
}

func (s *PostServiceImpl) Delete(id uint) error {
	err := s.postRepo.DeletePost(id)
	if err == nil {
		s.cacheRepo.Delete(fmt.Sprintf("post:%d", id))
		s.cacheRepo.Delete("posts:all")
	}
	return err
}

package ports

import "github.com/Nishishei01/Go_Hexagonal/internal/domains"

type PostRepository interface {
	CreatePost(postRequest *domains.PostRequest) error
	GetAllPost() ([]*domains.Post, error)
	GetPostByID(id uint) (*domains.Post, error)
	UpdatePost(id uint, postRequest *domains.PostRequest) error
	DeletePost(id uint) error
}

package gorm

import (
	"log"

	"github.com/Nishishei01/Go_Hexagonal/internal/domains"
	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
	"gorm.io/gorm"
)

type PostGormRepository struct {
	db *gorm.DB
}

func NewPostGormRepository(db *gorm.DB) ports.PostRepository {
	return &PostGormRepository{db: db}
}

func (p *PostGormRepository) CreatePost(postRequest *domains.PostRequest) error {
	post := &domains.Post{
		Title:   postRequest.Title,
		Content: postRequest.Content,
		UserID:  postRequest.UserID,
	}

	if result := p.db.Create(post); result.Error != nil {
		log.Printf("Error to create post: %v", result.Error)
		return result.Error
	}

	log.Println("Create post successfully!")
	return nil
}

func (p *PostGormRepository) GetAllPost() ([]*domains.Post, error) {
	var posts []*domains.Post

	if result := p.db.Order("create_at desc").Find(&posts); result.Error != nil {
		log.Printf("Error to get all post: %v", result.Error)
		return nil, result.Error
	}

	log.Println("Get post successfully!")
	return posts, nil
}

func (p *PostGormRepository) GetPostByID(id uint) (*domains.Post, error) {
	var post *domains.Post

	if result := p.db.Where("id = ?", id).First(&post); result.Error != nil {
		log.Printf("Error to get post by id: %v", result.Error)
		return nil, result.Error
	}

	log.Println("Get post successfully!")
	return post, nil
}

func (p *PostGormRepository) UpdatePost(id uint, postRequest *domains.PostRequest) error {
	result := p.db.Model(&domains.Post{}).Where("id = ?", id).Updates(postRequest)
	if result.Error != nil {
		log.Printf("Error to update post: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	log.Println("Update post successfully!")
	return nil
}

func (p *PostGormRepository) DeletePost(id uint) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ?", id).Delete(&domains.Post{})
		if result.Error != nil {
			log.Printf("Error to delete post: %v", result.Error)
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		log.Println("Delete post successfully!")
		return nil
	})
}

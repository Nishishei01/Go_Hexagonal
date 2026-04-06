package domains

import "time"

type Post struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"autoUpdateTime"`

	UserID uint `json:"user_id"`
}

type PostRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	UserID  uint   `json:"user_id" validate:"required"`
}

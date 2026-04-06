package domains

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreateAt  time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt  time.Time `json:"updateAt" gorm:"autoUpdateTime"`

	Posts []Post `json:"posts" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Roles []Role `json:"roles" gorm:"many2many:user_roles"`
}

type UserRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`

	RoleIDs []uint `json:"role_ids" gorm:"-"`
}

package domains

import "time"

type Role struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Role     string    `json:"role"`
	CreateAt time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"autoUpdateTime"`

	Users []User `json:"users" gorm:"many2many:user_roles"`
}

type RoleRequest struct {
	Role string `json:"role" validate:"required"`
}

package role

import "time"

type Role struct {
	ID        string    `gorm:"primarykey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp"`
	DeletedAt time.Time `gorm:"index;default:null"`
}

type RoleRequest struct {
	Name string `json:"name" validate:"required~Name cannot be empty"`
}

type RoleResponse struct {
	ID        string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

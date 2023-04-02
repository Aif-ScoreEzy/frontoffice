package role

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        string         `gorm:"primarykey"`
	Name      string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RoleRequest struct {
	ID        string    `json:"-"`
	Name      string    `json:"name" validate:"required~Name cannot be empty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type RoleResponse struct {
	ID        string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

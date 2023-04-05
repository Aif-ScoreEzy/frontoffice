package permission

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        string         `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"-"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PermissionRequest struct {
	ID        string    `json:"-"`
	Name      string    `json:"name" validate:"required~Name cannot be empty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

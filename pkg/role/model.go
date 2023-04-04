package role

import (
	"front-office/pkg/permission"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string                  `gorm:"primarykey"`
	Name        string                  `gorm:"not null"`
	Permissions []permission.Permission `gorm:"many2many:role_permissions"`
	CreatedAt   time.Time               `gorm:"not null;default:current_timestamp"`
	UpdatedAt   time.Time               `gorm:"not null;default:current_timestamp"`
	DeletedAt   gorm.DeletedAt          `gorm:"index"`
}

type RoleRequest struct {
	ID          string                  `json:"-"`
	Name        string                  `json:"name" validate:"required~Name cannot be empty"`
	Permissions []permission.Permission `json:"permissions"`
	CreatedAt   time.Time               `json:"-"`
	UpdatedAt   time.Time               `json:"-"`
	DeletedAt   time.Time               `json:"-"`
}

type RoleResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions"`
	CreatedAt   time.Time               `json:"-"`
	UpdatedAt   time.Time               `json:"-"`
	DeletedAt   gorm.DeletedAt          `json:"-"`
}

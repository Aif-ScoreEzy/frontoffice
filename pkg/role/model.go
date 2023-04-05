package role

import (
	"front-office/pkg/permission"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string                  `gorm:"primarykey" json:"id"`
	Name        string                  `gorm:"not null" json:"name"`
	Permissions []permission.Permission `gorm:"many2many:role_permissions" json:"permissions"`
	CreatedAt   time.Time               `gorm:"not null;default:current_timestamp" json:"-"`
	UpdatedAt   time.Time               `gorm:"not null;default:current_timestamp" json:"-"`
	DeletedAt   gorm.DeletedAt          `gorm:"index" json:"-"`
}

type RoleRequest struct {
	ID          string                  `json:"-"`
	Name        string                  `json:"name" validate:"required~Name cannot be empty"`
	Permissions []permission.Permission `json:"permissions"`
	CreatedAt   time.Time               `json:"-"`
	UpdatedAt   time.Time               `json:"-"`
	DeletedAt   time.Time               `json:"-"`
}

package role

import (
	"front-office/pkg/core/permission"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Id          string                  `gorm:"primarykey" json:"id"`
	Name        string                  `gorm:"not null" json:"name"`
	Permissions []permission.Permission `gorm:"many2many:role_permissions" json:"permissions"`
	TierLevel   uint                    `gorm:"not null" json:"tier_level"`
	CreatedAt   time.Time               `gorm:"not null;default:current_timestamp" json:"-"`
	UpdatedAt   time.Time               `gorm:"not null;default:current_timestamp" json:"-"`
	DeletedAt   gorm.DeletedAt          `gorm:"index" json:"-"`
}

type MstRole struct {
	RoleId      uint            `json:"role_id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name"`
	Permissions []MstPermission `json:"permissions" gorm:"many2many:ref_role_permissions"`
}

type MstPermission struct {
	PermissionId uint   `json:"permission_id" gorm:"primaryKey;autoIncrement"`
	Slug         string `json:"slug"`
	Name         string `json:"name"`
}

type CreateRoleRequest struct {
	Name        string                  `json:"name" validate:"required~Field Name is required"`
	Permissions []permission.Permission `json:"permissions" validate:"required~Field Permissions is required"`
	TierLevel   uint                    `json:"tier_level" validate:"required~Field Tier Level is required, range(0|2)~Field Tier Level is only available in the range of 0 to 2."`
}

type CreateRoleResponse struct {
	Id          string                  `json:"-"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions"`
	TierLevel   uint                    `json:"tier_level"`
	CreatedAt   time.Time               `json:"-"`
	UpdatedAt   time.Time               `json:"-"`
	DeletedAt   time.Time               `json:"-"`
}

type UpdateRoleRequest struct {
	Id          string                  `json:"-"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions"`
	TierLevel   uint                    `json:"tier_level" validate:"max=3~only available with tier level 1, 2"`
}

type AifResponse struct {
	Success bool    `json:"success"`
	Data    MstRole `json:"data"`
	Message string  `json:"message"`
	Meta    any     `json:"meta,omitempty"`
	Status  bool    `json:"status,omitempty"`
}

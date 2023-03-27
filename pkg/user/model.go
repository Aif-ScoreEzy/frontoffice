package user

import (
	"time"

	pkgRole "front-office/pkg/role"

	"gorm.io/gorm"
)

type User struct {
	ID        string        `gorm:"primarykey" json:"-"`
	Name      string        `gorm:"not null" json:"name"`
	Username  string        `gorm:"not null" json:"username"`
	Email     string        `gorm:"not null" json:"email"`
	Password  string        `gorm:"not null" json:"-"`
	RoleID    string        `gorm:"not null" json:"_"`
	Role      *pkgRole.Role `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time     `gorm:"not null;default:current_timestamp" json:"-"`
	UpdatedAt time.Time     `gorm:"not null;default:current_timestamp" json:"-"`
	DeletedAt time.Time     `gorm:"index" json:"-"`
}

type UserRequest struct {
	ID        string         `json:"-"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type UserResponse struct {
	ID        string         `json:"-"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	RoleID    string         `json:"-"`
	Role      *pkgRole.Role  `json:"role"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

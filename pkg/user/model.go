package user

import (
	"front-office/pkg/company"
	"front-office/pkg/role"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string          `json:"id"`
	Name      string          `json:"name" validate:"required,alphaspace, min(3)"`
	Username  string          `json:"username" gorm:"unique"`
	Email     string          `json:"email" gorm:"unique" validate:"required,email"`
	Password  string          `json:"password"`
	Phone     string          `json:"phone"`
	Key       string          `json:"key"`
	Active    bool            `json:"active"`
	CompanyID string          `json:"-"`
	Company   company.Company `json:"company" gorm:"foreignKey:CompanyID"`
	RoleID    string          `json:"-"`
	Role      role.Role       `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

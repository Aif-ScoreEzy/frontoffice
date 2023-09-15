package user

import (
	"front-office/pkg/company"
	"front-office/pkg/role"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          string          `gorm:"primarykey" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Email       string          `gorm:"unique" json:"email"`
	Password    string          `gorm:"not null" json:"password"`
	Phone       string          `json:"phone"`
	AccountType string          `json:"account_type"`
	Key         string          `json:"key"`
	Active      bool            `json:"active"`
	IsVerified  bool            `gorm:"default:false" json:"is_verified"`
	CompanyID   string          `json:"company_id"`
	Company     company.Company `gorm:"foreignKey:CompanyID" json:"company"`
	RoleID      string          `gorm:"not null" json:"role_id"`
	Role        role.Role       `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

type UserResponse struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Password   string          `json:"-"`
	Phone      string          `json:"phone"`
	Active     bool            `json:"active"`
	IsVerified bool            `json:"is_verified"`
	CompanyID  string          `json:"-"`
	Company    company.Company `json:"company"`
	RoleID     string          `json:"-"`
	Role       role.Role       `json:"role"`
	CreatedAt  time.Time       `json:"-"`
	UpdatedAt  time.Time       `json:"-"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

type RegisterMemberRequest struct {
	Name   string `json:"name" validate:"required~Field Name is required"`
	Email  string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
	RoleID string `json:"role_id" validate:"required~Field Name is required"`
}

type UpdateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email" validate:"email~Only email pattern are allowed"`
	Phone     string `string:"phone" validate:"phone"`
	RoleID    string `json:"role_id"`
	CompanyID string `json:"company_id"`
}

type UserUpdateResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Phone     string `json:"phone"`
	Active    bool   `json:"active"`
	CompanyID string `json:"company_id"`
	RoleID    string `json:"role_id"`
}

func SetPassword(password string) string {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	password = string(hashedPass)

	return password
}

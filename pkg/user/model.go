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
	Status      string          `gorm:"default:pending" json:"status"`
	Active      bool            `gorm:"default:false" json:"active"`
	IsVerified  bool            `gorm:"default:false" json:"is_verified"`
	Image       string          `gorm:"default:default-profile-image.jpg" json:"image"`
	CompanyID   string          `json:"company_id"`
	Company     company.Company `gorm:"foreignKey:CompanyID" json:"company"`
	RoleID      string          `gorm:"not null" json:"role_id"`
	Role        role.Role       `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

type ActivationToken struct {
	ID         string         `gorm:"primarykey" json:"id"`
	Token      string         `gorm:"not null" json:"token"`
	Activation bool           `gorm:"not null;default:false" json:"activation"`
	UserID     string         `json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserResponse struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Password   string          `json:"-"`
	Phone      string          `json:"phone"`
	Status     string          `json:"status"`
	Active     bool            `json:"active"`
	IsVerified bool            `json:"is_verified"`
	CompanyID  string          `json:"-"`
	Company    company.Company `json:"company"`
	RoleID     string          `json:"-"`
	Role       role.Role       `json:"role"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"-"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

type RegisterMemberRequest struct {
	Name   string `json:"name" validate:"required~Field Name is required"`
	Email  string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
	RoleID string `json:"role_id" validate:"required~Field Role is required"`
}

type ActivationAccountRequest struct {
	Email string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
}

type GetUsersResponse struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	Status     string         `json:"status"`
	Active     bool           `json:"active"`
	IsVerified bool           `json:"is_verified"`
	CompanyID  string         `json:"company_id"`
	RoleID     string         `json:"-"`
	Role       role.Role      `json:"role"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type UpdateUserRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email" validate:"email~Only email pattern are allowed"`
	RoleID *string `json:"role_id"`
	Active *bool   `json:"active"`
	Status *string `json:"status"`
}

type UserUpdateResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Status    string `json:"status"`
	Active    bool   `json:"active"`
	CompanyID string `json:"company_id"`
	RoleID    string `json:"role_id"`
}

func SetPassword(password string) string {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	password = string(hashedPass)

	return password
}

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
	IsVerified  bool            `json:"is_verified"`
	CompanyID   string          `json:"company_id"`
	Company     company.Company `gorm:"foreignKey:CompanyID" json:"company"`
	RoleID      string          `gorm:"not null" json:"role_id"`
	Role        role.Role       `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

type RegisterUserRequest struct {
	Name            string `json:"name" validate:"required~Name cannot be empty"`
	Email           string `json:"email" validate:"required~Email cannot be empty, email~Only email pattern are allowed"`
	Password        string `json:"password" validate:"required~Password cannot be empty, min(8)~Password must have at least 8 characters"`
	Phone           string `string:"phone" validate:"required~Phone cannot be empty, phone"`
	CompanyName     string `json:"company_name"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	AgreementNumber string `json:"agreement_number"`
	IndustryID      string `json:"industry_id"`
	PaymentScheme   string `json:"payment_scheme"`
	RoleID          string `json:"role_id" validate:"required~Role cannot be empty"`
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

type SendEmailVerificationRequest struct {
	Email string `json:"email" validate:"required~Email cannot be empty"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required~Email cannot be empty"`
}

type RegisterMemberRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleID string `json:"role_id" validate:"required~Role cannot be empty"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required~Email cannot be empty"`
	Password string `json:"password" validate:"required~Password cannot be empty"`
}

type UserLoginResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
	RoleID    string `json:"role_id"`
	Token     string `json:"access_token"`
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

package auth

import (
	"front-office/pkg/company"
	"front-office/pkg/role"
	"front-office/pkg/user"
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID         string         `gorm:"primarykey" json:"id"`
	Token      string         `gorm:"not null" json:"token"`
	Activation bool           `gorm:"not null;default:false" json:"activation"`
	UserID     string         `json:"user_id"`
	User       user.User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type RegisterAdminRequest struct {
	Name            string `json:"name" validate:"required~Field Name is required"`
	Email           string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
	Password        string `json:"password" validate:"required~Field Password is required, min(8)~Field Password must have at least 8 characters"`
	Phone           string `string:"phone" validate:"required~Field Phone is required, phone"`
	CompanyName     string `json:"company_name"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	AgreementNumber string `json:"agreement_number"`
	IndustryID      string `json:"industry_id"`
	PaymentScheme   string `json:"payment_scheme"`
	RoleID          string `json:"role_id" validate:"required~Field Role is required"`
}

type RegisterAdminResponse struct {
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
	CreatedAt  time.Time       `json:"-"`
	UpdatedAt  time.Time       `json:"-"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required~Field Email is required"`
	Password string `json:"password" validate:"required~Field Password is required"`
}

type UserLoginResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`
	TierLevel   uint   `json:"tier_level"`
	Image       string `json:"image"`
	Token       string `json:"access_token"`
}

type SendEmailVerificationRequest struct {
	Email string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
}

type PasswordResetRequest struct {
	Password        string `json:"password" validate:"required~Field Password is required, min(8)~Field Password must have at least 8 characters"`
	ConfirmPassword string `json:"confirm_password" validate:"required~Field Confirm Password is required"`
}

type ChangePasswordRequest struct {
	CurrentPassword    string `json:"password" validate:"required~Field Current Password is required"`
	NewPassword        string `json:"new_password" validate:"required~Field New Password is required, min(8)~Field Password must have at least 8 characters"`
	ConfirmNewPassword string `json:"confirm_password" validate:"required~Field Confirmation New Password is required"`
}

type UpdateUserAuth struct {
	Status string `json:"status"`
}

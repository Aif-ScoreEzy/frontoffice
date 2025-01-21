package auth

import (
	"front-office/pkg/core/company"
	"front-office/pkg/core/role"
	"time"

	"gorm.io/gorm"
)

type RegisterAdminRequest struct {
	Name            string `json:"name" validate:"required~Field Name is required"`
	Email           string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
	Password        string `json:"password" validate:"required~Field Password is required, min(8)~Field Password must have at least 8 characters"`
	Phone           string `string:"phone" validate:"required~Field Phone is required, phone"`
	CompanyName     string `json:"company_name"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	AgreementNumber string `json:"agreement_number"`
	IndustryId      string `json:"industry_id"`
	PaymentScheme   string `json:"payment_scheme"`
	RoleId          string `json:"role_id" validate:"required~Field Role is required"`
}

type RegisterAdminResponse struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Password   string          `json:"-"`
	Phone      string          `json:"phone"`
	Status     string          `json:"status"`
	Active     bool            `json:"active"`
	IsVerified bool            `json:"is_verified"`
	CompanyId  string          `json:"-"`
	Company    company.Company `json:"company"`
	RoleId     string          `json:"-"`
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
	Id                 uint        `json:"id"`
	Name               string      `json:"name"`
	Email              string      `json:"email"`
	CompanyId          uint        `json:"company_id"`
	CompanyName        string      `json:"company_name"`
	TierLevel          uint        `json:"tier_level"`
	Image              string      `json:"image"`
	SubscriberProducts interface{} `json:"subscriber_products"`
}

type dataLoginResponse struct {
	MemberId           uint        `json:"member_id"`
	Name               string      `json:"name"`
	Email              string      `json:"email"`
	CompanyId          uint        `json:"company_id"`
	CompanyName        string      `json:"company_name"`
	RoleId             uint        `json:"role_id"`
	Image              string      `json:"image"`
	SubscriberProducts interface{} `json:"subscriber_products"`
}

type aifcoreAuthMemberResponse struct {
	Data       *dataLoginResponse `json:"data"`
	StatusCode int                `json:"status_code"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
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

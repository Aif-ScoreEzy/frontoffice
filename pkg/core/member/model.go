package member

import (
	"front-office/pkg/core/company"
	"front-office/pkg/core/role"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MstMember struct {
	MemberId          uint               `json:"member_id" gorm:"primaryKey;autoIncrement"`
	Name              string             `json:"name" validate:"required,alphaspace, min(3)"`
	Username          string             `json:"username" gorm:"unique"`
	Email             string             `json:"email" gorm:"unique" validate:"required,email"`
	Password          string             `json:"password"`
	Phone             string             `json:"phone"`
	Key               string             `json:"-" gorm:"uniqueIndex"`
	Active            bool               `json:"active"`
	ParentId          string             `json:"parent_id"`
	CompanyId         uint               `json:"company_id"`
	MstCompany        company.MstCompany `json:"company" gorm:"foreignKey:CompanyId"`
	RoleId            uint               `json:"role_id"`
	Role              role.MstRole       `json:"role" gorm:"foreignKey:RoleId"`
	Status            bool               `json:"status"`
	MailStatus        string             `json:"mail_status" gorm:"default:pending"`
	AccountType       string             `json:"account_type"`
	ProductPermission string             `json:"product_permission"`
	IsVerified        bool               `json:"is_verified"`
	Image             string             `json:"image" gorm:"default:default-profile-image.jpg"`
	QuotaType         int8               `json:"quota_type"` //0: none, 1: Quota Total 2: Quota per product
	Quota             int                `json:"quota"`
	CreatedAt         time.Time          `json:"-"`
	UpdatedAt         time.Time          `json:"-"`
	DeletedAt         gorm.DeletedAt     `json:"-" gorm:"index"`
}

type UserResponse struct {
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
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"-"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

type RegisterMemberRequest struct {
	Name      string `json:"name" validate:"required~Field Name is required"`
	Email     string `json:"email" validate:"required~Field Email is required, email~Only email pattern are allowed"`
	Username  string `json:"username"`
	Phone     string `json:"phone" validate:"phone"`
	CompanyId uint   `json:"company_id"`
	RoleId    uint   `json:"role_id"`
	Key       string `json:"key"`
}

type dataRegisterMemberResponse struct {
	MemberId uint `json:"member_id"`
}

type RegisterMemberResponse struct {
	Data       *dataRegisterMemberResponse `json:"data"`
	StatusCode int                         `json:"-"`
}

type GetUsersResponse struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	Status     string         `json:"status"`
	Active     bool           `json:"active"`
	IsVerified bool           `json:"is_verified"`
	CompanyId  string         `json:"company_id"`
	RoleId     string         `json:"-"`
	Role       role.Role      `json:"role"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type UpdateUserRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email" validate:"email~Only email pattern are allowed"`
	RoleId *string `json:"role_id"`
	Active *bool   `json:"active"`
	Status *string `json:"status"`
}

type UserUpdateResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	CompanyId uint   `json:"company_id"`
	RoleId    uint   `json:"role_id"`
}

type UpdateProfileRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type UploadProfileImageRequest struct {
	Image *string `json:"image"`
}

func SetPassword(password string) string {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	password = string(hashedPass)

	return password
}

type FindUserQuery struct {
	Id       string
	Email    string
	Username string
	Key      string
}

type FindUserAifCoreResponse struct {
	Message    string     `json:"message"`
	Success    bool       `json:"success"`
	Data       *MstMember `json:"data"`
	StatusCode int        `json:"-"`
}

type AifResponse struct {
	Success bool      `json:"success"`
	Data    MstMember `json:"data"`
	Message string    `json:"message"`
	Meta    any       `json:"meta,omitempty"`
	Status  bool      `json:"status,omitempty"`
}

type Meta struct {
	Total      any `json:"total,omitempty"`
	Page       any `json:"page,omitempty"`
	TotalPages any `json:"total_pages,omitempty"`
	Visible    any `json:"visible,omitempty"`
	StartData  any `json:"start_data,omitempty"`
	EndData    any `json:"end_data,omitempty"`
	Size       any `json:"size,omitempty"`
	Message    any `json:"message,omitempty"`
}
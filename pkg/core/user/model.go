package user

import (
	"front-office/pkg/core/company"
	"front-office/pkg/core/role"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id          string          `gorm:"primarykey" json:"id"`
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
	CompanyId   string          `json:"company_id"`
	Company     company.Company `gorm:"foreignKey:CompanyId" json:"company"`
	RoleId      string          `gorm:"not null" json:"role_id"`
	Role        role.Role       `gorm:"foreignKey:RoleId" json:"role"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

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
	CompanyId uint   `json:"company_id"`
	RoleId    uint   `json:"role_id"`
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
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Status    string `json:"status"`
	Active    bool   `json:"active"`
	CompanyId string `json:"company_id"`
	RoleId    string `json:"role_id"`
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

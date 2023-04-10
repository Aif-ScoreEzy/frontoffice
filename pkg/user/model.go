package user

import (
	"front-office/pkg/company"
	"front-office/pkg/role"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Username  string          `json:"username" gorm:"unique"`
	Email     string          `json:"email" gorm:"unique"`
	Password  string          `json:"password"`
	Phone     string          `json:"phone"`
	Key       string          `json:"key"`
	Active    bool            `json:"active"`
	CompanyID string          `json:"company_id"`
	Company   company.Company `json:"company" gorm:"foreignKey:CompanyID"`
	RoleID    string          `json:"role_id"`
	Role      role.Role       `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

type RegisterUserRequest struct {
	Name            string `json:"name" validate:"required~Name cannot be empty"`
	Email           string `json:"email" gorm:"unique" validate:"required~Email cannot be empty, email~Only email pattern are allowed"`
	Password        string `json:"password" validate:"required~Password cannot be empty, length(8)~Password must have at least 8 characters"`
	Username        string `json:"username" gorm:"unique" validate:"required~Username cannot be empty, alphanum~Only alphabet and numeric values are allowed for username"`
	Phone           string `string:"phone" validate:"required~Phone cannot be empty, phone"`
	CompanyName     string `json:"company_name"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	AgreementNumber string `json:"agreement_number"`
	IndustryID      string `json:"industry_id"`
	PaymentScheme   string `json:"payment_scheme"`
	RoleID          string `json:"role_id" validate:"required~Role cannot be empty"`
}

type RegisterUserResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Username  string          `json:"username" gorm:"unique"`
	Email     string          `json:"email" gorm:"unique"`
	Password  string          `json:"-"`
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

type UserLoginRequest struct {
	Username string `json:"username" validate:"required~Username cannot be empty"`
	Password string `json:"password" validate:"required~Password cannot be empty"`
}

type UserLoginResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
	Role      string `json:"role_id"`
	Key       string `json:"key"`
	Token     string `json:"access_token"`
}

type UpdateUserRequest struct {
	Name      string `json:"name" validate:"required~Name cannot be empty"`
	Username  string `json:"username" gorm:"unique" validate:"required~Username cannot be empty, alphanum~Only alphabet and numeric values are allowed for username"`
	Email     string `json:"email" gorm:"unique" validate:"required~Email cannot be empty, email~Only email pattern are allowed"`
	Phone     string `string:"phone" validate:"required~Phone cannot be empty, phone"`
	RoleID    string `json:"role_id" validate:"required~Role cannot be empty"`
	CompanyID string `json:"company_id" validate:"required~Company cannot be empty"`
}

func (user *User) SetPassword(password string) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = string(hashedPass)
}

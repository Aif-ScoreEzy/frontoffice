package company

import (
	"front-office/pkg/industry"
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID              string            `json:"id"`
	CompanyName     string            `json:"company_name"`
	CompanyAddress  string            `json:"company_address"`
	CompanyPhone    string            `json:"company_phone"`
	AgreementNumber string            `json:"agreement_number"`
	PaymentScheme   string            `json:"payment_scheme"`
	PostpaidActive  bool              `json:"active"`
	IndustryID      string            `json:"-"`
	Industry        industry.Industry `json:"industry" gorm:"foreignKey:IndustryID"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
}

type UpdateCompanyRequest struct {
	CompanyName     string `json:"company_name"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	AgreementNumber string `json:"agreement_number"`
	PaymentScheme   string `json:"payment_scheme"`
	IndustryID      string `json:"industry_id"`
}

type UpdateCompanyResponse struct {
	ID              string            `json:"id"`
	CompanyName     string            `json:"company_name"`
	CompanyAddress  string            `json:"company_address"`
	CompanyPhone    string            `json:"company_phone"`
	AgreementNumber string            `json:"agreement_number"`
	PaymentScheme   string            `json:"payment_scheme"`
	PostpaidActive  bool              `json:"active"`
	IndustryID      string            `json:"industry_id"`
	Industry        industry.Industry `json:"-"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	DeletedAt       gorm.DeletedAt    `json:"-"`
}

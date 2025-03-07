package company

import (
	"front-office/pkg/core/industry"
	"time"

	"gorm.io/gorm"
)

type Company struct {
	Id              string            `json:"id"`
	CompanyName     string            `json:"company_name"`
	CompanyAddress  string            `json:"company_address"`
	CompanyPhone    string            `json:"company_phone"`
	AgreementNumber string            `json:"agreement_number"`
	PaymentScheme   string            `json:"payment_scheme"`
	PostpaidActive  bool              `json:"active"`
	IndustryId      string            `json:"-"`
	Industry        industry.Industry `json:"industry" gorm:"foreignKey:IndustryId"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
}

type MstCompany struct {
	CompanyId       uint    `json:"company_id" gorm:"primaryKey;autoIncrement"`
	CompanyName     string  `json:"company_name"`
	CompanyAddress  string  `json:"company_address"`
	CompanyPhone    string  `json:"company_phone"`
	AgreementNumber string  `json:"agreement_number"`
	PaymentScheme   string  `json:"payment_scheme"`
	PostpaidActive  bool    `json:"active"`
	IndustryId      uint    `json:"industry_id"`
	BasePricing     float64 `json:"base_pricing"`
	// Apiconfigs      []apiconfig.MstApiconfig `json:"apiconfigs" gorm:"foreignKey:CompanyId"`
	// Products        []MstSubscribedProduct   `json:"products" gorm:"foreignKey:CompanyId"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UpdateCompanyRequest struct {
	CompanyName     string `json:"company_name"`
	CompanyAddress  string `json:"company_address"`
	CompanyPhone    string `json:"company_phone"`
	AgreementNumber string `json:"agreement_number"`
	PaymentScheme   string `json:"payment_scheme"`
	IndustryId      string `json:"industry_id"`
}

type UpdateCompanyResponse struct {
	Id              string            `json:"id"`
	CompanyName     string            `json:"company_name"`
	CompanyAddress  string            `json:"company_address"`
	CompanyPhone    string            `json:"company_phone"`
	AgreementNumber string            `json:"agreement_number"`
	PaymentScheme   string            `json:"payment_scheme"`
	PostpaidActive  bool              `json:"active"`
	IndustryId      string            `json:"industry_id"`
	Industry        industry.Industry `json:"-"`
	CreatedAt       time.Time         `json:"-"`
	UpdatedAt       time.Time         `json:"-"`
	DeletedAt       gorm.DeletedAt    `json:"-"`
}

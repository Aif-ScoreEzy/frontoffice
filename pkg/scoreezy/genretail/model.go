package genretail

import (
	"front-office/pkg/core/company"
	"front-office/pkg/core/user"

	"time"

	"gorm.io/gorm"
)

type GenRetailRequest struct {
	LoanNo   string `json:"loan_no"`
	Name     string `json:"name"`
	IDCardNo string `json:"id_card_no"`
	PhoneNo  string `json:"phone_no"`
}

type GenRetailV3ModelResponse struct {
	Message      string                 `json:"message"`
	ErrorMessage string                 `json:"error_message"`
	Success      bool                   `json:"success"`
	Data         *GenRetailV3DataClient `json:"data"`
	StatusCode   int                    `json:"status_code"`
}

type GenRetailV3DataClient struct {
	TransactionID        string  `json:"transaction_id"`
	Name                 string  `json:"name"`
	IDCardNo             string  `json:"id_card_no"`
	PhoneNo              string  `json:"phone_no"`
	LoanNo               string  `json:"loan_no"`
	ProbabilityToDefault float64 `json:"probability_to_default"`
	Grade                string  `json:"grade"`
	Date                 string  `json:"date"`
}

type GenRetailV3ClientReturnSuccess struct {
	Message string                 `json:"message"`
	Success bool                   `json:"success"`
	Data    *GenRetailV3DataClient `json:"data"`
}

type GenRetailV3ClientReturnError struct {
	Message      string                 `json:"message"`
	ErrorMessage string                 `json:"error_message"`
	Data         *GenRetailV3DataClient `json:"data"`
}

type UploadScoringRequest struct {
	Files []byte `json:"files"`
}

type UploadScoringReturnError struct {
	Message string `json:"message"`
}

type BulkSearch struct {
	ID                   uint            `gorm:"primarykey;autoIncrement" json:"id"`
	UploadID             string          `gorm:"not null" json:"upload_id"`
	TransactionID        string          `gorm:"not null" json:"transaction_id"`
	Name                 string          `gorm:"not null" json:"name"`
	IDCardNo             string          `gorm:"not null" json:"id_card_no"`
	PhoneNo              string          `gorm:"not null" json:"phone_no"`
	LoanNo               string          `gorm:"not null" json:"loan_no"`
	ProbabilityToDefault float64         `gorm:"not null" json:"probability_to_default"`
	Grade                string          `gorm:"not null" json:"grade"`
	Date                 string          `gorm:"not null" json:"date"`
	Type                 string          `gorm:"not null" json:"type"`
	UserID               string          `gorm:"not null" json:"user_id"`
	User                 user.User       `gorm:"foreignKey:UserID" json:"user"`
	CompanyID            string          `json:"company_id"`
	Company              company.Company `gorm:"foreignKey:CompanyID" json:"company"`
	CreatedAt            time.Time       `json:"-"`
	UpdatedAt            time.Time       `json:"-"`
	DeletedAt            gorm.DeletedAt  `gorm:"index" json:"-"`
}

type BulkSearchRequest struct {
	LoanNo      string `json:"loan_no"`
	Name        string `json:"name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
}

type BulkSearchResponse struct {
	TransactionID        string  `json:"transaction_id"`
	Name                 string  `json:"name"`
	PIC                  string  `json:"pic"`
	IDCardNo             string  `json:"id_card_no"`
	PhoneNo              string  `json:"phone_no"`
	LoanNo               string  `json:"loan_no"`
	ProbabilityToDefault float64 `json:"probability_to_default"`
	Grade                string  `json:"grade"`
	Type                 string  `json:"type"`
	Date                 string  `json:"date"`
}

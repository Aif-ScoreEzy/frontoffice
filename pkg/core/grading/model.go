package grading

import (
	"front-office/pkg/core/company"
	"time"

	"gorm.io/gorm"
)

type Grading struct {
	Id           string          `gorm:"primarykey" json:"id"`
	GradingLabel string          `gorm:"not null" json:"grading_label"`
	MinGrade     float64         `gorm:"not null" json:"min_grade"`
	MaxGrade     float64         `gorm:"not null" json:"max_grade"`
	CompanyId    string          `json:"company_id"`
	Company      company.Company `gorm:"foreignKey:CompanyId" json:"-"`
	CreatedAt    time.Time       `json:"-"`
	UpdatedAt    time.Time       `json:"-"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
}

type MstGrade struct {
	Id    uint    `json:"id"`
	Grade string  `json:"grade"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type CreateGradingRequest struct {
	GradingLabel string   `json:"grading_label"`
	MinGrade     *float64 `json:"min_grade"`
	MaxGrade     *float64 `json:"max_grade"`
}

type CreateGradingsRequest struct {
	CreateGradingsRequest []*CreateGradingRequest `json:"gradings"`
}

type UpdateGradingRequest struct {
	Id           string    `json:"id"`
	GradingLabel string    `json:"grading_label"`
	MinGrade     *float64  `json:"min_grade"`
	MaxGrade     *float64  `json:"max_grade"`
	IsDeleted    bool      `json:"is_deleted"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type UpdateGradingsRequest struct {
	UpdateGradingsRequest []*UpdateGradingRequest `json:"gradings"`
}

type CreateGradingNewRequest struct {
	Grade string     `json:"grade"`
	Value []*float64 `json:"value"`
}

type CreateGradingsNewRequest struct {
	CreateGradingsNewRequest []*CreateGradingNewRequest `json:"gradings"`
}

type DataGradesResponse struct {
	ApiconfigId uint       `json:"apiconfig_id"`
	CompanyId   uint       `json:"company_id"`
	BasePrice   float64    `json:"base_price"`
	Grades      []MstGrade `json:"grades"`
	AddOns      []any      `json:"addons"`
}

type AifResponse struct {
	Success bool                `json:"success"`
	Data    *DataGradesResponse `json:"data"`
	Message string              `json:"message"`
	Meta    any                 `json:"meta,omitempty"`
	Status  bool                `json:"status,omitempty"`
}

package grading

import (
	"front-office/pkg/company"
	"time"

	"gorm.io/gorm"
)

type Grading struct {
	ID           string          `gorm:"primarykey" json:"id"`
	GradingLabel string          `gorm:"not null" json:"grading_label"`
	MinGrade     float64         `gorm:"not null" json:"min_grade"`
	MaxGrade     float64         `gorm:"not null" json:"max_grade"`
	CompanyID    string          `json:"company_id"`
	Company      company.Company `gorm:"foreignKey:CompanyID" json:"-"`
	CreatedAt    time.Time       `json:"-"`
	UpdatedAt    time.Time       `json:"-"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
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
	ID           string    `json:"id"`
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

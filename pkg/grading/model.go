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
	Company      company.Company `gorm:"foreignKey:CompanyID" json:"company"`
	CreatedAt    time.Time       `json:"-"`
	UpdatedAt    time.Time       `json:"-"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
}

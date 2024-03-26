package industry

import (
	"time"

	"gorm.io/gorm"
)

type Industry struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"-"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

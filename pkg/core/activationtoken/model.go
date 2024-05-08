package activationtoken

import (
	"front-office/pkg/core/user"
	"time"

	"gorm.io/gorm"
)

type ActivationToken struct {
	ID         string         `gorm:"primarykey" json:"id"`
	Token      string         `gorm:"not null" json:"token"`
	Activation bool           `gorm:"not null;default:false" json:"activation"`
	UserID     string         `json:"user_id"`
	User       user.User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

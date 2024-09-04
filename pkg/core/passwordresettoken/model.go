package passwordresettoken

import (
	"front-office/pkg/core/user"
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	Id         uint           `gorm:"primarykey" json:"id"`
	Token      string         `gorm:"not null" json:"token"`
	Activation bool           `gorm:"not null;default:false" json:"activation"`
	UserId     string         `json:"user_id"`
	User       user.User      `gorm:"foreignKey:UserId" json:"user"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreatePasswordResetTokenRequest struct {
	Token string `json:"token"`
}

type FindTokenResponse struct {
	Message    string              `json:"message"`
	Success    bool                `json:"success"`
	Data       *PasswordResetToken `json:"data"`
	StatusCode int                 `json:"-"`
}

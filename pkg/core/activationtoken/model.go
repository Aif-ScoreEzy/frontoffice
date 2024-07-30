package activationtoken

import (
	"front-office/pkg/core/user"
	"time"

	"gorm.io/gorm"
)

type MstActivationToken struct {
	Id         string         `gorm:"primarykey" json:"id"`
	Token      string         `gorm:"not null" json:"token"`
	Activation bool           `gorm:"not null;default:false" json:"activation"`
	UserId     string         `json:"user_id"`
	User       user.User      `gorm:"foreignKey:UserId" json:"user"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateActivationTokenRequest struct {
	Token string `json:"token"`
}

type CreateActivationTokenResponse struct {
	Message    string              `json:"message"`
	Success    bool                `json:"success"`
	Data       *MstActivationToken `json:"data"`
	StatusCode int                 `json:"-"`
}

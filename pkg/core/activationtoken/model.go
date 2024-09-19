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
	MemberId   uint           `json:"member_id"`
	Member     user.MstMember `gorm:"foreignKey:MemberId" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
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

type FindTokenResponse struct {
	Message    string              `json:"message"`
	Success    bool                `json:"success"`
	Data       *MstActivationToken `json:"data"`
	StatusCode int                 `json:"-"`
}

type AifResponse struct {
	Success bool                `json:"success"`
	Data    *MstActivationToken `json:"data"`
	Message string              `json:"message"`
	Meta    any                 `json:"meta,omitempty"`
	Status  bool                `json:"status,omitempty"`
}

type Meta struct {
	Total      any `json:"total,omitempty"`
	Page       any `json:"page,omitempty"`
	TotalPages any `json:"total_pages,omitempty"`
	Visible    any `json:"visible,omitempty"`
	StartData  any `json:"start_data,omitempty"`
	EndData    any `json:"end_data,omitempty"`
	Size       any `json:"size,omitempty"`
	Message    any `json:"message,omitempty"`
}

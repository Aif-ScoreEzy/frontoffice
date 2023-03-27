package role

import "time"

type Role struct {
	ID        string    `gorm:"primarykey" json:"-"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp" json:"-"`
	DeletedAt time.Time `gorm:"index" json:"-"`
}

type RoleRequest struct {
	ID        string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type RoleResponse struct {
	ID        string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

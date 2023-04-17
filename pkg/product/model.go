package product

import "time"

type Product struct {
	ID        string    `gorm:"primarykey" json:"id"`
	Slug      string    `gorm:"not null" json:"slug"`
	Name      string    `gorm:"not null" json:"name"`
	Version   string    `json:"version"`
	Url       string    `json:"url"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `gorm:"index" json:"-"`
}

type ProductRequest struct {
	Name    string `json:"name" validate:"required~Name cannot be empty"`
	Version string `json:"version"`
	Url     string `json:"url"`
	Key     string `json:"key"`
}

type ProductResponse struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Version string `json:"version"`
	Url     string `json:"url"`
	Key     string `json:"key"`
}

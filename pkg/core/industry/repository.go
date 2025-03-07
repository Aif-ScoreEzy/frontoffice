package industry

import (
	"fmt"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

type repository struct {
	DB *gorm.DB
}

type Repository interface {
	FindOneById(industry Industry) (Industry, error)
}

func (repo *repository) FindOneById(industry Industry) (Industry, error) {
	err := repo.DB.First(&industry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return industry, fmt.Errorf("Industry with Id %s not found", industry.Id)
		}

		return industry, fmt.Errorf("Failed to find industry with Id %s: %v", industry.Id, err)
	}

	return industry, nil
}

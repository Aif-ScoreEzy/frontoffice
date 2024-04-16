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
	FindOneByID(industry Industry) (Industry, error)
}

func (repo *repository) FindOneByID(industry Industry) (Industry, error) {
	err := repo.DB.First(&industry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return industry, fmt.Errorf("Industry with ID %s not found", industry.ID)
		}

		return industry, fmt.Errorf("Failed to find industry with ID %s: %v", industry.ID, err)
	}

	return industry, nil
}

package company

import (
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

type repository struct {
	DB *gorm.DB
}

type Repository interface {
	FindOneByID(id string) (*Company, error)
	UpdateOneByID(req Company, id string) (Company, error)
}

func (repo *repository) FindOneByID(id string) (*Company, error) {
	var company *Company

	err := repo.DB.First(&company, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (repo *repository) UpdateOneByID(req Company, id string) (Company, error) {
	var company Company
	repo.DB.First(&company, "id = ?", id)

	err := repo.DB.Debug().Model(&company).Updates(req).Error
	if err != nil {
		return company, err
	}

	return company, nil
}

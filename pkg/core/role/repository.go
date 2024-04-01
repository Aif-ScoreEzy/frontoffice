package role

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
	Create(role Role) (Role, error)
	FindAll() ([]Role, error)
	FindOneByID(id string) (*Role, error)
	FindOneByName(name string) (*Role, error)
	UpdateByID(req *Role, id string) (*Role, error)
	Delete(id string) error
}

func (repo *repository) Create(role Role) (Role, error) {
	result := repo.DB.Create(&role)

	repo.DB.Preload("Permissions").First(&role, "id = ?", role.ID)

	return role, result.Error
}

func (repo *repository) FindAll() ([]Role, error) {
	var roles []Role

	result := repo.DB.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return roles, result.Error
	}

	return roles, nil
}

func (repo *repository) FindOneByID(id string) (*Role, error) {
	var role *Role

	err := repo.DB.Preload("Permissions").First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (repo *repository) FindOneByName(name string) (*Role, error) {
	var role *Role
	result := repo.DB.First(&role, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	return role, nil
}

func (repo *repository) UpdateByID(req *Role, id string) (*Role, error) {
	var role *Role

	result := repo.DB.Model(&role).
		Where("id = ?", id).Updates(req)
	if result.Error != nil {
		return nil, result.Error
	}

	return role, nil
}

func (repo *repository) Delete(id string) error {
	var role Role
	err := repo.DB.Where("id = ?", id).Delete(&role).Error

	return err
}

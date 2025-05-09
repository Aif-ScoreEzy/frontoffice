package permission

import (
	"errors"
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
	Create(permission Permission) (Permission, error)
	FindOneById(permission Permission) (Permission, error)
	FindOneByName(name string) (Permission, error)
	UpdateById(req PermissionRequest, id string) (Permission, error)
	Delete(id string) error
}

func (repo *repository) Create(permission Permission) (Permission, error) {
	result := repo.DB.Create(&permission)

	repo.DB.First(&permission, "id = ?", permission.Id)

	return permission, result.Error
}

func (repo *repository) FindOneById(permission Permission) (Permission, error) {
	err := repo.DB.First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return permission, fmt.Errorf("Permission with Id %s not found", permission.Id)
		}

		return permission, fmt.Errorf("Failed to find permission with Id %s: %v", permission.Id, err)
	}

	return permission, nil
}

func (repo *repository) FindOneByName(name string) (Permission, error) {
	var permission Permission
	result := repo.DB.First(&permission, "name = ?", name)
	if result.RowsAffected != 0 {
		return permission, errors.New("Permission with the same name already exists")
	}

	return permission, nil
}

func (repo *repository) UpdateById(req PermissionRequest, id string) (Permission, error) {
	var permission Permission
	result := repo.DB.Model(&permission).
		Where("id = ?", id).Updates(req)
	if result.Error != nil {
		return permission, result.Error
	}

	repo.DB.First(&permission, "id = ?", id)

	return permission, nil
}

func (repo *repository) Delete(id string) error {
	var permission Permission
	result := repo.DB.Where("id = ?", id).Delete(&permission)

	return result.Error
}

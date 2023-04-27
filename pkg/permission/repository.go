package permission

import (
	"errors"
	"fmt"
	"front-office/config/database"

	"gorm.io/gorm"
)

func Create(permission Permission) (Permission, error) {
	result := database.DBConn.Debug().Create(&permission)

	database.DBConn.First(&permission, "id = ?", permission.ID)

	return permission, result.Error
}

func FindOneByID(permission Permission) (Permission, error) {
	err := database.DBConn.Debug().First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return permission, fmt.Errorf("Permission with ID %s not found", permission.ID)
		}

		return permission, fmt.Errorf("Failed to find permission with ID %s: %v", permission.ID, err)
	}

	return permission, nil
}

func FindOneByName(name string) (Permission, error) {
	var permission Permission
	result := database.DBConn.First(&permission, "name = ?", name)
	if result.RowsAffected != 0 {
		return permission, errors.New("Permission with the same name already exists")
	}

	return permission, nil
}

func UpdateByID(req PermissionRequest, id string) (Permission, error) {
	var permission Permission
	result := database.DBConn.Debug().Model(&permission).
		Where("id = ?", id).Updates(req)
	if result.Error != nil {
		return permission, result.Error
	}

	database.DBConn.First(&permission, "id = ?", id)

	return permission, nil
}

func Delete(id string) error {
	var permission Permission
	result := database.DBConn.Debug().Where("id = ?", id).Delete(&permission)

	return result.Error
}

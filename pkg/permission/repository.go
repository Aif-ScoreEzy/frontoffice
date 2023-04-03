package permission

import (
	"front-office/config/database"
)

func Create(permission Permission) (Permission, error) {
	result := database.DBConn.Debug().Create(&permission)

	database.DBConn.First(&permission, "id = ?", permission.ID)

	return permission, result.Error
}

func FindOneByID(id string) (Permission, error) {
	var permission Permission
	result := database.DBConn.First(&permission, "id = ?", id)
	if result.Error != nil {
		return permission, result.Error
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

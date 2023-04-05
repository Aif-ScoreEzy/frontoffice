package role

import (
	"errors"
	"front-office/config/database"
)

func Create(role Role) (Role, error) {

	result := database.DBConn.Debug().Create(&role)

	database.DBConn.Preload("Permissions").First(&role, "id = ?", role.ID)

	return role, result.Error
}

func FindAll() ([]Role, error) {
	var roles []Role

	result := database.DBConn.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		return roles, result.Error
	}

	return roles, nil
}

func FindOneByID(id string) (Role, error) {
	var role Role

	result := database.DBConn.Preload("Permissions").First(&role, "id = ?", id)
	if result.Error != nil {
		return role, result.Error
	}

	return role, nil
}

func FindOneByName(name string) (Role, error) {
	var role Role
	result := database.DBConn.First(&role, "name = ?", name)
	if result.RowsAffected != 0 {
		return role, errors.New("Role with the same name already exists")
	}

	return role, nil
}

func UpdateByID(req Role, id string) (Role, error) {
	var role Role
	database.DBConn.Preload("Permissions").First(&role, "id = ?", id)

	result := database.DBConn.Debug().Model(&role).
		Where("id = ?", id).Updates(req)
	if result.Error != nil {
		return role, result.Error
	}

	return role, nil
}

func Delete(id string) error {
	var role Role

	result := database.DBConn.Debug().Where("id = ?", id).Delete(&role)

	return result.Error
}

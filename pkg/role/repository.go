package role

import (
	"front-office/config/database"
)

func Create(roleReq Role) (Role, error) {
	var role Role

	result := database.DBConn.Debug().Create(&roleReq)

	database.DBConn.First(&role, "id = ?", roleReq.ID)

	return role, result.Error
}

func FindOneByID(id string) (Role, error) {
	var role Role
	result := database.DBConn.First(&role, "id = ?", id)
	if result.Error != nil {
		return role, result.Error
	}

	return role, nil
}

func Update(roleReq RoleRequest, id string) (Role, error) {
	var role Role

	result := database.DBConn.Debug().Model(&role).
		Where("id = ?", id).Updates(roleReq)
	if result.Error != nil {
		return role, result.Error
	}

	database.DBConn.First(&role, "id = ?", id)

	return role, nil
}

func Delete(id string) error {
	var role Role

	result := database.DBConn.Debug().Where("id = ?", id).Delete(&role)

	return result.Error
}

package role

import (
	"front-office/config/database"
)

func Create(roleReq Role) (RoleResponse, error) {
	var role RoleResponse

	result := database.DBConn.Debug().Create(&roleReq)

	database.DBConn.Debug().Model(&Role{}).Where("id = ?", roleReq.ID).Find(&role)

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

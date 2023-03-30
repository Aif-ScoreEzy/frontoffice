package role

import (
	"front-office/config/database"
)

func Create(roleReq Role) (RoleResponse, error) {
	var role RoleResponse

	result := database.DBConn.Debug().Create(&roleReq)
	if result.RowsAffected < 1 {
		return role, result.Error
	}

	database.DBConn.Debug().Model(&Role{}).Where("id = ?", roleReq.ID).Find(&role)

	return role, nil
}

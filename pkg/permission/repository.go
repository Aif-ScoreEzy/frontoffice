package permission

import "front-office/config/database"

func Create(permission Permission) (Permission, error) {
	// var permission Permission

	result := database.DBConn.Debug().Create(&permission)

	database.DBConn.First(&permission, "id = ?", permission.ID)

	return permission, result.Error
}

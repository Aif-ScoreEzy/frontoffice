package config

import (
	"front-office/config/database"
	"front-office/pkg/role"
)

func Migrate() {
	db := database.DBConn

	db.AutoMigrate(&role.Role{})
}

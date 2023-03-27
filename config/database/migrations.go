package database

import (
	Role "front-office/pkg/role"
	User "front-office/pkg/user"
)

func Migrate() {
	db := DBConn

	db.AutoMigrate(&User.User{}, &Role.Role{})
}

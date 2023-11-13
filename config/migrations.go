package config

import (
	"front-office/config/database"
	activation_token "front-office/pkg/activation-token"
	"front-office/pkg/company"
	"front-office/pkg/industry"
	"front-office/pkg/password_reset_token"
	"front-office/pkg/permission"
	"front-office/pkg/product"
	"front-office/pkg/role"
	"front-office/pkg/user"
	"log"
)

func Migrate() {
	db := database.DBConn

	log.Println("Running Migrations")
	db.AutoMigrate(
		&role.Role{},
		&permission.Permission{},
		&user.User{},
		&activation_token.ActivationToken{},
		&password_reset_token.PasswordResetToken{},
		&company.Company{},
		&industry.Industry{},
		&product.Product{})
}

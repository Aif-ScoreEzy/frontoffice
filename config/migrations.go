package config

import (
	"front-office/config/database"
	"front-office/pkg/auth"
	"front-office/pkg/company"
	"front-office/pkg/industry"
	"front-office/pkg/permission"
	"front-office/pkg/product"
	"front-office/pkg/role"
	"front-office/pkg/user"
	"log"
)

func Migrate() {
	db := database.DBConn

	log.Println("Running Migrations")
	db.AutoMigrate(&role.Role{}, &permission.Permission{}, &user.User{}, &user.ActivationToken{}, &auth.PasswordResetToken{}, &company.Company{}, &industry.Industry{}, &product.Product{})
}

package migrate

import (
	"front-office/app/database"
	"front-office/pkg/activationtoken"
	"front-office/pkg/company"
	"front-office/pkg/genretail"
	"front-office/pkg/grading"
	"front-office/pkg/industry"
	"front-office/pkg/passwordresettoken"
	"front-office/pkg/permission"
	"front-office/pkg/product"
	"front-office/pkg/role"
	"front-office/pkg/user"
	"log"
)

func PostgreDB(dbase database.Database) {
	db := dbase.GetDB()

	log.Println("Running Migrations")
	err := db.AutoMigrate(
		&role.Role{},
		&permission.Permission{},
		&user.User{},
		&activationtoken.ActivationToken{},
		&passwordresettoken.PasswordResetToken{},
		&company.Company{},
		&industry.Industry{},
		&product.Product{},
		&grading.Grading{},
		&genretail.BulkSearch{},
	)

	if err != nil {
		log.Println("error migrate")
	}
}

package migrate

import (
	"front-office/app/database"
	"front-office/pkg/core/activationtoken"
	"front-office/pkg/core/company"
	"front-office/pkg/core/grading"
	"front-office/pkg/core/industry"
	"front-office/pkg/core/passwordresettoken"
	"front-office/pkg/core/permission"
	"front-office/pkg/core/role"
	"front-office/pkg/core/user"
	"front-office/pkg/scoreezy/genretail"
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
		&grading.Grading{},
		&genretail.BulkSearch{},
	)

	if err != nil {
		log.Println("error migrate")
	}
}

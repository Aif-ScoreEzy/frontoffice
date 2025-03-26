package migrate

import (
	"front-office/app/database"
)

func PostgreDB(dbase database.Database) {
	// db := dbase.GetDB()

	// log.Println("Running Migrations")
	// err := db.AutoMigrate(
	// 	&role.Role{},
	// 	&permission.Permission{},
	// 	&user.User{},
	// 	&activationtoken.MstActivationToken{},
	// 	&passwordresettoken.PasswordResetToken{},
	// 	&company.Company{},
	// 	&industry.Industry{},
	// 	&grading.Grading{},
	// 	&genretail.BulkSearch{},
	// 	&livestatus.Job{},
	// 	&livestatus.JobDetail{},
	// )

	// if err != nil {
	// 	log.Println("error migrate")
	// }
}

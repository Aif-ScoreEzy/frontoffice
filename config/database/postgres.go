package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func ConnectPostgres() error {
	dsn := "user=postgres password=srlyw4Ty38d5hHgs host=db.wzwsqnwncocfxlppuomq.supabase.co port=5432 dbname=postgres sslmode=disable Timezone=Asia/Jakarta"

	// dsn := "host=" + os.Getenv("DB_HOST") +
	// 	" user=" + os.Getenv("DB_USER") +
	// 	" password=" + os.Getenv("DB_PASSWORD") +
	// 	" dbname=" + os.Getenv("DB_NAME") +
	// 	" port=" + os.Getenv("DB_PORT") +
	// 	" sslmode=disable Timezone=Asia/Jakarta"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Cannot connect to database because: ", err.Error())
		return err
	}

	DBConn = database
	log.Println("🚀 Connected Successfully to the Database")
	return nil
}

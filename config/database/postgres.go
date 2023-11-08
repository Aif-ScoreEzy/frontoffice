package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DBConn *gorm.DB

func ConnectPostgres() error {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable Timezone=Asia/Jakarta"

	if strings.ToLower(os.Getenv("GCP_CLOUD_SQL")) == "true" {

		db, err := gorm.Open(postgres.New(
			postgres.Config{
				DriverName: "cloudsqlpostgres",
				DSN:        dsn,
			}),
			&gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "frontoffice.", // schema name
					SingularTable: false,
				},
			})

		if err != nil {
			log.Println("Can't Connect to DB on GCP because : ", err.Error())
			return err
		}
		DBConn = db
		fmt.Println("Success Connect to DB on GCP")

		return nil
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Cannot connect to database because: ", err.Error())
		return err
	}

	DBConn = database
	log.Println("🚀 Connected Successfully to the Database")
	return nil
}

package database

import (
	"fmt"
	"front-office/app/config"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

type postgresDb struct {
	Db *gorm.DB
}

func NewPostgresDb(cfg *config.Config) Database {
	dsn := "host=" + cfg.Db.Host +
		" user=" + cfg.Db.User +
		" password=" + cfg.Db.Password +
		" dbname=" + cfg.Db.Name +
		" port=" + cfg.Db.Port +
		" sslmode=" + cfg.Db.SSLMode +
		" Timezone=" + cfg.Db.TimeZone

	db, _ := ConnectPg(dsn, cfg)
	return &postgresDb{Db: db}
}

func ConnectPg(dsn string, cfg *config.Config) (*gorm.DB, error) {

	if strings.ToLower(cfg.Env.CloudProvider) == "gcp_cloudsql" {

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
				// Logger: logger.Default.LogMode(logger.Info)
			})

		if err != nil {
			log.Println("Can't Connect to DB on GCP because : ", err.Error())
			return db, err
		}
		// PG = db
		fmt.Println("Success Connect to DB on GCP")

		return db, err
	}

	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN: dsn,
		},
	),
		&gorm.Config{
			// Logger: logger.Default.LogMode(logger.Info),
		},
	)

	if err != nil {
		log.Println("Can't connect to database because: ", err.Error())
	}

	log.Println("🚀 Success connect to the Database")

	return db, nil
}

func (p *postgresDb) GetDB() *gorm.DB {
	return p.Db
}

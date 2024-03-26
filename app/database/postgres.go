package database

import (
	"front-office/app/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN: dsn,
		},
	),
		&gorm.Config{},
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

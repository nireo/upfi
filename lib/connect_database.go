package lib

import (
	"fmt"

	"github.com/nireo/upfi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	User string
	Port string
	Name string
	Host string
}

func ConnectToDatabase(conf *DatabaseConfig) error {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			conf.Host, conf.Port, conf.User, conf.Name),
	}), &gorm.Config{})

	if err != nil {
		return err
	}

	models.MigrateModels(db)
	SetDatabase(db)

	return nil
}

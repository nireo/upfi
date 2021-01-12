package models

import (
	"fmt"

	"github.com/nireo/upfi/lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds all the fields needed to connect to a database.
type DatabaseConfig struct {
	User string
	Port string
	Name string
	Host string
}

// ConnectToDatabase sets up a gorm connection to a database given a database config. Also this functions migrates
// all the models and sets the global database variable.
func ConnectToDatabase(conf *DatabaseConfig) error {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			conf.Host, conf.Port, conf.User, conf.Name),
	}), &gorm.Config{})

	if err != nil {
		return err
	}

	MigrateModels(db)
	lib.SetDatabase(db)

	return nil
}

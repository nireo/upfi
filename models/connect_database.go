package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nireo/upfi/lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	User string
	Port string
	Name string
	Host string
}

func SetupTestDatabase() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Could not load env file")
	}

	conf := &DatabaseConfig{
		User: os.Getenv("db_username"),
		Port: os.Getenv("db_port"),
		Host: os.Getenv("db_host"),
		Name: os.Getenv("db_name"),
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			conf.Host, conf.Port, conf.User, conf.Name),
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	lib.SetDatabase(db)

	return db, nil
}

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

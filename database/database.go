package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nireo/upfi/models"
)

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.File{})
}

func InitializeDatabase() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost user=postgre dbname=upfi")
	if err != nil {
		panic("Cannot connect to database")
	}

	defer db.Close()
	return db
}

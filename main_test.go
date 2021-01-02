package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nireo/upfi/models"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		// Stop the execution, since we need all of the environment varialbes
		log.Fatal("Could not load the env file")
	}

	// Store most of the enviroment varialbes into normal variables, so that the database connection
	// line becomes more easier to read.
	databaseConfig := &models.DatabaseConfig{
		User: os.Getenv("db_username"),
		Port: os.Getenv("db_port"),
		Host: os.Getenv("db_host"),
		Name: os.Getenv("db_name"),
	}

	// Use a library function to setup the database connection. Also migrates the models
	// and sets a database constanst in the lib package.
	if err := models.ConnectToDatabase(databaseConfig); err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	models.ResetInformation()

	os.Exit(exitCode)
}

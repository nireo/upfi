package main

import (
	"log"
	"os"
	"testing"

	"github.com/nireo/upfi/middleware"

	"github.com/joho/godotenv"
	"github.com/nireo/upfi/models"
)

// TestMain setups the database connection such that the http handlers work properly. Also disabled the http logging,
// since test output looks quite confusing with http logging in the same place.
func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		// Stop the execution, since we need all of the environment varialbes
		log.Fatal("Could not load the env file")
	}

	// Store most of the environment varialbes into normal variables, so that the database connection
	// line becomes more easier to read.
	databaseConfig := &models.DatabaseConfig{
		User: os.Getenv("db_username"),
		Port: os.Getenv("db_port"),
		Host: os.Getenv("db_host"),
		Name: os.Getenv("db_name"),
	}

	// Use a library function to setup the database connection. Also migrates the models
	// and sets a global database variable in the lib package.
	if err := models.ConnectToDatabase(databaseConfig); err != nil {
		log.Fatal(err)
	}

	// Disable http logging
	middleware.SetHTTPLogging(false)

	// Run all the tests
	exitCode := m.Run()

	// Remove all the information created during the testing. NOTE: Also removed all of the files and database data
	// from every user.
	models.ResetInformation()
	os.Exit(exitCode)
}

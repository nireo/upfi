package main

import (
	"github.com/nireo/upfi/jsonapi"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/templateapi"
)

func main() {
	// Load all of the environment variables listed in the .env file, in the project root directory
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

	// Use the optimized version of the api, which uses the fasthttp package to improve performance
	// Is its own function, since before there was a older implementation which used net/http.
	serverPort := os.Getenv("port")

	if len(os.Args) == 2 && os.Args[1] == "reset_information" {
		models.ResetInformation()
		return
	}

	if len(os.Args) == 2 && os.Args[1] == "json" {
		jsonapi.RunJSONApi(serverPort)
	} else {
		templateapi.SetupTemplateApi(serverPort)
	}
}

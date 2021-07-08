package main

import (
	"log"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/web"

	"github.com/joho/godotenv"
	"github.com/nireo/upfi/models"
)

func main() {
	// Load all of the environment variables listed in the .env file, in the project root directory
	if err := godotenv.Load(); err != nil {
		// Stop the execution, since we need all of the environment varialbes
		log.Fatal("Could not load the env file")
	}

	filesPath := lib.AddRootToPath("files")
	if _, err := os.Stat(filesPath); os.IsNotExist(err) {
		if err := os.Mkdir(filesPath, 0755); err != nil {
			log.Fatalf("could not create file directory: %s", err)
		}
	}

	// The temp directory is used to hold the decrypted files for some time
	// until they can properly be sent to the user.
	tempPath := lib.AddRootToPath("temp")
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		if err := os.Mkdir(tempPath, 0755); err != nil {
			log.Fatalf("could not create temp directory: %s", err)
		}
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
	// and sets a global database varialbe in the lib package.
	if err := models.ConnectToDatabase(databaseConfig); err != nil {
		log.Fatal(err)
	}

	// Use the optimized version of the api, which uses the fasthttp package to improve performance
	// Is its own function, since before there was a older implementation which used net/http.
	serverPort := os.Getenv("port")

	web.StartServer(serverPort)
}

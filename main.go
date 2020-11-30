package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/middleware"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/optimized_api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nireo/upfi/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load all of the environment variables listed in the .env file, in the project root directory
	if err := godotenv.Load(); err != nil {
		// Stop the execution, since we need all of the environment varialbes
		log.Fatal("Could not load the env file")
	}

	// Store most of the enviroment varialbes into normal variables, so that the database connection
	// line becomes more easier to read.
	user := os.Getenv("db_username")
	dbPort := os.Getenv("db_port")
	host := os.Getenv("db_host")
	dbName := os.Getenv("db_name")
	serverPort := os.Getenv("port")

	// Connect to the PostgreSQL database using gorm; which returns a pointer to the database, which we
	// store in a utility file.
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, dbPort, user, dbName),
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Execute some setup functions which we need
	models.MigrateModels(db)
	lib.SetDatabase(db)

	// Use the optimized version of the api, which uses the fasthttp package to improve performance
	// Is its own function, since before there was a older implementation which used net/http.
	optimized_api.SetupOptimizedApi(serverPort)
}

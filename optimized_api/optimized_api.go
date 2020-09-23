package optimized_api

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/jinzhu/gorm"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

func SetupOptimizedApi() {
	router := fasthttprouter.New()

	// Load database
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=upfi sslmode=disable")
	if err != nil {
		panic(err)
	}
	models.MigrateModels(db)
	defer db.Close()
	lib.SetDatabase(db)

	// setup routes
	router.POST("/upload", UploadFile)
	router.POST("/register", Register)

	// start the http server
	if err := fasthttp.ListenAndServe("localhost:8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe %s", err)
	}
}

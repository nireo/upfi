package optimized_api

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/jinzhu/gorm"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/middleware"
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
	router.POST("/upload", middleware.CheckAuthentication(UploadFile))
	router.GET("/upload", middleware.CheckAuthentication(ServeUploadPage))

	router.POST("/register", Register)
	router.POST("/login", Login)
	router.GET("/login", ServeLoginPage)
	router.GET("/register", ServeRegisterPage)
	router.GET("/file/:file", middleware.CheckAuthentication(GetSingleFile))
	router.GET("/files", middleware.CheckAuthentication(GetUserFiles))
	router.PATCH("/file/:file", middleware.CheckAuthentication(UpdateFile))
	router.DELETE("/file/:file", middleware.CheckAuthentication(DeleteFile))
	router.GET("/settings", middleware.CheckAuthentication(ServeSettingsPage))
	router.POST("/settings", middleware.CheckAuthentication(HandleSettingChange))
	router.DELETE("/remove", middleware.CheckAuthentication(DeleteUser))
	router.PATCH("/password", middleware.CheckAuthentication(UpdatePassword))

	// start the http server
	if err := fasthttp.ListenAndServe("localhost:8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe %s", err)
	}
}

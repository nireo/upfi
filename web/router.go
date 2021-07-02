package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nireo/upfi/middleware"
)

func StartServer(port string) {
	router := httprouter.New()

	// misc
	router.GET("/", ServeHomePage)

	// auth
	router.GET("/login", ServeLoginPage)
	router.GET("/register", ServeRegisterPage)
	router.POST("/login", Login)
	router.POST("/register", Register)

	// files
	router.GET("/file", middleware.CheckToken(GetSingleFile))
	router.GET("/upload", middleware.CheckToken(ServeUploadPage))
	router.POST("/upload", middleware.CheckToken(UploadFile))
	router.GET("/files", middleware.CheckToken(GetUserFiles))
	router.PATCH("/file", middleware.CheckToken(UpdateFile))
	router.POST("/download", middleware.CheckToken(DownloadFile))
	router.DELETE("/shared", middleware.CheckToken(DeleteSharedContract))
	router.GET("/shared_by", middleware.CheckToken(GetSharedByUser))
	router.GET("/shared_to", middleware.CheckToken(GetSharedToUser))
	router.GET("/shared", middleware.CheckToken(ServeCreateSharedPage))
	router.POST("/shared", middleware.CheckToken(CreateSharedFile))

	// user
	router.DELETE("/remove", middleware.CheckToken(DeleteUser))
	router.PATCH("/password", middleware.CheckToken(UpdatePassword))
	router.GET("/settings", middleware.CheckToken(ServeSettingsPage))
	router.POST("/settings", middleware.CheckToken(HandleSettingChange))

	log.Fatal(http.ListenAndServe(":"+port, router))
}

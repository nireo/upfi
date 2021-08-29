package web

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/julienschmidt/httprouter"
	"github.com/nireo/upfi/middleware"
)

func StartServer(port string) {
	router := httprouter.New()

	// misc
	router.GET("/", middleware.SecureHeaders(ServeHomePage))

	// auth
	router.GET("/login", middleware.SecureHeaders(ServeLoginPage))
	router.GET("/register", middleware.SecureHeaders(ServeRegisterPage))
	router.POST("/login", middleware.SecureHeaders(Login))
	router.POST("/register", middleware.SecureHeaders(ServeRegisterPage))

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
	router.GET("/share", middleware.CheckToken(ServeCreateSharedPage))
	router.POST("/share", middleware.CheckToken(CreateSharedFile))

	// user
	router.DELETE("/remove", middleware.CheckToken(DeleteUser))
	router.PATCH("/password", middleware.CheckToken(UpdatePassword))
	router.GET("/settings", middleware.CheckToken(ServeSettingsPage))
	router.POST("/settings", middleware.CheckToken(HandleSettingChange))

	csrfSecret := os.Getenv("csrfkey")

	CSRF := csrf.Protect([]byte(csrfSecret), nil)
	log.Fatal(http.ListenAndServe(":"+port, CSRF(router)))
}

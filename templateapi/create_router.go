package templateapi

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/upfi/middleware"
)

// CreateRouter puts all the routes together and returns the fasthttp handler.
func CreateRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.POST("/register", middleware.TinyLogger(Register))
	router.POST("/login", middleware.TinyLogger(Login))
	router.GET("/login", middleware.TinyLogger(ServeLoginPage))
	router.GET("/register", middleware.TinyLogger(ServeRegisterPage))
	router.GET("/", middleware.TinyLogger(ServeHomePage))

	router.POST("/upload", middleware.CheckAuthentication(UploadFile))
	router.GET("/upload", middleware.CheckAuthentication(ServeUploadPage))
	router.GET("/file/:file", middleware.CheckAuthentication(GetSingleFile))
	router.GET("/files", middleware.CheckAuthentication(GetUserFiles))
	router.PATCH("/file/:file", middleware.CheckAuthentication(UpdateFile))
	router.DELETE("/file/:file", middleware.CheckAuthentication(DeleteFile))
	router.GET("/settings", middleware.CheckAuthentication(ServeSettingsPage))
	router.POST("/settings", middleware.CheckAuthentication(HandleSettingChange))
	router.DELETE("/remove", middleware.CheckAuthentication(DeleteUser))
	router.PATCH("/password", middleware.CheckAuthentication(UpdatePassword))

	return router
}

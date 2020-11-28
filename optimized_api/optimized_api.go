package optimized_api

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/upfi/middleware"
	"github.com/valyala/fasthttp"
)

func SetupOptimizedApi(port string) {
	router := fasthttprouter.New()

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
	router.GET("/paste", middleware.CheckAuthentication(ServeCreatePage))
	router.POST("/paste", middleware.CheckAuthentication(CreatePaste))
	router.GET("/pastes", middleware.CheckAuthentication(DeletePaste))
	router.PATCH("/paste", middleware.CheckAuthentication(UpdatePastePrivacy))

	// Start a HTTP server listening on the port from the environment variable
	if err := fasthttp.ListenAndServe(fmt.Sprintf("localhost:%s", port), router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe %s", err)
	}
}

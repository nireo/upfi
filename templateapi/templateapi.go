package templateapi

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/upfi/middleware"
	"github.com/valyala/fasthttp"
)

// SetupTemplateAPI starts a HTTP listener on a given port with all the routes in the program.
func SetupTemplateAPI(port string) {
	// Create a new instance of a router. (Not really necessary, because we can route ourselves, but
	// this way the implementation is cleaner.)
	router := fasthttprouter.New()

	// Setup all the endpoints and add the CheckAuthentication middleware to protected routes,
	// in which we need the user's username.
	router.POST("/register", middleware.TinyLogger(Register))
	router.POST("/login", middleware.TinyLogger(Login))
	router.GET("/login", middleware.TinyLogger(ServeLoginPage))
	router.GET("/register", middleware.TinyLogger(ServeRegisterPage))
	router.GET("/", middleware.TinyLogger(ServeHomePage))

	router.POST("/upload", middleware.CheckAuthentication(UploadFile))
	// router.GET("/upload", middleware.CheckAuthentication(ServeUploadPage))
	router.GET("/upload", ServeUploadPage)

	router.POST("/download/:file", middleware.CheckAuthentication(DownloadFile))
	router.GET("/file/:file", middleware.CheckAuthentication(GetSingleFile))
	router.GET("/files", middleware.CheckAuthentication(GetUserFiles))
	router.PATCH("/file/:file", middleware.CheckAuthentication(UpdateFile))
	router.DELETE("/file/:file", middleware.CheckAuthentication(DeleteFile))
	router.GET("/settings", middleware.CheckAuthentication(ServeSettingsPage))
	router.POST("/settings", middleware.CheckAuthentication(HandleSettingChange))
	router.DELETE("/remove", middleware.CheckAuthentication(DeleteUser))
	router.POST("/password", middleware.CheckAuthentication(UpdatePassword))

	log.Printf("current running on port: %s", port)

	// Start a HTTP server listening on the port from the environment variable
	if err := fasthttp.ListenAndServe(fmt.Sprintf("localhost:%s", port), router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe %s", err)
	}
}

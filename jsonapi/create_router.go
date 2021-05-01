package jsonapi

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/upfi/middleware"
)

// CreateJSONRouter creates a fasthttp router which contains all the json handlers. Some of the json handlers
// are pretty much the same as the templateapi only exception being that data is given in json form.
func CreateJSONRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.POST("/api/register", Register)
	router.POST("/api/login", Login)

	router.POST("/api/file", middleware.CheckAuthentication(UploadFile))
	router.GET("/api/file/:file", middleware.CheckAuthentication(GetSingleFile))
	router.GET("/api/download/:file", middleware.CheckAuthentication(DownloadFile))

	router.DELETE("/api/file/:file", middleware.CheckAuthentication(DeleteFile))
	router.GET("/api/files", middleware.CheckAuthentication(GetUserFiles))

	router.PUT("/api/file/:file", middleware.CheckAuthentication(UpdateFile))

	router.PUT("/api/password", middleware.CheckAuthentication(UpdatePassword))
	router.PUT("/api/settings", middleware.CheckAuthentication(HandleSettingsChange))

	router.GET("/api/me", middleware.CheckAuthentication(WhoAmI))

	return router
}

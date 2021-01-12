package jsonapi

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/upfi/middleware"
	"github.com/nireo/upfi/templateapi"
)

// CreateJSONRouter creates a fasthttp router which contains all the json handlers. Some of the json handlers
// are pretty much the same as the templateapi only exception being that data is given in json form.
func CreateJSONRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.POST("/api/register", Register)
	router.POST("/api/login", Login)

	router.POST("/api/file", middleware.CheckAuthentication(UploadFile))
	router.GET("/api/file", middleware.CheckAuthentication(GetSingleFile))

	router.DELETE("/api/file", middleware.CheckAuthentication(DeleteFile))
	router.GET("/api/files", middleware.CheckAuthentication(GetUserFiles))

	router.PATCH("/api/file", middleware.CheckAuthentication(UpdateFile))

	router.DELETE("/api/username", middleware.CheckAuthentication(templateapi.DeleteUser))
	router.PATCH("/api/password", middleware.CheckAuthentication(UpdatePassword))
	router.PATCH("/api/settings", middleware.CheckAuthentication(HandleSettingsChange))

	router.GET("/api/me", middleware.CheckAuthentication(WhoAmI))

	return router
}

package json_api

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/upfi/middleware"
	"github.com/nireo/upfi/optimized_api"
)

func CreateJSONRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.POST("/api/register", Register)
	router.POST("/api/login", Login)

	// Use the optimized api file upload function since we can't really upload files using json.
	// So we just reuse the old since no need to copy-and-paste.
	router.POST("/api/upload", middleware.CheckAuthentication(optimized_api.UploadFile))

	router.DELETE("/api/username", middleware.CheckAuthentication(optimized_api.DeleteUser))
	router.PATCH("/api/password", middleware.CheckAuthentication(UpdatePassword))
	router.PATCH("/api/settings", middleware.CheckAuthentication(HandleSettingsChange))

	router.GET("/api/me", middleware.CheckAuthentication(WhoAmI))

	return router
}

package json_api

import "github.com/buaazp/fasthttprouter"

func CreateJSONRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.POST("/api/register", Register)

	return router
}

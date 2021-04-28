package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func StartService(port string) {
	router := httprouter.New()

	// auth
	router.GET("/login", ServeLoginPage)
	router.GET("/register", ServeRegisterPage)
	router.POST("/login", Login)
	router.POST("/register", Register)

	// files
	router.GET("/file", GetSingleFile)
	router.GET("/upload", UploadFile)
	router.POST("/upload", UploadFile)
	router.GET("/files", GetUserFiles)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

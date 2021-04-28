package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func StartService(port string) {
	router := httprouter.New()
	router.GET("/login", ServeLoginPage)
	router.GET("/register", ServeRegisterPage)

	router.POST("/login", Login)
	router.POST("/register", Register)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

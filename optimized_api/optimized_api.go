package optimized_api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/buaazp/fasthttprouter"
	"github.com/jinzhu/gorm"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/server"
	"github.com/valyala/fasthttp"
)

func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	// prompt the user
	fmt.Fprintf(ctx, "Cannot: '%s' route: '%s'", ctx.Method(), ctx.RequestURI())
	ctx.SetContentType("text/plain; charset=utf-8")

	// set not found status
	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
}

func Register(ctx *fasthttp.RequestCtx) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	if len(form.Value["username"]) == 0 || len(form.Value["password"][0]) == 0 {
		ctx.Error("Both username and password fields must be added", fasthttp.StatusNotFound)
		return
	}

	username := form.Value["username"][0]
	password := form.Value["password"][0]

	// check for conflicting users
	_, err = models.FindOneUser(&models.User{Username: username})
	if err == nil {
		ctx.Error("User already exists with that username", fasthttp.StatusConflict)
		return
	}

	hash, err := server.HashPassword(password)
	if err != nil {
		ctx.Error("Internal server error", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Username: username,
		Password: hash,
		UUID:     lib.GenerateUUID(),
	}

	// create a directory which stores all of the user's files
	err = os.Mkdir("./files/"+newUser.UUID, os.ModePerm)
	if err != nil {
		ctx.Error("Failed user directory creation", fasthttp.StatusInternalServerError)
		return
	}

	db := lib.GetDatabase()
	db.NewRecord(newUser)
	db.Create(&newUser)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func Login(ctx *fasthttp.RequestCtx) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	if len(form.Value["username"]) == 0 || len(form.Value["password"][0]) == 0 {
		ctx.Error("Both username and password fields must be added", fasthttp.StatusNotFound)
		return
	}

	username := form.Value["username"][0]
	password := form.Value["password"][0]

	// validate the input
	user, err := models.FindOneUser(&models.User{Username: username})
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		return
	}

	if !server.CheckPasswordHash(password, user.Password) {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusForbidden), fasthttp.StatusForbidden)
		return
	}

	// create token string
	token, err := lib.CreateToken(user.Username)
	if err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	// create token cookie
	var cookie fasthttp.Cookie
	cookie.SetKey("token")
	cookie.SetValue(token)

	ctx.Response.Header.SetCookie(&cookie)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
}

func UploadFile(ctx *fasthttp.RequestCtx) {
	// get file
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error("File could not be parsed", fasthttp.StatusInternalServerError)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	db := lib.GetDatabase()
	// find user
	user, err := models.FindOneUser(&models.User{Username: form.Value["username"][0]})
	if err != nil {
		ctx.Error("User not found", fasthttp.StatusNotFound)
		return
	}

	newFileEntry := &models.File{
		Filename:    header.Filename,
		UUID:        lib.GenerateUUID(),
		Description: form.Value["description"][0],
		Size:        header.Size,
		UserID:      user.ID,
		Extension:   filepath.Ext(header.Filename),
	}

	fileDirectory := fmt.Sprintf("./files/%s/%s%s", user.UUID, newFileEntry.UUID, newFileEntry.Extension)
	if err := fasthttp.SaveMultipartFile(header, fileDirectory); err != nil {
		ctx.Error("File could not be stored", fasthttp.StatusInternalServerError)
		return
	}

	db.NewRecord(newFileEntry)
	db.Create(newFileEntry)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func SetupOptimizedApi() {
	router := fasthttprouter.New()

	// Load database
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=upfi sslmode=disable")
	if err != nil {
		panic(err)
	}
	models.MigrateModels(db)
	defer db.Close()
	lib.SetDatabase(db)

	// setup routes
	router.POST("/upload", UploadFile)
	router.POST("/register", Register)

	// start the http server
	if err := fasthttp.ListenAndServe("localhost:8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe %s", err)
	}
}

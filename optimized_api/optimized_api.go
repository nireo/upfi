package optimized_api

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/buaazp/fasthttprouter"
	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/valyala/fasthttp"
)

func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	// prompt the user
	fmt.Fprintf(ctx, "Cannot: '%s' route: '%s'", ctx.Method(), ctx.RequestURI())
	ctx.SetContentType("text/plain; charset=utf-8")

	// set not found status
	ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
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
		ctx.Error("Content type need to multipart", fasthttp.StatusBadRequest)
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
	router.POST("/upload", UploadFile)

	if err := fasthttp.ListenAndServe("localhost:8080", router.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe %s", err)
	}
}

package optimized_api

import (
	"html/template"
	"os"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
	"github.com/nireo/upfi/server"
	"github.com/valyala/fasthttp"
)

func ServeRegisterPage(ctx *fasthttp.RequestCtx) {
	tmpl := template.Must(template.ParseFiles("./static/register.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func ServeLoginPage(ctx *fasthttp.RequestCtx) {
	tmpl := template.Must(template.ParseFiles("./static/login.html"))
	if err := tmpl.Execute(ctx, nil); err != nil {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}
}

func Register(ctx *fasthttp.RequestCtx) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("Content type needs to be multipart", fasthttp.StatusBadRequest)
		return
	}

	if len(form.Value["username"]) == 0 || len(form.Value["password"][0]) == 0 {
		ctx.Error("Both username and password fields must be added", fasthttp.StatusBadRequest)
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
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
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

package templates

import (
	"embed"
	"html/template"
	"io"

	"github.com/nireo/upfi/lib"
	"github.com/nireo/upfi/models"
)

// This file contains definitions for all of the html pages. Utilizing the go embed's
// to easily contain all of the html files in the compiled binary. This makes the task
// of deploying the application easier.

// For some pages the Authenticated parameter is not needed, since it is obvious that the user
// is authenticated if he/she can access the site.

//go:embed *
var files embed.FS

// Define all the template pages such that using all the templates is easier.
var (
	home = parse("home.html")

	filesPage  = parse("files_template.html")
	fileSingle = parse("single_file_template.html")
	upload     = parse("upload.html")

	settings = parse("settings_template.html")

	login    = parse("login.html")
	register = parse("register.html")

	errorPage = parse("error_page.html")

	successPage = parse("success_page.html")
)

// FilesParams contains all of the parameters to the files page.
type FilesParams struct {
	Title         string
	Files         []models.File
	Authenticated bool
}

// Files renders the files template file
func Files(w io.Writer, params FilesParams) error {
	return filesPage.Execute(w, params)
}

// SingleFileParams contains all of the parameters to the single file page.
type SingleFileParams struct {
	Title         string
	File          models.File
	Authenticated bool
}

// SingleFile renders the single file template file
func SingleFile(w io.Writer, params SingleFileParams) error {
	return fileSingle.Execute(w, params)
}

type HomeParams struct {
	Title         string
	Authenticated bool
}

// Register renders the register template file
func Home(w io.Writer, params HomeParams) error {
	return home.Execute(w, params)
}

// SettingsParams contains all of the parameters to the settings page.
type SettingsParams struct {
	Title         string
	User          *models.User
	Authenticated bool
}

// Settings renders the settings template file
func Settings(w io.Writer, params SettingsParams) error {
	return settings.Execute(w, params)
}

// LoginParams contains parameters for the login page
type LoginParams struct {
	Authenticated bool
	Title         string
}

// Login renders the login template file
func Login(w io.Writer, params LoginParams) error {
	return login.Execute(w, params)
}

// RegisterParams contains parameters for the register page
type RegisterParams struct {
	Authenticated bool
	Title         string
}

// Register renders the register template file
func Register(w io.Writer, params RegisterParams) error {
	return register.Execute(w, params)
}

// UploadParams contains all of the parameters to the upload page
type UploadParams struct {
	Authenticated bool
	Title         string
}

// Upload renders the upload template file
func Upload(w io.Writer, params UploadParams) error {
	return upload.Execute(w, params)
}

// ErrorParams contains all of the parameters to the error pages
type ErrorParams struct {
	Title         string
	Authenticated bool
	Error         lib.ErrorPageContent
}

// ErrorPage contains the renderer for the error pages
func ErrorPage(w io.Writer, params ErrorParams) error {
	return errorPage.Execute(w, params)
}

// SuccessPage exists, since we cannot redirect the user's request from POST -> GET, so we render a
// success page in which the user clicks the button that actually redirects them to the link. It also
// makes the experience more responsive.
type SuccessPage struct {
	Title        string
	Description  string
	RedirectPath string
}

// Success renders the success_page.html template with the given success page parameters.
func Success(w io.Writer, params SuccessPage) error {
	return successPage.Execute(w, params)
}

// parse takes in a file path and parses the embedded template files for the file and returns a
// template pointer *template.Template. Mostly used to make defining pages more elegant and clear.
func parse(file string) *template.Template {
	return template.Must(template.New("layout.html").ParseFS(files, "layout.html", file))
}

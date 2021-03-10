package templates

import (
	"embed"
	"html/template"
	"io"

	"github.com/nireo/upfi/models"
)

// This file contains definitions for all of the html pages. Utilizing the go embed's
// to easily contain all of the html files in the compiled binary. This makes the task
// of deploying the application easier.

//go:embed templates/*
var files embed.FS

// Define all the template pages such that using all the templates is easier.
var (
	dashboard = parse("dashboard.html")
	home      = parse("home.html")

	filesPage  = parse("files.html")
	fileSingle = parse("single_file.html")

	settings = parse("settings.html")

	login    = parse("login.html")
	register = parse("register.html")
)

// DashboardParams contains all of the parameters to the dashboard page.
type DashboardParams struct {
	Title string
}

// Dashboard renders the dashboard template file
func Dashboard(w io.Writer) error {
	return dashboard.Execute(w, nil)
}

// FilesParams contains all of the parameters to the files page.
type FilesParams struct {
	Title string
	Files []models.File
}

// Files renders the files template file
func Files(w io.Writer, params FilesParams) error {
	return filesPage.Execute(w, params)
}

// SingleFileParams contains all of the parameters to the single file page.
type SingleFileParams struct {
	Title string
	File  models.File
}

// SingleFile renders the single file template file
func SingleFile(w io.Writer, params SingleFileParams) error {
	return fileSingle.Execute(w, params)
}

// Register renders the register template file
func Home(w io.Writer) error {
	return home.Execute(w, nil)
}

// SettingsParams contains all of the parameters to the settings page.
type SettingsParams struct {
	Title string
	User  models.User
}

// Settings renders the settings template file
func Settings(w io.Writer, params SettingsParams) error {
	return settings.Execute(w, params)
}

// Login renders the login template file
func Login(w io.Writer) error {
	return login.Execute(w, nil)
}

// Register renders the register template file
func Register(w io.Writer) error {
	return register.Execute(w, nil)
}

// parse takes in a file path and parses the embedded template files for the file and returns a
// template pointer *template.Template. Mostly used to make defining pages more elegant and clear.
func parse(file string) *template.Template {
	return template.Must(template.New("layout.html").ParseFS(files, "layout.html", file))
}

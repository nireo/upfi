package lib

import (
	"net/http"
	"text/template"
)

func HttpNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/not_found_template.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

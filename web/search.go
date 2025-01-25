package web

import (
	"github.com/jeffscottbrown/applemusic/model"
	"github.com/jeffscottbrown/applemusic/search"
	"html/template"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/search.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	bandName := r.FormValue("band_name")
	results, errorMessage := search.SearchApple(bandName)

	tmpl.Execute(w, struct {
		Success    bool
		Error      string
		Results    model.SearchResult
		SearchTerm string
	}{errorMessage == "", errorMessage, results, bandName})
}

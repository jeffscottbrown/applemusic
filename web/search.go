package web

import (
	"github.com/jeffscottbrown/applemusic/controllers"
	"github.com/jeffscottbrown/applemusic/model"
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
	results := controllers.SearchApple(bandName)

	tmpl.Execute(w, struct {
		Success    bool
		Results    model.SearchResult
		SearchTerm string
	}{true, results["searchResults"].(model.SearchResult), bandName})

}

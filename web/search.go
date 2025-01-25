package web

import (
	"fmt"
	"github.com/jeffscottbrown/applemusic/commit"
	"github.com/jeffscottbrown/applemusic/model"
	"github.com/jeffscottbrown/applemusic/search"
	"html/template"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/search.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, struct {
			CommitHash string
			BuildTime  string
			Error      bool
			Success    bool
		}{
			commit.Hash,
			commit.BuildTime,
			false,
			false,
		})
		return
	}

	bandName := r.FormValue("band_name")
	results, errorMessage := search.SearchApple(bandName)

	jsonUrl := fmt.Sprintf("/search/%s", url.QueryEscape(bandName))

	tmpl.Execute(w, struct {
		Success    bool
		Error      string
		Results    model.SearchResult
		SearchTerm string
		JsonUrl    string
		CommitHash string
		BuildTime  string
	}{errorMessage == "",
		errorMessage,
		results,
		bandName,
		jsonUrl,
		commit.Hash,
		commit.BuildTime,
	})
}

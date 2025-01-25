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
	searchModel := SearchModel{CommitHash: commit.Hash,
		BuildTime: commit.BuildTime}

	if r.Method != http.MethodPost {

		tmpl.Execute(w, searchModel)
		return
	}

	bandName := r.FormValue("band_name")
	results, errorMessage := search.SearchApple(bandName)

	searchModel.SearchTerm = bandName
	searchModel.Error = errorMessage
	searchModel.Results = results
	searchModel.JsonUrl = fmt.Sprintf("/search/%s", url.QueryEscape(bandName))

	tmpl.Execute(w, searchModel)
}

type SearchModel struct {
	BuildTime  string
	CommitHash string
	Error      string
	JsonUrl    string
	Results    model.SearchResult
	SearchTerm string
}

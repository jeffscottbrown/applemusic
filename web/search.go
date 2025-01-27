package web

import (
	"fmt"
	"github.com/jeffscottbrown/applemusic/commit"
	"github.com/jeffscottbrown/applemusic/constants"
	"github.com/jeffscottbrown/applemusic/model"
	"github.com/jeffscottbrown/applemusic/search"
	"github.com/markbates/goth/gothic"
	"html/template"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/search.html"))
	searchModel := createModel(r)

	if r.Method != http.MethodPost {

		tmpl.Execute(w, searchModel)
		return
	}

	bandName := r.FormValue("band_name")
	limit := r.FormValue("limit")
	results, errorMessage := search.SearchApple(bandName, limit)

	searchModel.SearchTerm = bandName
	searchModel.Error = errorMessage
	searchModel.Results = results
	searchModel.JsonUrl = fmt.Sprintf("/search/%s", url.QueryEscape(bandName))

	tmpl.Execute(w, searchModel)
}

func createModel(r *http.Request) SearchModel {
	searchModel := SearchModel{CommitHash: commit.Hash,
		BuildTime:  commit.BuildTime,
		GitHubRepo: constants.GitHubRepo}

	authenticatedUserName, _ := gothic.GetFromSession("authenticatedUser", r)
	searchModel.AuthenticatedUserName = authenticatedUserName
	return searchModel
}

type SearchModel struct {
	BuildTime             string
	CommitHash            string
	Error                 string
	JsonUrl               string
	Results               model.SearchResult
	SearchTerm            string
	GitHubRepo            string
	AuthenticatedUserName string
}

package search

import (
	"encoding/json"
	"github.com/jeffscottbrown/applemusic/model"
	"log/slog"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.PathValue("term")
	data, err := SearchApple(searchTerm)
	w.Header().Add("Content-Type", "application/json")
	if err == "" {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err})
	}
}

func SearchApple(searchTerm string) (model.SearchResult, string) {
	fullURL := createSearchUrl(searchTerm)

	slog.Debug("Querying Apple API", "url", fullURL)

	resp, err := http.Get(fullURL)

	var errorMessage string
	var result model.SearchResult
	if err != nil {
		errorMessage = "Failed to fetch data"
	} else {
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			errorMessage = "Failed to parse JSON"
		}
	}
	return result, errorMessage
}

func createSearchUrl(searchTerm string) string {
	apiURL := "https://itunes.apple.com/search"

	params := createRequestParameters(searchTerm)

	fullURL := apiURL + "?" + params.Encode()
	return fullURL
}

func createRequestParameters(searchTerm string) url.Values {
	params := url.Values{}
	params.Add("term", searchTerm)
	params.Add("media", "music")
	params.Add("entity", "album")
	return params
}

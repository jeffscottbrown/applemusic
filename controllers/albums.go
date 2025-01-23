package controllers

import (
	"encoding/json"
	"github.com/jeffscottbrown/applemusic/model"
	"log/slog"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.PathValue("term")
	data := SearchApple(searchTerm)
	json.NewEncoder(w).Encode(data)
}

func SearchApple(searchTerm string) map[string]interface{} {
	apiURL := "https://itunes.apple.com/search"

	params := url.Values{}
	params.Add("term", searchTerm)
	params.Add("media", "music")
	params.Add("entity", "album")

	fullURL := apiURL + "?" + params.Encode()

	slog.Debug("Querying Apple API", "url", fullURL)

	resp, err := http.Get(fullURL)

	data := make(map[string]interface{})

	if err != nil {
		data["error"] = "Failed to fetch data"
	} else {
		defer resp.Body.Close()

		var result model.SearchResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			data["error"] = "Failed to parse JSON"
		} else {
			data["searchResults"] = result
		}
	}
	return data
}

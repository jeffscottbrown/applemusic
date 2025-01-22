package controllers

import (
	"encoding/json"
	"github.com/jeffscottbrown/applemusic/model"
	"net/http"
	"net/url"
)

func Search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.PathValue("term")
	apiURL := "https://itunes.apple.com/search"

	params := url.Values{}
	params.Add("term", searchTerm)
	params.Add("media", "music")
	params.Add("entity", "album")

	fullURL := apiURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch data"})
		return
	}
	defer resp.Body.Close()

	var result model.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse JSON"})
		return
	}

	json.NewEncoder(w).Encode(result)
}

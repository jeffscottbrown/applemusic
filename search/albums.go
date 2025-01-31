package search

import (
	"encoding/json"
	"fmt"
	"github.com/jeffscottbrown/applemusic/constants"
	"github.com/jeffscottbrown/applemusic/model"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

var searchCache = cache.New(5*time.Minute, 10*time.Minute)

func Search(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.PathValue("term")
	data, err := SearchApple(searchTerm, "25")
	w.Header().Add("Content-Type", "application/json")
	if err == "" {
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err})
	}
}

func SearchApple(searchTerm string, limit string) (model.SearchResult, string) {
	cachKey := fmt.Sprintf("%s-%s", searchTerm, limit)
	if cachedData, found := searchCache.Get(cachKey); found {
		slog.Debug("Cache hit", "searchTerm", searchTerm)
		return cachedData.(model.SearchResult), ""
	}

	fullURL := createSearchUrl(searchTerm, limit)

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
		searchCache.Set(cachKey, result, cache.DefaultExpiration)
	}
	return result, errorMessage
}

func createSearchUrl(searchTerm string, limit string) string {
	params := createRequestParameters(searchTerm, limit)

	fullURL := constants.AppleMusicAPI + "?" + params.Encode()
	return fullURL
}

func createRequestParameters(searchTerm string, limit string) url.Values {
	params := url.Values{}
	params.Add("term", searchTerm)
	params.Add("media", "music")
	params.Add("entity", "album")
	params.Add("limit", limit)
	return params
}

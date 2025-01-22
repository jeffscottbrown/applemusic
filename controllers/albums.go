package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/jeffscottbrown/applemusic/model"
)

func Search(c *gin.Context) {
	searchTerm := c.Param("term")
	apiURL := "https://itunes.apple.com/search"

	params := url.Values{}
	params.Add("term", searchTerm)
	params.Add("media", "music")
	params.Add("entity", "album")

	fullURL := apiURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer resp.Body.Close()

	var result model.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
		return
	}

	c.JSON(http.StatusOK, result)
}

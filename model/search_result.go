package model

type SearchResult struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		ArtistName string `json:"artistName"`
		AlbumTitle string `json:"collectionName"`
	} `json:"results"`
}

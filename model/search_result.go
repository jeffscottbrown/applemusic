package model

type SearchResult struct {
	Results []struct {
		ArtistName string `json:"artistName"`
		AlbumTitle string `json:"collectionName"`
	} `json:"results"`
}

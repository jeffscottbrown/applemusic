package server

import (
	"github.com/jeffscottbrown/applemusic/search"
	"github.com/jeffscottbrown/applemusic/web"
	"net/http"
)

func RunServer() {
	router := http.NewServeMux()
	router.HandleFunc("/search/{term}", search.Search)
	router.HandleFunc("/", web.Search)

	fs := http.FileServer(http.Dir("assets/"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}

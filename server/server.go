package server

import (
	"github.com/jeffscottbrown/applemusic/controllers"
	"net/http"
)

func RunServer() {
	router := http.NewServeMux()
	router.HandleFunc("/search/{term}", controllers.Search)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}

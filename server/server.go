package server

import (
	"github.com/jeffscottbrown/applemusic/auth"
	"github.com/jeffscottbrown/applemusic/search"
	"github.com/jeffscottbrown/applemusic/web"
	"net/http"
)

func Run() {
	router := createAndConfigureRouter()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}

func createAndConfigureRouter() *http.ServeMux {
	router := http.NewServeMux()

	configureApplicationHandlers(router)
	configureStaticResourceHandler(router)

	auth.ConfigureAuthorizationHandlers(router)

	return router
}

func configureApplicationHandlers(router *http.ServeMux) {
	router.HandleFunc("/search/{term}", search.Search)
	router.HandleFunc("/", web.Search)
}

func configureStaticResourceHandler(router *http.ServeMux) {
	dir := http.Dir("./assets/")
	fileServer := http.FileServer(dir)
	foo := http.StripPrefix("/static/", fileServer)
	router.Handle("/static/", foo)
}

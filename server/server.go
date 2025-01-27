package server

import (
	"github.com/gorilla/pat"
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

func createAndConfigureRouter() *pat.Router {
	router := pat.New()

	configureApplicationHandlers(router)
	configureStaticResourceHandler(router)

	auth.ConfigureAuthorizationHandlers(router)

	return router
}

func configureApplicationHandlers(router *pat.Router) {
	router.Get("/search/{term}", search.Search)
	router.HandleFunc("/", web.Search)
}

func configureStaticResourceHandler(router *pat.Router) {
	dir := http.Dir("./assets/")
	fileServer := http.FileServer(dir)
	foo := http.StripPrefix("/static/", fileServer)
	router.Handle("/static/{filepath:.*}", foo)
}

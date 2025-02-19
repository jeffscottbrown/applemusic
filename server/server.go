package server

import (
	"net/http"

	"github.com/jeffscottbrown/applemusic/auth"
	"github.com/jeffscottbrown/applemusic/search"
	"github.com/jeffscottbrown/applemusic/templates"
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
	router.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		bandName := r.FormValue("band_name")
		limit := r.FormValue("limit")
		searchResult, _ := search.SearchApple(bandName, limit)
		templates.Results(searchResult).Render(r.Context(), w)

	})
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		templates.Home(r).Render(r.Context(), w)
	})
}

func configureStaticResourceHandler(router *http.ServeMux) {
	dir := http.Dir("./web/assets/")
	fileServer := http.FileServer(dir)
	handler := http.StripPrefix("/static/", fileServer)
	router.Handle("GET /static/", handler)
}

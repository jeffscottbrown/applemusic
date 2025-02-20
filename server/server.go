package server

import (
	"net/http"

	"github.com/go-chi/chi"
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

func createAndConfigureRouter() *chi.Mux {
	router := chi.NewRouter()
	configureApplicationHandlers(router)
	configureStaticResourceHandler(router)

	auth.ConfigureAuthorizationHandlers(router)

	return router
}

func configureApplicationHandlers(router *chi.Mux) {
	router.Post("/search", func(w http.ResponseWriter, r *http.Request) {
		bandName := r.FormValue("band_name")
		limit := r.FormValue("limit")
		searchResult, _ := search.SearchApple(bandName, limit)
		templates.Results(searchResult).Render(r.Context(), w)

	})
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Home(auth.IsAuthenticated(r)).Render(r.Context(), w)
	})
}

func configureStaticResourceHandler(router *chi.Mux) {
	dir := http.Dir("./web/assets/")
	fs := http.FileServer(dir)
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
}

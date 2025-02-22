package server

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeffscottbrown/applemusic/auth"
	"github.com/jeffscottbrown/applemusic/templates"
	"github.com/jeffscottbrown/goapple/music"
)

func Run() {
	router := createAndConfigureRouter()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}

func createAndConfigureRouter() *gin.Engine {
	router := gin.Default()
	configureApplicationHandlers(router)
	configureStaticResourceHandler(router)

	auth.ConfigureAuthorizationHandlers(router)

	return router
}

func configureApplicationHandlers(router *gin.Engine) {
	router.POST("/search", auth.AuthRequired(), func(c *gin.Context) {
		req := c.Request
		res := c.Writer
		bandName := c.PostForm("band_name")
		limit := c.PostForm("limit")
		searchResult, _ := music.SearchApple(bandName, limit)
		templates.Results(searchResult).Render(req.Context(), res)

	})
	router.GET("/", func(c *gin.Context) {
		req := c.Request
		res := c.Writer
		isAuthenticated := auth.IsAuthenticated(req)

		templates.Home(isAuthenticated).Render(req.Context(), res)
	})
}

//go:embed css
var staticFiles embed.FS

func configureStaticResourceHandler(router *gin.Engine) {
	fileServer := http.FileServer(http.FS(staticFiles))
	x := http.StripPrefix("/static", fileServer)

	router.GET("/static/*filepath", func(c *gin.Context) {
		x.ServeHTTP(c.Writer, c.Request)
	})
}

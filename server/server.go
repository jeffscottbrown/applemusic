package server

import (
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
		r := c.Request
		w := c.Writer
		bandName := c.PostForm("band_name")
		limit := c.PostForm("limit")
		searchResult, _ := music.SearchApple(bandName, limit)
		templates.Results(searchResult).Render(r.Context(), w)

	})
	router.GET("/", func(c *gin.Context) {
		r := c.Request
		w := c.Writer
		templates.Home(auth.IsAuthenticated(r)).Render(r.Context(), w)
	})
}

func configureStaticResourceHandler(router *gin.Engine) {
	router.Static("/static", "./web/assets")
}

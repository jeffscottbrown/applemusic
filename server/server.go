package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffscottbrown/applemusic/controllers"
)

func RunServer() {
	r := gin.Default()

	r.GET("/search/:term", controllers.Search)

	r.Run()
}

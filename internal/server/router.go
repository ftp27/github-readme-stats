package server

import (
	"github.com/ftp27/github-readme-stats/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	api := router.Group("/api")
	{
		api.GET("/", handlers.Stats)
		api.GET("/pin", handlers.Pin)
		api.GET("/top-langs", handlers.TopLangs)
		api.GET("/wakatime", handlers.Wakatime)
		api.GET("/gist", handlers.Gist)

		status := api.Group("/status")
		status.GET("/up", handlers.StatusUp)
		status.GET("/pat-info", handlers.PATInfo)
	}

	return router
}

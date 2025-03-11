package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (h *MusicHandler) InitRoutes() *gin.Engine {
	r := gin.New()

	api := r.Group("/api/v1")

	api.POST("/add", h.AddSong)
	api.GET("/songs", h.GetSongsList)
	api.GET("/:songID/lyrics", h.GetSongLyrics)
	api.PUT("/:songID", h.UpdateSong)
	api.DELETE("/:songID", h.DeleteSong)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

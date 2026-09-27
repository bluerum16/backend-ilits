package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend-ilits/handlers"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	aktaHandler := handlers.NewAktaHandler(db)

	akta := r.Group("/api/akta")
	{
		akta.GET("", aktaHandler.GetAllAkta)
		akta.GET("/:id", aktaHandler.GetAkta)
		akta.POST("", aktaHandler.CreateAkta)
		akta.PUT("/:id", aktaHandler.UpdateAkta)
		akta.DELETE("/:id", aktaHandler.DeleteAkta)
	}

	return r
}

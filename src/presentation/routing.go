package presentation

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(transferHandlers *TransferHandlers) *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: false,
		AllowMethods:     []string{"GET", "POST", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api/v1")
	{
		transfers := api.Group("/transfer")
		{
			transfers.POST("/", transferHandlers.CreateTransfer)
			transfers.POST("/cancel", transferHandlers.CancelTransfer)
			transfers.GET("/status/:id", transferHandlers.GetStatus)
		}
	}

	return router
}

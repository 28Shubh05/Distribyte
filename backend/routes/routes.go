package routes

import (
	"Distribyte/backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Distribyte Backend Running",
		})
	})

	router.POST("/upload", handlers.UploadFile)

	router.GET("/files", handlers.GetFiles)

	router.GET("/deleted-files", handlers.GetDeletedFiles)

	router.GET("/download/:id", handlers.DownloadFile)

	router.DELETE("/files/:id", handlers.DeleteFile)

	router.POST("/restore/:id", handlers.RestoreFile)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
}

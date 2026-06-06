package routes

import (
	"Distribyte/backend/handlers"
	"Distribyte/backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "Distribyte Backend Running",
		})
	})

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	router.GET("/health", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	auth := router.Group("/")

	auth.Use(
		middleware.AuthMiddleware(),
	)

	{
		auth.POST(
			"/upload",
			handlers.UploadFile,
		)

		auth.GET(
			"/files",
			handlers.GetFiles,
		)

		auth.GET(
			"/deleted-files",
			handlers.GetDeletedFiles,
		)

		auth.GET(
			"/download/:id",
			handlers.DownloadFile,
		)

		auth.DELETE(
			"/files/:id",
			handlers.DeleteFile,
		)

		auth.POST(
			"/restore/:id",
			handlers.RestoreFile,
		)
	}
}

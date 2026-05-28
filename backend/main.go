package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Distribyte Backend Running",
		})
	})

	router.POST("/upload", func(c *gin.Context) {

		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "File upload failed",
			})
			return
		}

		filename := filepath.Base(file.Filename)

		savePath := "../storage/" + filename

		err = c.SaveUploadedFile(file, savePath)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not save file",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
			"file":    filename,
		})
	})

	router.Run(":8080")
}

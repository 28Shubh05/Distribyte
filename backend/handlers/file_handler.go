package handlers

import (
	"net/http"
	"path/filepath"

	"log"

	"Distribyte/backend/database"
	"Distribyte/backend/models"
	"Distribyte/backend/services"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "File upload failed",
		})
		return
	}

	filename := filepath.Base(file.Filename)

	savePath := "../storage/" + filename

	err = c.SaveUploadedFile(file, savePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Could not save file",
		})
		return
	}

	fileData, err := services.SaveFileMetadata(
		filename,
		savePath,
		file.Size,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database insert failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"data":    fileData,
	})

	log.Println("File uploaded:", filename)
}

func GetFiles(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT id, filename, filepath, size, uploaded_at
		FROM files
		ORDER BY uploaded_at DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch files",
		})
		return
	}

	defer rows.Close()

	var files []models.File

	for rows.Next() {

		var file models.File

		err := rows.Scan(
			&file.ID,
			&file.Filename,
			&file.Filepath,
			&file.Size,
			&file.UploadedAt,
		)

		if err != nil {
			continue
		}

		files = append(files, file)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   files,
	})
}

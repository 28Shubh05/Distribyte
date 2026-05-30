package handlers

import (
	"net/http"
	"path/filepath"

	"log"

	"Distribyte/backend/database"
	"Distribyte/backend/models"
	"Distribyte/backend/services"
	"Distribyte/backend/utils"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {

		log.Println("DB ERROR:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	filename := filepath.Base(file.Filename)
	storedName := utils.GenerateStoredName(filename)

	savePath := "../storage/" + storedName

	err = c.SaveUploadedFile(file, savePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Could not save file",
		})
		return
	}

	hash, err := utils.GenerateSHA256(savePath)

	if err != nil {
		log.Println(err)
	}

	log.Println("SHA256:", hash)

	fileData, err := services.SaveFileMetadata(
		filename,
		storedName,
		savePath,
		file.Size,
		hash,
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
    	SELECT
    	    id,
    	    original_name,
    	    stored_name,
    	    filepath,
    	    size,
    	    file_hash,
    	    uploaded_at
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
			&file.OriginalName,
			&file.StoredName,
			&file.Filepath,
			&file.Size,
			&file.FileHash,
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

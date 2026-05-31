package handlers

import (
	"net/http"
	"path/filepath"

	"log"
	"os"

	"Distribyte/backend/config"
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

	if file.Size > config.MaxFileSize {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "File exceeds maximum size limit (10 MB)",
		})
		return
	}

	if !utils.IsAllowedFileType(file.Filename) {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Unsupported file type",
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

	exists, err := services.HashExists(hash)

	if exists {
		os.Remove(savePath)
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Duplicate file detected",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Duplicate check failed",
		})
		return
	}

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
		WHERE is_deleted = FALSE
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

func DownloadFile(c *gin.Context) {

	id := c.Param("id")

	file, err := services.GetFileByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "File not found",
		})
		return
	}

	c.FileAttachment(
		file.Filepath,
		file.OriginalName,
	)
}

func DeleteFile(c *gin.Context) {

	id := c.Param("id")

	// Verify file exists
	_, err := services.GetFileByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "File not found",
		})
		return
	}

	// Soft delete metadata
	err = services.SoftDeleteFile(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

func RestoreFile(c *gin.Context) {

	id := c.Param("id")

	err := services.RestoreFile(id)

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File restored successfully",
	})
}

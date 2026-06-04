package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

	existingFile, found, err := services.GetFileByHash(hash)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Duplicate check failed",
		})

		return
	}

	if found {

		os.Remove(savePath)

		if existingFile.IsDeleted {

			err = services.RestoreFile(
				strconv.Itoa(existingFile.ID),
			)

			if err != nil {

				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "Failed to restore existing file",
				})

				return
			}

			services.ClearFileCaches()

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Existing deleted file restored",
			})

			return
		}

		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Duplicate file detected",
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

	services.ClearFileCaches()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"data":    fileData,
	})

	log.Println("File uploaded:", filename)
}

func GetFiles(c *gin.Context) {

	cachedFiles, err := database.RedisClient.Get(
		database.Ctx,
		"files:list",
	).Result()

	if err == nil {

		log.Println("CACHE HIT: files:list")

		var files []models.File

		json.Unmarshal(
			[]byte(cachedFiles),
			&files,
		)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"files":   files,
		})

		return
	}

	log.Println("CACHE MISS: files:list")

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

	jsonData, _ := json.Marshal(files)

	database.RedisClient.Set(
		database.Ctx,
		"files:list",
		jsonData,
		5*time.Minute,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   files,
	})
}

func GetDeletedFiles(c *gin.Context) {

	cachedFiles, err := database.RedisClient.Get(
		database.Ctx,
		"deleted_files:list",
	).Result()

	if err == nil {

		log.Println("CACHE HIT: deleted_files:list")

		var files []models.File

		json.Unmarshal(
			[]byte(cachedFiles),
			&files,
		)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"files":   files,
		})

		return
	}

	log.Println("CACHE MISS: deleted_files:list")

	files, err := services.GetDeletedFiles()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch deleted files",
		})

		return
	}

	jsonData, _ := json.Marshal(files)

	database.RedisClient.Set(
		database.Ctx,
		"deleted_files:list",
		jsonData,
		5*time.Minute,
	)

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

	services.ClearFileCaches()

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

	services.ClearFileCaches()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File restored successfully",
	})
}

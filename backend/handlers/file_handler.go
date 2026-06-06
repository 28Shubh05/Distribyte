package handlers

import (
	"database/sql"
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
	"Distribyte/backend/utils"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file uploaded",
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

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Could not save file",
		})
		return
	}

	hash, err := utils.GenerateSHA256(savePath)
	if err != nil {
		_ = os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate hash",
		})
		return
	}

	existingFile, found, err := getUserFileByHash(userID, hash)
	if err != nil {
		_ = os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Duplicate check failed",
		})
		return
	}

	if found {
		_ = os.Remove(savePath)

		if existingFile.IsDeleted {
			restored, err := restoreUserFile(userID, itoa(existingFile.ID))
			if err != nil || !restored {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "Failed to restore existing file",
				})
				return
			}

			clearUserCaches(userID)

			existingFile.IsDeleted = false
			existingFile.DeletedAt = nil

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Existing deleted file restored",
				"data":    existingFile,
			})

			userIDInterface, _ := c.Get("user_id")

			log.Println("USER ID:", userIDInterface)
			return
		}

		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Duplicate file detected",
		})
		return
	}

	log.Println("SHA256:", hash)

	fileData, err := insertFileMetadata(
		userID,
		filename,
		storedName,
		savePath,
		file.Size,
		hash,
	)
	if err != nil {
		_ = os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database insert failed",
		})
		return
	}

	clearUserCaches(userID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"data":    fileData,
	})

	log.Println("File uploaded:", filename)

	userIDInterface, _ := c.Get("user_id")

	log.Println("USER ID:", userIDInterface)
}

func GetFiles(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	cacheKey := activeFilesCacheKey(userID)

	cachedFiles, err := database.RedisClient.Get(
		database.Ctx,
		cacheKey,
	).Result()

	if err == nil {
		log.Println("CACHE HIT:", cacheKey)

		var files []models.File
		_ = json.Unmarshal([]byte(cachedFiles), &files)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"files":   files,
		})
		return
	}

	log.Println("CACHE MISS:", cacheKey)

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
		WHERE
			is_deleted = FALSE
			AND user_id = $1
		ORDER BY uploaded_at DESC
	`, userID)

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

	_ = database.RedisClient.Set(
		database.Ctx,
		cacheKey,
		jsonData,
		5*time.Minute,
	).Err()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   files,
	})
}

func GetDeletedFiles(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	cacheKey := deletedFilesCacheKey(userID)

	cachedFiles, err := database.RedisClient.Get(
		database.Ctx,
		cacheKey,
	).Result()

	if err == nil {
		log.Println("CACHE HIT:", cacheKey)

		var files []models.File
		_ = json.Unmarshal([]byte(cachedFiles), &files)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"files":   files,
		})
		return
	}

	log.Println("CACHE MISS:", cacheKey)

	rows, err := database.DB.Query(`
		SELECT
			id,
			original_name,
			stored_name,
			filepath,
			size,
			file_hash,
			uploaded_at,
			is_deleted,
			deleted_at
		FROM files
		WHERE
			is_deleted = TRUE
			AND user_id = $1
		ORDER BY deleted_at DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch deleted files",
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
			&file.IsDeleted,
			&file.DeletedAt,
		)
		if err != nil {
			continue
		}

		files = append(files, file)
	}

	jsonData, _ := json.Marshal(files)

	_ = database.RedisClient.Set(
		database.Ctx,
		cacheKey,
		jsonData,
		5*time.Minute,
	).Err()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"files":   files,
	})
}

func DownloadFile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	id := c.Param("id")

	file, found, err := getUserFileByID(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch file",
		})
		return
	}

	if !found {
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
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	id := c.Param("id")

	found, err := softDeleteUserFile(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete file",
		})
		return
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "File not found",
		})
		return
	}

	clearUserCaches(userID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

func RestoreFile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	id := c.Param("id")

	found, err := restoreUserFile(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to restore file",
		})
		return
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "File not found",
		})
		return
	}

	clearUserCaches(userID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File restored successfully",
	})
}

func getCurrentUserID(c *gin.Context) (int, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}

	switch t := v.(type) {
	case int:
		return t, true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}

func activeFilesCacheKey(userID int) string {
	return "files:list:user:" + itoa(userID)
}

func deletedFilesCacheKey(userID int) string {
	return "deleted_files:list:user:" + itoa(userID)
}

func clearUserCaches(userID int) {
	_ = database.RedisClient.Del(
		database.Ctx,
		activeFilesCacheKey(userID),
		deletedFilesCacheKey(userID),
	).Err()

	log.Println("CACHE CLEARED for user", userID)
}

func insertFileMetadata(
	userID int,
	originalName string,
	storedName string,
	savePath string,
	size int64,
	fileHash string,
) (models.File, error) {

	query := `
		INSERT INTO files (
			original_name,
			stored_name,
			filepath,
			size,
			file_hash,
			user_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uploaded_at
	`

	var file models.File

	err := database.DB.QueryRow(
		query,
		originalName,
		storedName,
		savePath,
		size,
		fileHash,
		userID,
	).Scan(
		&file.ID,
		&file.UploadedAt,
	)

	if err != nil {
		return file, err
	}

	file.OriginalName = originalName
	file.StoredName = storedName
	file.Filepath = savePath
	file.Size = size
	file.FileHash = fileHash

	return file, nil
}

func getUserFileByID(userID int, id string) (models.File, bool, error) {
	var file models.File

	query := `
		SELECT
			id,
			original_name,
			stored_name,
			filepath,
			size,
			file_hash,
			uploaded_at,
			is_deleted,
			deleted_at
		FROM files
		WHERE id = $1
		AND user_id = $2
		AND is_deleted = FALSE
		LIMIT 1
	`

	err := database.DB.QueryRow(
		query,
		id,
		userID,
	).Scan(
		&file.ID,
		&file.OriginalName,
		&file.StoredName,
		&file.Filepath,
		&file.Size,
		&file.FileHash,
		&file.UploadedAt,
		&file.IsDeleted,
		&file.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return file, false, nil
	}

	if err != nil {
		return file, false, err
	}

	return file, true, nil
}

func getUserFileByHash(userID int, hash string) (models.File, bool, error) {
	var file models.File

	query := `
		SELECT
			id,
			original_name,
			stored_name,
			filepath,
			size,
			file_hash,
			uploaded_at,
			is_deleted,
			deleted_at
		FROM files
		WHERE file_hash = $1
		AND user_id = $2
		LIMIT 1
	`

	err := database.DB.QueryRow(
		query,
		hash,
		userID,
	).Scan(
		&file.ID,
		&file.OriginalName,
		&file.StoredName,
		&file.Filepath,
		&file.Size,
		&file.FileHash,
		&file.UploadedAt,
		&file.IsDeleted,
		&file.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return file, false, nil
	}

	if err != nil {
		return file, false, err
	}

	return file, true, nil
}

func softDeleteUserFile(userID int, id string) (bool, error) {
	query := `
		UPDATE files
		SET
			is_deleted = TRUE,
			deleted_at = NOW()
		WHERE id = $1
		AND user_id = $2
		AND is_deleted = FALSE
	`

	result, err := database.DB.Exec(
		query,
		id,
		userID,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func restoreUserFile(userID int, id string) (bool, error) {
	query := `
		UPDATE files
		SET
			is_deleted = FALSE,
			deleted_at = NULL
		WHERE id = $1
		AND user_id = $2
		AND is_deleted = TRUE
	`

	result, err := database.DB.Exec(
		query,
		id,
		userID,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func itoa(v int) string {
	return strconv.FormatInt(int64(v), 10)
}

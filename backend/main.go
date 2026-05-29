package main

import (
	"net/http"
	"path/filepath"

	"Distribyte/backend/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type FileMetadata struct {
	ID         int
	Filename   string
	Filepath   string
	Size       int64
	UploadedAt string
}

func main() {

	err := godotenv.Load("../.env")

	if err != nil {
		panic("Error loading .env file")
	}

	database.ConnectDatabase()

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

		query := `
		INSERT INTO files (filename, filepath, size)
		VALUES ($1, $2, $3)
		RETURNING id, uploaded_at
		`

		var id int
		var uploadedAt string

		err = database.DB.QueryRow(
			query,
			filename,
			savePath,
			file.Size,
		).Scan(&id, &uploadedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database insert failed",
			})
			return
		}

		fileData := FileMetadata{
			ID:         id,
			Filename:   filename,
			Filepath:   savePath,
			Size:       file.Size,
			UploadedAt: uploadedAt,
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
			"data":    fileData,
		})
	})

	router.GET("/files", func(c *gin.Context) {

		rows, err := database.DB.Query(`
			SELECT id, filename, filepath, size, uploaded_at
			FROM files
			ORDER BY uploaded_at DESC
		`)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch files",
			})
			return
		}

		defer rows.Close()

		var files []FileMetadata

		for rows.Next() {

			var file FileMetadata

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
			"files": files,
		})
	})

	router.Run(":8080")
}

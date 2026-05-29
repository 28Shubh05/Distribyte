package services

import (
	"path/filepath"

	"Distribyte/backend/database"
	"Distribyte/backend/models"
)

func SaveFileMetadata(
	filename string,
	savePath string,
	size int64,
) (models.File, error) {

	query := `
	INSERT INTO files (filename, filepath, size)
	VALUES ($1, $2, $3)
	RETURNING id, uploaded_at
	`

	var file models.File

	err := database.DB.QueryRow(
		query,
		filename,
		savePath,
		size,
	).Scan(
		&file.ID,
		&file.UploadedAt,
	)

	if err != nil {
		return file, err
	}

	file.Filename = filepath.Base(filename)
	file.Filepath = savePath
	file.Size = size

	return file, nil
}

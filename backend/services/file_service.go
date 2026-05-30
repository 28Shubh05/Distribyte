package services

import (
	"Distribyte/backend/database"
	"Distribyte/backend/models"
)

func SaveFileMetadata(
	originalName string,
	storedName string,
	savePath string,
	size int64,
) (models.File, error) {

	query := `
	INSERT INTO files (
    original_name,
    stored_name,
    filepath,
    size
)
VALUES ($1,$2,$3,$4)
RETURNING id, uploaded_at
	`

	var file models.File

	err := database.DB.QueryRow(
		query,
		originalName,
		storedName,
		savePath,
		size,
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

	return file, nil
}

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
	fileHash string,
) (models.File, error) {

	query := `
		INSERT INTO files (
    	original_name,
    	stored_name,
    	filepath,
    	size,
    	file_hash
	)
	VALUES ($1,$2,$3,$4,$5)
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

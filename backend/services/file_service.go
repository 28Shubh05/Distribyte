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

func HashExists(hash string) (bool, error) {

	var count int

	err := database.DB.QueryRow(
		"SELECT COUNT(*) FROM files WHERE file_hash=$1",
		hash,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetFileByID(id string) (models.File, error) {

	var file models.File

	query := `
	SELECT
		id,
		original_name,
		stored_name,
		filepath,
		size,
		file_hash,
		uploaded_at
	FROM files
	WHERE id = $1
	`

	err := database.DB.QueryRow(
		query,
		id,
	).Scan(
		&file.ID,
		&file.OriginalName,
		&file.StoredName,
		&file.Filepath,
		&file.Size,
		&file.FileHash,
		&file.UploadedAt,
	)

	return file, err
}

func DeleteFileMetadata(id string) error {

	query := `
	DELETE FROM files
	WHERE id = $1
	`

	_, err := database.DB.Exec(
		query,
		id,
	)

	return err
}

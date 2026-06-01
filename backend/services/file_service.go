package services

import (
	"Distribyte/backend/database"
	"Distribyte/backend/models"
	"database/sql"
	"errors"
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
	AND is_deleted = FALSE
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

func GetFileByHash(hash string) (models.File, bool, error) {

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
	LIMIT 1
	`

	err := database.DB.QueryRow(
		query,
		hash,
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

func SoftDeleteFile(id string) error {

	query := `
	UPDATE files
	SET
		is_deleted = TRUE,
		deleted_at = NOW()
	WHERE id = $1
	`

	_, err := database.DB.Exec(
		query,
		id,
	)

	return err
}

func RestoreFile(id string) error {

	query := `
	UPDATE files
	SET
		is_deleted = FALSE,
		deleted_at = NULL
	WHERE id = $1
	`

	result, err := database.DB.Exec(
		query,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("file not found")
	}

	return nil
}

func GetDeletedFiles() ([]models.File, error) {

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
		WHERE is_deleted = TRUE
		ORDER BY deleted_at DESC
	`)

	if err != nil {
		return nil, err
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

	return files, nil
}

package services

import (
	"Distribyte/backend/database"
	"io"
	"os"
	"path/filepath"
)

func SaveReplica(
	fileID int,
	nodeID int,
	replicaPath string,
) error {

	query := `
		INSERT INTO file_replicas (
			file_id,
			node_id,
			replica_path
		)
		VALUES ($1,$2,$3)
	`

	_, err := database.DB.Exec(
		query,
		fileID,
		nodeID,
		replicaPath,
	)

	return err
}

func CopyFile(
	source string,
	destination string,
) error {

	srcFile, err := os.Open(source)

	if err != nil {
		return err
	}

	defer srcFile.Close()

	err = os.MkdirAll(
		filepath.Dir(destination),
		os.ModePerm,
	)

	if err != nil {
		return err
	}

	dstFile, err := os.Create(destination)

	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(
		dstFile,
		srcFile,
	)

	return err
}

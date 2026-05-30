package utils

import (
	"path/filepath"

	"github.com/google/uuid"
)

func GenerateStoredName(filename string) string {

	ext := filepath.Ext(filename)

	return uuid.New().String() + ext
}

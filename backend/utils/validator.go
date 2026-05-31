package utils

import (
	"path/filepath"
	"strings"
)

var AllowedExtensions = map[string]bool{
	".pdf":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".txt":  true,
}

func IsAllowedFileType(filename string) bool {

	ext := strings.ToLower(
		filepath.Ext(filename),
	)

	return AllowedExtensions[ext]
}

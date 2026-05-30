package models

type File struct {
	ID           int    `json:"id"`
	OriginalName string `json:"original_name"`
	StoredName   string `json:"stored_name"`
	Filepath     string `json:"filepath"`
	Size         int64  `json:"size"`
	UploadedAt   string `json:"uploaded_at"`
}

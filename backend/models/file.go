package models

type File struct {
	ID           int     `json:"id"`
	OriginalName string  `json:"original_name"`
	StoredName   string  `json:"stored_name"`
	Filepath     string  `json:"filepath"`
	Size         int64   `json:"size"`
	FileHash     string  `json:"file_hash"`
	UploadedAt   string  `json:"uploaded_at"`
	IsDeleted    bool    `json:"is_deleted"`
	DeletedAt    *string `json:"deleted_at"`
}

package models

type File struct {
	ID         int    `json:"id"`
	Filename   string `json:"filename"`
	Filepath   string `json:"filepath"`
	Size       int64  `json:"size"`
	UploadedAt string `json:"uploaded_at"`
}

package models

type StorageNode struct {
	ID          int    `json:"id"`
	NodeName    string `json:"node_name"`
	StoragePath string `json:"storage_path"`

	TotalSpace int64 `json:"total_space"`
	UsedSpace  int64 `json:"used_space"`

	IsActive bool `json:"is_active"`
}

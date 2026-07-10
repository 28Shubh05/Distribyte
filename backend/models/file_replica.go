package models

type FileReplica struct {
	ID          int    `json:"id"`
	FileID      int    `json:"file_id"`
	NodeID      int    `json:"node_id"`
	ReplicaPath string `json:"replica_path"`
}

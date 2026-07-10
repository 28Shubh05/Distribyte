package services

import (
	"Distribyte/backend/database"
	"Distribyte/backend/models"
)

func GetAvailableNode() (models.StorageNode, error) {

	var node models.StorageNode

	query := `
		SELECT
			id,
			node_name,
			storage_path,
			total_space,
			used_space,
			is_active
		FROM storage_nodes
		WHERE
			is_active = TRUE
			AND node_type = 'PRIMARY'
		ORDER BY used_space ASC
		LIMIT 1
	`

	err := database.DB.QueryRow(query).Scan(
		&node.ID,
		&node.NodeName,
		&node.StoragePath,
		&node.TotalSpace,
		&node.UsedSpace,
		&node.IsActive,
	)

	if err != nil {
		return node, err
	}

	return node, nil
}

func GetNodeByID(nodeID int) (models.StorageNode, error) {

	var node models.StorageNode

	query := `
		SELECT
			id,
			node_name,
			storage_path,
			total_space,
			used_space,
			is_active
		FROM storage_nodes
		WHERE id = $1
	`

	err := database.DB.QueryRow(
		query,
		nodeID,
	).Scan(
		&node.ID,
		&node.NodeName,
		&node.StoragePath,
		&node.TotalSpace,
		&node.UsedSpace,
		&node.IsActive,
	)

	if err != nil {
		return node, err
	}

	return node, nil
}

func UpdateNodeUsage(
	nodeID int,
	fileSize int64,
) error {

	query := `
		UPDATE storage_nodes
		SET used_space = used_space + $1
		WHERE id = $2
	`

	_, err := database.DB.Exec(
		query,
		fileSize,
		nodeID,
	)

	return err
}

func ReduceNodeUsage(
	nodeID int,
	fileSize int64,
) error {

	query := `
		UPDATE storage_nodes
		SET used_space = GREATEST(
			used_space - $1,
			0
		)
		WHERE id = $2
	`

	_, err := database.DB.Exec(
		query,
		fileSize,
		nodeID,
	)

	return err
}

func GetAllNodes() ([]models.StorageNode, error) {

	rows, err := database.DB.Query(`
		SELECT
			id,
			node_name,
			storage_path,
			total_space,
			used_space,
			is_active
		FROM storage_nodes
		ORDER BY id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var nodes []models.StorageNode

	for rows.Next() {

		var node models.StorageNode

		err := rows.Scan(
			&node.ID,
			&node.NodeName,
			&node.StoragePath,
			&node.TotalSpace,
			&node.UsedSpace,
			&node.IsActive,
		)

		if err != nil {
			continue
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func MarkNodeOffline(
	nodeID int,
) error {

	query := `
		UPDATE storage_nodes
		SET is_active = FALSE
		WHERE id = $1
	`

	_, err := database.DB.Exec(
		query,
		nodeID,
	)

	return err
}

func MarkNodeOnline(
	nodeID int,
) error {

	query := `
		UPDATE storage_nodes
		SET is_active = TRUE
		WHERE id = $1
	`

	_, err := database.DB.Exec(
		query,
		nodeID,
	)

	return err
}

func GetReplicaForNode(
	primaryNodeID int,
) (models.StorageNode, error) {

	var node models.StorageNode

	query := `
		SELECT
			sn.id,
			sn.node_name,
			sn.storage_path,
			sn.total_space,
			sn.used_space,
			sn.is_active
		FROM storage_nodes sn
		INNER JOIN node_replica_mapping m
			ON sn.id = m.replica_node_id
		WHERE
			m.primary_node_id = $1
			AND sn.is_active = TRUE
		LIMIT 1
	`

	err := database.DB.QueryRow(
		query,
		primaryNodeID,
	).Scan(
		&node.ID,
		&node.NodeName,
		&node.StoragePath,
		&node.TotalSpace,
		&node.UsedSpace,
		&node.IsActive,
	)

	if err != nil {
		return node, err
	}

	return node, nil
}

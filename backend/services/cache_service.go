package services

import (
	"log"

	"Distribyte/backend/database"
)

func ClearFileCaches() {

	database.RedisClient.Del(
		database.Ctx,
		"files:list",
	)

	database.RedisClient.Del(
		database.Ctx,
		"deleted_files:list",
	)

	log.Println("CACHE CLEARED")
}

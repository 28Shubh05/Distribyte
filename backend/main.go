package main

import (
	"Distribyte/backend/config"
	"Distribyte/backend/database"
	"Distribyte/backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	database.ConnectDatabase()

	router := gin.Default()

	router.Use(cors.Default())

	routes.SetupRoutes(router)

	router.Run(":8080")
}

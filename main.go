package main

import (
	"test-golang/config"
	"test-golang/models"
	"test-golang/routes" // Import routes package

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()
	r.Use(cors.Default()) // Enable CORS
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Company{})
	routes.SetupRoutes(r)
	r.Run(":8080")
}

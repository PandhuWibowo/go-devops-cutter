package main

import (
	"log"
	"os"

	"github.com/PandhuWibowo/go-devops-cutter/internal/api"
	"github.com/PandhuWibowo/go-devops-cutter/internal/database"
	"github.com/gin-gonic/gin"
)

var version = "0.1.0"

func main() {
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://devops:devops123@localhost:5432/devops_cutter?sslmode=disable")

	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": version})
	})

	api.SetupRoutes(router, db)

	log.Printf("Starting API server on port %s (version: %s)", port, version)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

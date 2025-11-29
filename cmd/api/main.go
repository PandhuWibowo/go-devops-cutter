package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var version = "0.1.0"

func main() {
	port := getEnv("PORT", "8080")

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": version})
	})

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

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PandhuWibowo/go-devops-cutter/internal/database"
	"github.com/PandhuWibowo/go-devops-cutter/internal/queue"
	"github.com/hibiken/asynq"
)

var version = "0.1.0"

func main() {
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	dbURL := getEnv("DATABASE_URL", "postgres://devops:devops123@localhost:5432/devops_cutter?sslmode=disable")

	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 5},
	)

	mux := queue.NewHandler(db)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("Could not run worker: %v", err)
		}
	}()

	log.Printf("Worker started (version: %s)", version)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	srv.Shutdown()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

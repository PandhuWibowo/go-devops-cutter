package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PandhuWibowo/go-devops-cutter/internal/database"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const TypeBackupDatabase = "backup:database"

type BackupPayload struct {
	RequestID    uint   `json:"request_id"`
	DatabaseName string `json:"database_name"`
	Format       string `json:"format"`
	Compression  bool   `json:"compression"`
	BackupType   string `json:"backup_type"`
}

func NewHandler(db *gorm.DB) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	handler := &Handler{db: db}
	mux.HandleFunc(TypeBackupDatabase, handler.HandleBackupDatabase)
	return mux
}

type Handler struct {
	db *gorm.DB
}

func (h *Handler) HandleBackupDatabase(ctx context.Context, t *asynq.Task) error {
	var payload BackupPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	log.Printf("Processing backup request ID: %d", payload.RequestID)

	var req database.BackupRequest
	h.db.First(&req, payload.RequestID)
	req.Status = "processing"
	h.db.Save(&req)

	// TODO: Implement actual backup logic
	time.Sleep(5 * time.Second)

	expiresAt := time.Now().Add(24 * time.Hour)
	req.Status = "completed"
	req.DownloadURL = fmt.Sprintf("https://storage.example.com/backups/%d.sql.gz", payload.RequestID)
	req.ExpiresAt = &expiresAt
	req.FileSize = 1024 * 1024 * 100
	h.db.Save(&req)

	log.Printf("Backup request ID %d completed", payload.RequestID)
	return nil
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/PandhuWibowo/go-devops-cutter/internal/database"
	"github.com/PandhuWibowo/go-devops-cutter/internal/queue"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type BackupHandler struct {
	db     *gorm.DB
	client *asynq.Client
}

func NewBackupHandler(db *gorm.DB) *BackupHandler {
	redisAddr := "localhost:6379"
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &BackupHandler{db: db, client: client}
}

type CreateBackupRequest struct {
	DatabaseName string `json:"database_name" binding:"required"`
	Format       string `json:"format" binding:"required,oneof=sql custom directory"`
	Compression  bool   `json:"compression"`
	BackupType   string `json:"backup_type" binding:"required,oneof=full schema data"`
}

func (h *BackupHandler) Create(c *gin.Context) {
	var req CreateBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	var count int64
	h.db.Model(&database.BackupRequest{}).
		Where("user_id = ? AND created_at >= CURRENT_DATE", userID).
		Count(&count)

	if count >= 5 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "daily quota exceeded (max 5 backups per day)"})
		return
	}

	compression := "none"
	if req.Compression {
		compression = "gzip"
	}

	backupReq := database.BackupRequest{
		UserID:       userID,
		DatabaseName: req.DatabaseName,
		Format:       req.Format,
		Compression:  compression,
		BackupType:   req.BackupType,
		Status:       "pending",
	}

	if err := h.db.Create(&backupReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create backup request"})
		return
	}

	payload, _ := json.Marshal(queue.BackupPayload{
		RequestID:    backupReq.ID,
		DatabaseName: req.DatabaseName,
		Format:       req.Format,
		Compression:  req.Compression,
		BackupType:   req.BackupType,
	})

	task := asynq.NewTask(queue.TypeBackupDatabase, payload)
	info, err := h.client.Enqueue(task, asynq.Queue("default"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue backup job"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "backup request queued",
		"request_id": backupReq.ID,
		"job_id":     info.ID,
	})
}

func (h *BackupHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var backups []database.BackupRequest
	h.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&backups)
	c.JSON(http.StatusOK, backups)
}

func (h *BackupHandler) GetStatus(c *gin.Context) {
	requestID := c.Param("id")
	userID := c.MustGet("user_id").(uint)
	var backup database.BackupRequest
	if err := h.db.Where("id = ? AND user_id = ?", requestID, userID).First(&backup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "backup request not found"})
		return
	}
	c.JSON(http.StatusOK, backup)
}

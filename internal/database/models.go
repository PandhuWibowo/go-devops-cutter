package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint             `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
	Username       string           `gorm:"uniqueIndex;not null" json:"username"`
	Email          string           `gorm:"uniqueIndex;not null" json:"email"`
	Password       string           `gorm:"not null" json:"-"`
	TelegramID     string           `json:"telegram_id,omitempty"`
	Role           string           `gorm:"default:'developer'" json:"role"`
	BackupRequests []BackupRequest  `json:"-"`
	DatabaseAccess []DatabaseAccess `json:"-"`
}

type BackupRequest struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `json:"user_id"`
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DatabaseName string         `gorm:"not null" json:"database_name"`
	Format       string         `gorm:"not null" json:"format"`
	Compression  string         `gorm:"not null" json:"compression"`
	BackupType   string         `gorm:"not null" json:"backup_type"`
	Status       string         `gorm:"default:'pending'" json:"status"`
	DownloadURL  string         `json:"download_url,omitempty"`
	ObjectKey    string         `json:"-"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	FileSize     int64          `json:"file_size"`
}

type DatabaseAccess struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `json:"user_id"`
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DatabaseName string         `gorm:"not null" json:"database_name"`
	AccessLevel  string         `gorm:"not null" json:"access_level"`
	Status       string         `gorm:"default:'pending'" json:"status"`
	DBUsername   string         `json:"db_username,omitempty"`
	DBPassword   string         `json:"-"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty"`
}

type LogExportRequest struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `json:"user_id"`
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ServiceName  string         `gorm:"not null" json:"service_name"`
	FromDate     time.Time      `gorm:"not null" json:"from_date"`
	ToDate       time.Time      `gorm:"not null" json:"to_date"`
	Keyword      string         `json:"keyword,omitempty"`
	Status       string         `gorm:"default:'pending'" json:"status"`
	DownloadURL  string         `json:"download_url,omitempty"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	FileSize     int64          `json:"file_size"`
}

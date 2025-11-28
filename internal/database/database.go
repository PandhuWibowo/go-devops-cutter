package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Connect(dsn string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&BackupRequest{},
		&DatabaseAccess{},
		&LogExportRequest{},
	)
}

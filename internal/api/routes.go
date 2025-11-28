package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	public := r.Group("/api/v1")
	{
		auth := NewAuthHandler(db)
		public.POST("/auth/login", auth.Login)
		public.POST("/auth/register", auth.Register)
	}

	protected := r.Group("/api/v1")
	protected.Use(AuthMiddleware())
	{
		backup := NewBackupHandler(db)
		protected.POST("/backups", backup.Create)
		protected.GET("/backups", backup.List)
		protected.GET("/backups/:id", backup.GetStatus)
	}
}

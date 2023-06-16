package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/rickyromansyah2045/halocat-backend-go/helper"
	"gorm.io/gorm"
)

type Repository interface {
	SaveActivityLog(ActivityLog) (ActivityLog, error)
	DeleteActivityLog(ActivityLog) (bool, error)

	SaveActivityWebhook(ActivityWebhook) (ActivityWebhook, error)

	AdminDataTablesActivityLogs(ctx *gin.Context) (helper.DataTables, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}

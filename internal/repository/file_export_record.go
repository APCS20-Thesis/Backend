package repository

import (
	"github.com/APCS20-Thesis/Backend/internal/model"
	"gorm.io/gorm"
)

type FileExportRecordRepository interface {
}

type fileExportRecordRepo struct {
	*gorm.DB
	TableName string
}

func NewFileExportRecordRepository(db *gorm.DB) FileExportRecordRepository {
	return &fileExportRecordRepo{db, model.FileExportRecord{}.TableName()}
}

package model

import (
	"github.com/google/uuid"
	"time"
)

type FileType string

const (
	FileType_CSV     FileType = "CSV"
	FileType_PARQUET FileType = "PARQUET"
)

type FileExportRecord struct {
	ID              int64 `gorm:"primaryKey"`
	DataTableId     int64
	Type            FileType
	AccountUuid     uuid.UUID
	DataActionId    int64
	DataActionRunId int64
	DownloadUrl     string
	S3Key           string
	ExpireTime      time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (FileExportRecord) TableName() string {
	return "file_export_record"
}

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
	Format          FileType
	AccountUuid     uuid.UUID
	DataActionId    int64
	DataActionRunId int64
	Status          string
	DownloadUrl     string
	S3Key           string `gorm:"column:s3_key"`
	ExpirationTime  time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (FileExportRecord) TableName() string {
	return "file_export_record"
}

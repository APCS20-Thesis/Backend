package model

import (
	"github.com/google/uuid"
	"time"
)

type DataActionRunStatus string

const (
	DataActionRunStatus_Success    = "DONE"
	DataActionRunStatus_Processing = "PROCESSING"
	DataActionRunStatus_Failed     = "FAILED"
)

type DataActionRun struct {
	ID          int64 `gorm:"primaryKey"`
	EventId     int64
	RunId       int64
	Status      string
	Error       string
	AccountUuid uuid.UUID
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataActionRun) TableName() string {
	return "data_action_run"
}

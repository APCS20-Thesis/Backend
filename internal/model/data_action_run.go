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
	DataActionRunStatus_Canceled   = "CANCELED"
)

type DataActionRun struct {
	ID          int64 `gorm:"primaryKey"`
	ActionId    int64
	RunId       int64
	Status      DataActionRunStatus
	Error       string
	AccountUuid uuid.UUID
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataActionRun) TableName() string {
	return "data_action_run"
}

package model

import (
	"github.com/google/uuid"
	"time"
)

type DataActionRunStatus string

const (
	DataActionRunStatus_Success    DataActionRunStatus = "SUCCESS"
	DataActionRunStatus_Processing DataActionRunStatus = "PROCESSING"
	DataActionRunStatus_Failed     DataActionRunStatus = "FAILED"
	DataActionRunStatus_Canceled   DataActionRunStatus = "CANCELED"
	DataActionRunStatus_Creating   DataActionRunStatus = "CREATING"
)

type DataActionRun struct {
	ID          int64 `gorm:"primaryKey"`
	ActionId    int64
	RunId       int64
	DagRunId    string
	Status      DataActionRunStatus
	Error       string
	AccountUuid uuid.UUID
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataActionRun) TableName() string {
	return "data_action_run"
}

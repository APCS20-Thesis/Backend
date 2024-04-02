package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type ActionType string

const (
	ActionType_ImportDataFromFile ActionType = "IMPORT_DATA_FROM_FILE"
	ActionType_ImportDataFromS3   ActionType = "IMPORT_DATA_FROM_S3"
	ActionType_ExportDataToCSV    ActionType = "EXPORT_DATA_TO_CSV"
)

type DataActionStatus string

const (
	DataActionStatus_Triggered DataActionStatus = "TRIGGERED"
	DataActionStatus_Pending   DataActionStatus = "PENDING"
)

type DataAction struct {
	ID               int64 `gorm:"primaryKey"`
	ActionType       ActionType
	Payload          pqtype.NullRawMessage
	Status           DataActionStatus
	RunCount         int64
	Schedule         string
	DagId            string
	SourceTableMapId int64
	AccountUuid      uuid.UUID
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

func (DataAction) TableName() string {
	return "data_action"
}

package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type ActionType string

const (
	ActionType_UploadDataFromFile ActionType = "IMPORT_DATA_FROM_FILE"
)

type DataActionStatus string

const (
	DataActionStatus_Done     = "DONE"
	DataActionStatus_Pending  = "PENDING"
	DataActionStatus_Failed   = "FAILED"
	DataActionStatus_Canceled = "CANCELED"
)

type DataAction struct {
	ID          int64 `gorm:"primaryKey"`
	ActionType  ActionType
	Payload     pqtype.NullRawMessage
	Status      DataActionStatus
	RunCount    int
	Schedule    string
	DagId       string
	AccountUuid uuid.UUID
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataAction) TableName() string {
	return "data_action"
}

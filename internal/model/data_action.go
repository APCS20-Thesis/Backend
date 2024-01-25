package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type EventType string

const (
	EventType_UploadDataFromFile EventType = "IMPORT_DATA_FROM_FILE"
)

type DataAction struct {
	ID             int64 `gorm:"primaryKey"`
	EventType      EventType
	Payload        pqtype.NullRawMessage
	Configuration  string
	MappingOptions string
	DeltaTableName string
	AccountUuid    uuid.UUID
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (DataAction) TableName() string {
	return "data_action"
}

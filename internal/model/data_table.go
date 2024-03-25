package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type TableStatus string

const (
	TableStatus_DRAFT        TableStatus = "DRAFT"
	TableStatus_UPDATING     TableStatus = "UPDATING"
	TableStatus_NEED_TO_SYNC TableStatus = "NEED_TO_SYNC"
	TableStatus_UP_TO_DATE   TableStatus = "UP_TO_DATE"
)

type DataTable struct {
	ID          int64 `gorm:"primaryKey"`
	Schema      pqtype.NullRawMessage
	Name        string
	AccountUuid uuid.UUID
	Status      TableStatus
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataTable) TableName() string {
	return "data_table"
}

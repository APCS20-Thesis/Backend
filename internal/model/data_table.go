package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type DataTableStatus string

const (
	DataTableStatus_DRAFT        DataTableStatus = "DRAFT"
	DataTableStatus_UPDATING     DataTableStatus = "UPDATING"
	DataTableStatus_NEED_TO_SYNC DataTableStatus = "NEED_TO_SYNC"
	DataTableStatus_UP_TO_DATE   DataTableStatus = "UP_TO_DATE"
)

type DataTable struct {
	ID          int64 `gorm:"primaryKey"`
	Schema      pqtype.NullRawMessage
	Name        string
	AccountUuid uuid.UUID
	Status      DataTableStatus
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataTable) TableName() string {
	return "data_table"
}

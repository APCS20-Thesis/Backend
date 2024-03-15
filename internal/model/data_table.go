package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type DataTable struct {
	ID          int64 `gorm:"primaryKey"`
	Schema      pqtype.NullRawMessage
	Name        string
	AccountUuid uuid.UUID
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataTable) TableName() string {
	return "data_table"
}

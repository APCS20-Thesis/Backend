package model

import (
	"github.com/sqlc-dev/pqtype"
	"time"
)

type DestTableMap struct {
	ID             int64 `gorm:"primaryKey"`
	TableId        int64
	DestinationId  int64
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	MappingOptions pqtype.NullRawMessage
}

func (DestTableMap) TableName() string {
	return "dest_table_map"
}

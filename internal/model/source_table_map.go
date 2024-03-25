package model

import (
	"github.com/sqlc-dev/pqtype"
	"time"
)

type SourceTableMap struct {
	ID                int64 `gorm:"primaryKey"`
	TableId           int64
	SourceId          int64
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
	MappingOptions    pqtype.NullRawMessage
	TableNameInSource string
}

func (SourceTableMap) TableName() string {
	return "source_table_map"
}

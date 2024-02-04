package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type DataSourceType string

const (
	DataSourceType_File  DataSourceType = "FILE"
	DataSourceType_MySQL DataSourceType = "MYSQL"
)

type DataSource struct {
	ID             int64 `gorm:"primaryKey"`
	Name           string
	Description    string
	Type           DataSourceType
	Configuration  pqtype.NullRawMessage
	MappingOptions pqtype.NullRawMessage
	DeltaTableName string
	AccountUuid    uuid.UUID
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (DataSource) TableName() string {
	return "data_source"
}

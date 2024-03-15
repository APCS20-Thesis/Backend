package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type ConnectionType string

const (
	ConnectionType_S3    DataSourceType = "AWS-S3"
	ConnectionType_MySQL DataSourceType = "MYSQL"
)

type SourceConnection struct {
	ID             int64 `gorm:"primaryKey"`
	Name           string
	Type           ConnectionType
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (SourceConnection) TableName() string {
	return "source_connection"
}

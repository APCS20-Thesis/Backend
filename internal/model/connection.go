package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type ConnectionType string

const (
	ConnectionType_S3    ConnectionType = "AWS-S3"
	ConnectionType_MySQL ConnectionType = "MYSQL"
)

type Connection struct {
	ID             int64 `gorm:"primaryKey"`
	Name           string
	Type           ConnectionType
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type (
	S3Configurations struct {
		AccessKeyId     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		BucketName      string `json:"bucket_name"`
		Region          string `json:"region"`
	}
)

func (Connection) TableName() string {
	return "connection"
}

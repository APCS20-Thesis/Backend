package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type ConnectionType string

const (
	ConnectionType_S3      ConnectionType = "AWS-S3"
	ConnectionType_MySQL   ConnectionType = "MYSQL"
	ConnectionType_Gophish ConnectionType = "GOPHISH"
)

var ConnectionTypes = []ConnectionType{ConnectionType_S3, ConnectionType_MySQL, ConnectionType_Gophish}

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

	GophishConfiguration struct {
		Host   string `json:"host"`
		Port   string `json:"port"`
		ApiKey string `json:"api_key"`
	}

	MySQLConfiguration struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
)

func (Connection) TableName() string {
	return "connection"
}

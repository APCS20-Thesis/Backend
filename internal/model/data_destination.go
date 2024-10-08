package model

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type (
	DataDestinationType string
	//DataDestinationStatus string
)

const (
	DataDestinationType_GOPHISH   DataDestinationType = "GOPHISH"
	DataDestinationType_MYSQL     DataDestinationType = "MYSQL"
	DataDestinationType_S3FileCSV DataDestinationType = "S3-CSV"

	//DataDestinationStatus_Success DataDestinationStatus = "SUCCESS"
)

type DataDestination struct {
	ID             int64 `gorm:"primaryKey"`
	Name           string
	Type           DataDestinationType
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
	ConnectionId   int64
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (DataDestination) TableName() string {
	return "data_destination"
}

type (
	GophishDestinationConfiguration struct {
		UserGroupName string                     `json:"user_group_name"`
		Mapping       *api.MappingGophishProfile `json:"mapping"`
	}

	S3FileDestinationConfiguration struct {
		FileName string `json:"file_name"`
		FileType string `json:"file_type"`
		FilePath string `json:"file_path"`
	}

	MySQLDestinationConfiguration struct {
		TableName string `json:"table_name"`
	}
)

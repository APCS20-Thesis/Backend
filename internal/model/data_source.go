package model

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

const (
	DataSourceType_FileCsv   DataSourceType = "CSV"
	DataSourceType_FileExcel DataSourceType = "EXCEL"
	DataSourceType_MySQL     DataSourceType = "MYSQL"

	DataSourceStatus_Processing DataSourceStatus = "PROCESSING"
	DataSourceStatus_Success    DataSourceStatus = "SUCCESS"
	DataSourceStatus_Failed     DataActionStatus = "FAILED"
)

type DataSource struct {
	ID             int64 `gorm:"primaryKey"`
	Name           string
	Description    string
	Type           DataSourceType
	Status         DataSourceStatus
	Configurations pqtype.NullRawMessage
	AccountUuid    uuid.UUID
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type (
	CsvConfigurations struct {
		FileName      string                                        `json:"file_name"`
		ConnectionId  int64                                         `json:"connection_id"`
		Key           string                                        `json:"key"`
		CsvReadOption *api.ImportCsvRequest_ImportCsvConfigurations `json:"csv_read_option"`
	}

	DataSourceType string

	DataSourceStatus string
)

func (DataSource) TableName() string {
	return "data_source"
}

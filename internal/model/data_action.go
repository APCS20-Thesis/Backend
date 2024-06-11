package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type ActionType string

const (
	ActionType_ImportDataFromFile  ActionType = "IMPORT_DATA_FROM_FILE"
	ActionType_ImportDataFromS3    ActionType = "IMPORT_DATA_FROM_S3"
	ActionType_ImportDataFromMySQL ActionType = "IMPORT_DATA_FROM_MYSQL"
	ActionType_ExportDataToCSV     ActionType = "EXPORT_DATA_TO_CSV"
	ActionType_ExportToMySQL       ActionType = "EXPORT_TABLE_TO_MYSQL"
	ActionType_CreateMasterSegment ActionType = "CREATE_MS_SEGMENT"
)

type DataActionTargetTable string

const (
	TargetTable_SourceTableMap DataActionTargetTable = "source_table_map"
	TargetTable_DestTableMap   DataActionTargetTable = "dest_table_map"
	TargetTable_DataTable      DataActionTargetTable = "data_table"
)

type DataActionStatus string

const (
	DataActionStatus_Triggered  DataActionStatus = "TRIGGERED"
	DataActionStatus_Pending    DataActionStatus = "PENDING"
	DataActionStatus_Processing DataActionStatus = "PROCESSING"
	DataActionStatus_Success    DataActionStatus = "SUCCESS"
	DataActionStatus_Failed     DataActionStatus = "FAILED"
)

type DataAction struct {
	ID          int64 `gorm:"primaryKey"`
	ActionType  ActionType
	Payload     pqtype.NullRawMessage
	Status      DataActionStatus
	RunCount    int64
	Schedule    string
	DagId       string
	TargetTable DataActionTargetTable
	ObjectId    int64
	AccountUuid uuid.UUID
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (DataAction) TableName() string {
	return "data_action"
}

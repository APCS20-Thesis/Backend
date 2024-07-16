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
	ActionType_ExportDataToS3CSV   ActionType = "EXPORT_DATA_TO_S3_CSV"
	ActionType_ExportToMySQL       ActionType = "EXPORT_TABLE_TO_MYSQL"
	ActionType_CreateMasterSegment ActionType = "CREATE_MS_SEGMENT"
	ActionType_CreateSegment       ActionType = "CREATE_SEGMENT"
	ActionType_TrainPredictModel   ActionType = "TRAIN_PREDICT_MODEL"
	ActionType_ApplyPredictModel   ActionType = "APPLY_PREDICT_MODEL"
)

type DataActionTargetTable string

const (
	TargetTable_SourceTableMap       DataActionTargetTable = "source_table_map"
	TargetTable_DestTableMap         DataActionTargetTable = "dest_table_map"
	TargetTable_DestSegmentMap       DataActionTargetTable = "dest_segment_map"
	TargetTable_DestMasterSegmentMap DataActionTargetTable = "dest_ms_segment_map"
	TargetTable_DataTable            DataActionTargetTable = "data_table"
	TargetTable_Segment              DataActionTargetTable = "segment"
	TargetTable_PredictModel         DataActionTargetTable = "predict_model"
)

type DataActionStatus string

const (
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

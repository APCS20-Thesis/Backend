package utils

import (
	"fmt"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"strconv"
	"strings"
	"time"
)

const (
	ExportDataFileLocationFormat_AccountUuid            = "account_uuid"
	ExportDataFileLocationFormat_DateTimeFormatYYYYMMDD = "datetime"
	ExportDataFileLocationFormat_FileName               = "file_name"
	ExportDataFileLocationFormat_FileFormat             = "file_format"
	ExportDataFileLocationFormat                        = "export/account_uuid/datetime/file_name.file_format"

	CreateMasterSegmentFormat_DagId           = "create_ms_segment_account_uuid_master_segment_id"
	Common_AccountUuid                        = "account_uuid"
	CreateMasterSegmentFormat_MasterSegmentId = "master_segment_id"

	QueryDeltaTable_Path      = "data/bronze/account_uuid/table_name"
	QueryDeltaTable_TableName = "table_name"

	DeltaAudience_Path            = "segments/id/id_audience"
	DeltaAudience_MasterSegmentId = "id"

	DeltaBehavior_Path            = "segments/id/id_name"
	DeltaBehavior_MasterSegmentId = "id"
	DeltaBehavior_TableName       = "name"
)

func GenerateExportDataFileLocation(accountUuid string, fileName string, fileFormat string) string {
	dateTime := time.Now().UTC().Format("2006-01-02")
	path := strings.Replace(ExportDataFileLocationFormat, ExportDataFileLocationFormat_AccountUuid, accountUuid, 1)
	path = strings.Replace(path, ExportDataFileLocationFormat_DateTimeFormatYYYYMMDD, dateTime, 1)
	path = strings.Replace(path, ExportDataFileLocationFormat_FileName, fileName, 1)
	path = strings.Replace(path, ExportDataFileLocationFormat_FileFormat, fileFormat, 1)

	return path
}

func GenerateDagIdForCreateMasterSegment(accountUuid string, masterSegmentId int64) string {
	path := strings.Replace(CreateMasterSegmentFormat_DagId, Common_AccountUuid, accountUuid, 1)
	path = strings.Replace(path, CreateMasterSegmentFormat_MasterSegmentId, strconv.FormatInt(masterSegmentId, 10), 1)

	return path
}

func GenerateDeltaTablePath(accountUuid string, tableName string) string {
	path := strings.Replace(QueryDeltaTable_Path, Common_AccountUuid, accountUuid, 1)
	path = strings.Replace(path, QueryDeltaTable_TableName, tableName, 1)
	return path
}

func GenerateDeltaAudiencePath(masterSegmentId int64) string {
	return strings.Replace(DeltaAudience_Path, DeltaAudience_MasterSegmentId,
		strconv.FormatInt(masterSegmentId, 10), 2)
}

func GenerateDeltaBehaviorPath(masterSegmentId int64, behaviorTableName string) string {
	path := strings.Replace(DeltaBehavior_Path, DeltaBehavior_MasterSegmentId,
		strconv.FormatInt(masterSegmentId, 10), 2)
	path = strings.Replace(path, DeltaBehavior_TableName, behaviorTableName, 1)
	return path
}

func GenerateDagId(accountUuid string, dataActionType model.ActionType) string {
	dateTime := strconv.FormatInt(time.Now().Unix(), 10)
	var prefix string
	switch dataActionType {
	case model.ActionType_ImportDataFromS3:
		prefix = "import_csv_s3"
	case model.ActionType_ImportDataFromMySQL:
		prefix = "import_mysql"
	case model.ActionType_ImportDataFromFile:
		prefix = "import_csv"
	case model.ActionType_ExportToMySQL:
		prefix = "export_mysql"
	case model.ActionType_ExportDataToS3CSV:
		prefix = "export_csv"
	case model.ActionType_CreateSegment:
		prefix = "create_segment"
	case model.ActionType_TrainPredictModel:
		prefix = "train_predict_model"
	default:
		return ""
	}
	return fmt.Sprintf("%s_%s_%s", prefix, accountUuid, dateTime)
}

func GenerateDeltaSegmentPath(masterSegmentId int64, segmentId int64) string {
	return fmt.Sprintf("segments/%d/%d_segment", masterSegmentId, segmentId)
}

func GenerateDeltaPredictModelFilePath(masterSegmentId int64, predictModelId int64) string {
	return fmt.Sprintf("segments/%d/%d_model", masterSegmentId, predictModelId)
}

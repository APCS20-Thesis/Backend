package utils

import (
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
	CreateMasterSegmentFormat_AccountUuid     = "account_uuid"
	CreateMasterSegmentFormat_MasterSegmentId = "master_segment_id"
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
	path := strings.Replace(CreateMasterSegmentFormat_DagId, CreateMasterSegmentFormat_AccountUuid, accountUuid, 1)
	path = strings.Replace(path, CreateMasterSegmentFormat_MasterSegmentId, strconv.FormatInt(masterSegmentId, 10), 1)

	return path
}

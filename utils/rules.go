package utils

import (
	"fmt"
	"github.com/APCS20-Thesis/Backend/api"
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

func BuildSQLCondition(condition []*api.CreateSegmentRequest_GroupCondition) string {
	if len(condition) == 0 {
		return ""
	}

	sqlCondition := ""
	for _, each := range condition {
		clause := fmt.Sprintf("%s %s %s", each.Condition.ColumnName, each.Condition.Operator, each.Condition.Value)
		if len(each.Condition.Condition) > 0 {
			for _, alternateCondition := range each.Condition.Condition {
				clause += fmt.Sprintf(" %s %s %s %s",
					alternateCondition.Combinator,
					alternateCondition.Condition.ColumnName,
					alternateCondition.Condition.Operator,
					alternateCondition.Condition.Value,
				)
			}
			if each.Combinator == "" {
				sqlCondition += fmt.Sprintf("(%s)", clause)
				continue
			}
			sqlCondition += fmt.Sprintf(" %s (%s)", each.Combinator, clause)
		} else {
			if each.Combinator == "" {
				sqlCondition += fmt.Sprintf("%s", clause)
				continue
			}
			sqlCondition += fmt.Sprintf(" %s %s", each.Combinator, clause)
		}
	}

	return sqlCondition
}

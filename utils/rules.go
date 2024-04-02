package utils

import (
	"strings"
	"time"
)

const (
	ExportDataFileLocationFormat_AccountUuid            = "account_uuid"
	ExportDataFileLocationFormat_DateTimeFormatYYYYMMDD = "datetime"
	ExportDataFileLocationFormat_FileName               = "file_name"
	ExportDataFileLocationFormat_FileFormat             = "file_format"
	ExportDataFileLocationFormat                        = "export/account_uuid/datetime/file_name.file_format"
)

func GenerateExportDataFileLocation(accountUuid string, fileName string, fileFormat string) string {
	dateTime := time.Now().UTC().Format("2006-01-02")
	path := strings.Replace(ExportDataFileLocationFormat, ExportDataFileLocationFormat_AccountUuid, accountUuid, 1)
	path = strings.Replace(path, ExportDataFileLocationFormat_DateTimeFormatYYYYMMDD, dateTime, 1)
	path = strings.Replace(path, ExportDataFileLocationFormat_FileName, fileName, 1)
	path = strings.Replace(path, ExportDataFileLocationFormat_FileFormat, fileFormat, 1)

	return path
}

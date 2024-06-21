package config

func AccessibleRoles() map[string][]string {
	const rootServicePath = "/api.CDPService/"
	const rootServiceFilePath = "/api.CDPServiceFile/"
	return map[string][]string{
		rootServicePath + "Admin":                             {"admin"},
		rootServicePath + "GetAccountInfo":                    {"admin", "user"},
		rootServiceFilePath + "ImportCsv":                     {"admin", "user"},
		rootServicePath + "GetListDataSources":                {"admin", "user"},
		rootServicePath + "GetDataSource":                     {"admin", "user"},
		rootServicePath + "GetListDataTables":                 {"admin", "user"},
		rootServicePath + "GetDataTable":                      {"admin", "user"},
		rootServicePath + "GetListConnections":                {"admin", "user"},
		rootServicePath + "GetConnection":                     {"admin", "user"},
		rootServicePath + "CreateConnection":                  {"admin", "user"},
		rootServicePath + "UpdateConnection":                  {"admin", "user"},
		rootServicePath + "DeleteConnection":                  {"admin", "user"},
		rootServicePath + "ExportDataToFile":                  {"admin", "user"},
		rootServicePath + "GetListFileExportRecords":          {"admin", "user"},
		rootServicePath + "ImportCsvFromS3":                   {"admin", "user"},
		rootServicePath + "CreateMasterSegment":               {"admin", "user"},
		rootServicePath + "GetListMasterSegments":             {"admin", "user"},
		rootServicePath + "GetMasterSegmentDetail":            {"admin", "user"},
		rootServicePath + "CreateSegment":                     {"admin", "user"},
		rootServicePath + "GetListSegments":                   {"admin", "user"},
		rootServicePath + "GetQueryDataTable":                 {"admin", "user"},
		rootServicePath + "GetSegmentDetail":                  {"admin", "user"},
		rootServicePath + "CreateGophishUserGroupFromSegment": {"admin", "user"},
		rootServicePath + "ImportFromMySQLSource":             {"admin", "user"},
		rootServicePath + "ExportToMySQLDestination":          {"admin", "user"},
		rootServicePath + "GetListDataDestinations":           {"admin", "user"},
		rootServicePath + "GetListDataActions":                {"admin", "user"},
	}
}

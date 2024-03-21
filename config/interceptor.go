package config

func AccessibleRoles() map[string][]string {
	const rootServicePath = "/api.CDPService/"
	const rootServiceFilePath = "/api.CDPServiceFile/"
	return map[string][]string{
		rootServicePath + "Admin":              {"admin"},
		rootServicePath + "GetAccountInfo":     {"admin", "user"},
		rootServiceFilePath + "ImportFile":     {"admin", "user"},
		rootServicePath + "GetListDataSources": {"admin", "user"},
		rootServicePath + "GetDataSource":      {"admin", "user"},
		rootServicePath + "GetListDataTables":  {"admin", "user"},
		rootServicePath + "GetDataTable":       {"admin", "user"},
		rootServicePath + "GetListConnections": {"admin", "user"},
		rootServicePath + "GetConnection":      {"admin", "user"},
		rootServicePath + "CreateConnection":   {"admin", "user"},
		rootServicePath + "UpdateConnection":   {"admin", "user"},
		rootServicePath + "DeleteConnection":   {"admin", "user"},
	}
}

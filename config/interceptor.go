package config

func AccessibleRoles() map[string][]string {
	const rootServicePath = "/api.CDPService/"
	return map[string][]string{
		rootServicePath + "Admin":                 {"admin"},
		rootServicePath + "GetAccountInfo":        {"admin", "user"},
		rootServicePath + "CreateDataSourceMySQL": {"admin", "user"},
	}
}

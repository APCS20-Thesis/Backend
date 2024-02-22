package config

func AccessibleRoles() map[string][]string {
	const rootServicePath = "/api.CDPService/"
	const rootServiceFilePath = "/api.CDPServiceFile/"
	return map[string][]string{
		rootServicePath + "Admin":          {"admin"},
		rootServicePath + "GetAccountInfo": {"admin", "user"},
		rootServiceFilePath + "ImportFile": {"admin", "user"},
	}
}

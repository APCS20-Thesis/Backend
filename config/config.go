package config

type Config struct {
	Base         `mapstructure:",squash"`
	ServerConfig ServerConfig `json:"server_config" mapstructure:"server_config"`

	PostgreSQL      PostgreSQLConfig `json:"postgresql" mapstructure:"postgresql"`
	MigrationFolder string           `json:"migration_folder" mapstructure:"migration_folder"`
	S3StorageConfig S3StorageConfig  `json:"s3_storage_config" mapstructure:"s3_storage_config"`
}

type Base struct {
	Env string    `json:"env" mapstructure:"env"`
	Log LogConfig `json:"log" mapstructure:"log"`
}

func loadDefaultConfig() *Config {
	c := &Config{
		ServerConfig: ServerConfig{
			HttpServerAddress: ":11080",
			GrpcServerAddress: ":10443",
		},
		PostgreSQL: PostgreSQLConfig{
			DBConfig: DBConfig{
				Host:     "127.0.0.1",
				Database: "cdp_service",
				Port:     5432,
				Username: "cdp_service",
				Password: "postgres",
				Options:  "?sslmode=disable",
			},
		},
		MigrationFolder: "file://sql/migrations",
		S3StorageConfig: S3StorageConfig{
			Host:            "https://cdp-thesis-apcs.s3.ap-southeast-1.amazonaws.com",
			AccessKeyID:     "AKIASPHW355ITDXOYKO6",
			SecretAccessKey: "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",
			Region:          "ap-southeast-1",
		},
	}
	return c
}

type ServerConfig struct {
	HttpServerAddress string `json:"http_server_address" mapstructure:"http_server_address"`
	GrpcServerAddress string `json:"grpc_server_address" mapstructure:"grpc_server_address"`
}

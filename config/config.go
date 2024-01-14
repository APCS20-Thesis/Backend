package config

type Config struct {
	Base         `mapstructure:",squash"`
	ServerConfig ServerConfig     `json:"server_config" mapstructure:"server_config"`
	PostgreSQL   PostgreSQLConfig `json:"postgresql" mapstructure:"postgresql"`
}

type Base struct {
	Env string    `json:"env" mapstructure:"env"`
	Log LogConfig `json:"log" mapstructure:"log"`
}

func loadDefaultConfig() *Config {
	c := &Config{
		ServerConfig: ServerConfig{
			HttpServerAddress: ":10080",
			GrpcServerAddress: ":10443",
		},
		PostgreSQL: PostgreSQLConfig{},
	}
	return c
}

type ServerConfig struct {
	HttpServerAddress string `json:"http_server_address" mapstructure:"http_server_address"`
	GrpcServerAddress string `json:"grpc_server_address" mapstructure:"grpc_server_address"`
}

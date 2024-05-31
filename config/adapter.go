package config

type AirflowConfig struct {
	Address  string `json:"address" mapstructure:"address"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

type QueryConfig struct {
	Address string `json:"address" mapstructure:"address"`
}

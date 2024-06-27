package config

type AirflowConfig struct {
	Address  string `json:"address" mapstructure:"address"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

type QueryConfig struct {
	Address string `json:"address" mapstructure:"address"`
}

type MqttConfig struct {
	Host string `json:"host" mapstructure:"host"`
	Port int32  `json:"port" mapstructure:"port"`
}

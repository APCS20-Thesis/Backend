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
	Host     string `json:"host" mapstructure:"host"`
	Port     int32  `json:"port" mapstructure:"port"`
	ClientID string `json:"client_id" mapstructure:"client_id"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

type AlertConfig struct {
	EnableAlert bool   `json:"enable_alert" mapstructure:"enable_alert"`
	Webhook     string `json:"webhook" mapstructure:"webhook"`
}

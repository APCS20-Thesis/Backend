package config

type S3StorageConfig struct {
	Host            string `json:"host" mapstructure:"host" yaml:"host"`
	AccessKeyID     string `json:"access_key_id" mapstructure:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" mapstructure:"secret_access_key" yaml:"secret_access_key"`
	Region          string `json:"region" mapstructure:"region" yaml:"region"`
}

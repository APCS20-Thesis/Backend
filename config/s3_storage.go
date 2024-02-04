package config

type S3StorageConfig struct {
	Host            string `json:"host" mapstructure:"host" yaml:"host"`
	AccessKeyID     string `json:"access_key_id" mapstructure:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" mapstructure:"secret_access_key" yaml:"secret_access_key"`
}

// PostgresSQLDefaultConfig returns default config for mysql, usually use on development.
func S3StorageDefaultConfig() S3StorageConfig {
	return S3StorageConfig{
		Host:            "localhost:4566",
		AccessKeyID:     "foo",
		SecretAccessKey: "bar",
	}
}

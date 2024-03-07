package config

import (
	"fmt"
	"net/url"
)

type DBConfig struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstructure:"options" yaml:"options"`
}

// DBConfig used to set config for database.
type IDBConfig interface {
	String() string
	DSN() string
}

// DSN returns the Domain Source Name.
func (c DBConfig) DSN() string {
	options := c.Options
	if options != "" {
		if options[0] != '?' {
			options = "?" + options
		}
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		c.Username,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Database,
		options)
}

// PostgreSQLConfig used to set config for Postgres.
type PostgreSQLConfig struct {
	DBConfig `mapstructure:",squash"`
}

func (c PostgreSQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@%s:%d/%s%s", c.Username, url.QueryEscape(c.Password), c.Host, c.Port, c.Database, c.Options)
}

// String returns Postgres connection URI.
func (c PostgreSQLConfig) String() string {
	return fmt.Sprintf("postgresql://%s", c.DSN())
}

// PostgresSQLDefaultConfig returns default config for mysql, usually use on development.
func PostgresSQLDefaultConfig() PostgreSQLConfig {
	return PostgreSQLConfig{DBConfig{
		Host:     "127.0.0.1",
		Port:     5433,
		Database: "sample",
		Username: "default",
		Password: "secret",
		Options:  "",
	}}
}

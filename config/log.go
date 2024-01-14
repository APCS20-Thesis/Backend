package config

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// Config holding log config
type LogConfig struct {
	Level    string `json:"level" mapstructure:"level"`
	Mode     string `json:"mode" mapstructure:"mode"`
	Encoding string `json:"encoding" mapstructure:"encoding"`
}

func (c LogConfig) MustBuildLogR() logr.Logger {
	var log logr.Logger

	zapConfig := zap.NewDevelopmentConfig()
	if c.Mode == "production" {
		zapConfig = zap.NewProductionConfig()
	}
	zapLog, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	return log
}

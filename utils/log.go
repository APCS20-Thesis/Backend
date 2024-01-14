package utils

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// New Log
func NewLog() logr.Logger {
	log := zapr.NewLogger(zap.NewNop())
	return log
}

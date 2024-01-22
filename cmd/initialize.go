package main

import (
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logr "gorm.io/gorm/logger"
	oslog "log"
	"os"
	"time"
)

const (
	DefaultMaxIdlesConst = 10
	DefaultMaxOpenConst  = 100
	DefaultSlowThreshold = 10 * time.Second
)

func newDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}
	// force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

//nolint:gomnd
func ConnectPostgresql(dsn string) (*gorm.DB, error) {
	newLogger := gorm_logr.New(
		oslog.New(os.Stderr, "", oslog.LstdFlags),
		gorm_logr.Config{
			SlowThreshold: DefaultSlowThreshold,
			LogLevel:      gorm_logr.Warn,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(DefaultMaxIdlesConst)
	sqlDB.SetMaxOpenConns(DefaultMaxOpenConst)

	err = db.Raw("SELECT 1").Error
	if err != nil {
		return nil, err
	}

	return db, nil
}

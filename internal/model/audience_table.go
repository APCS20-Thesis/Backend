package model

import (
	"github.com/sqlc-dev/pqtype"
	"time"
)

type AudienceTable struct {
	ID                 int64 `gorm:"primaryKey"`
	MasterSegmentId    int64
	Schema             pqtype.NullRawMessage
	BuildConfiguration pqtype.NullRawMessage
	Name               string
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (AudienceTable) TableName() string {
	return "audience_table"
}

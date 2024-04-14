package model

import (
	"github.com/sqlc-dev/pqtype"
	"time"
)

type BehaviorTable struct {
	ID                 int64 `gorm:"primaryKey"`
	MasterSegmentId    int64
	Schema             pqtype.NullRawMessage
	DataTableId        int64
	AudienceForeignKey string
	JoinKey            string
	Name               string
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (BehaviorTable) TableName() string {
	return "behavior_table"
}

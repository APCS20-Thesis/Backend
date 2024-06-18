package model

import (
	"github.com/sqlc-dev/pqtype"
	"time"
)

type DestMasterSegmentMap struct {
	ID              int64 `gorm:"primaryKey"`
	MasterSegmentId int64
	DestinationId   int64
	MappingOptions  pqtype.NullRawMessage
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (DestMasterSegmentMap) TableName() string {
	return "dest_ms_segment_map"
}

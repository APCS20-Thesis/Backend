package model

import (
	"github.com/sqlc-dev/pqtype"
	"time"
)

type DestSegmentMap struct {
	ID             int64 `gorm:"primaryKey"`
	SegmentId      int64
	DestinationId  int64
	MappingOptions pqtype.NullRawMessage
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (DestSegmentMap) TableName() string {
	return "dest_segment_map"
}

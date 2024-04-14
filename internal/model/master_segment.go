package model

import (
	"github.com/google/uuid"
	"time"
)

type MasterSegmentStatus string

const (
	MasterSegmentStatus_DRAFT      MasterSegmentStatus = "DRAFT"
	MasterSegmentStatus_UPDATING   MasterSegmentStatus = "UPDATING"
	MasterSegmentStatus_UP_TO_DATE MasterSegmentStatus = "UP_TO_DATE"
)

type MasterSegment struct {
	ID          int64 `gorm:"primaryKey"`
	Description string
	Name        string
	AccountUuid uuid.UUID
	Status      MasterSegmentStatus
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (MasterSegment) TableName() string {
	return "master_segment"
}

package model

import (
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type Segment struct {
	ID              int64 `gorm:"primaryKey"`
	MasterSegmentId int64
	Condition       pqtype.NullRawMessage
	Description     string
	Name            string
	AccountUuid     uuid.UUID
	Status          MasterSegmentStatus
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (Segment) TableName() string {
	return "segment"
}

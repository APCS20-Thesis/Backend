package model

import (
	"github.com/APCS20-Thesis/Backend/api"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

type Segment struct {
	ID              int64 `gorm:"primaryKey"`
	MasterSegmentId int64
	Condition       pqtype.NullRawMessage
	SqlCondition    string
	Description     string
	Name            string
	AccountUuid     uuid.UUID
	Status          SegmentStatus
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (Segment) TableName() string {
	return "segment"
}

type SegmentBuildConditions struct {
	AudienceCondition  *api.Rule                `json:"audience_condition"`
	BehaviorConditions []*api.BehaviorCondition `json:"behavior_conditions"`
}

type SegmentStatus string

const (
	SegmentStatus_DRAFT      SegmentStatus = "DRAFT"
	SegmentStatus_UPDATING   SegmentStatus = "UPDATING"
	SegmentStatus_UP_TO_DATE SegmentStatus = "UP_TO_DATE"
	SegmentStatus_FAILED     SegmentStatus = "FAILED"
)

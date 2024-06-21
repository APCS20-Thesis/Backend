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
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (Segment) TableName() string {
	return "segment"
}

type SegmentBuildConditions struct {
	AudienceCondition  *api.Rule
	BehaviorConditions []*api.BehaviorCondition
}

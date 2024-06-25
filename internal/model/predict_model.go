package model

import (
	"github.com/sqlc-dev/pqtype"
	"time"
)

type PredictModelStatus string

const (
	PredictModelStatus_DRAFT      PredictModelStatus = "DRAFT"
	PredictModelStatus_UP_TO_DATE PredictModelStatus = "UP_TO_DATE"
)

type PredictModel struct {
	ID                  int64 `gorm:"primaryKey"`
	Name                string
	MasterSegmentId     int64
	TrainConfigurations pqtype.NullRawMessage
	Status              PredictModelStatus
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}

func (PredictModel) TableName() string {
	return "predict_model"
}

type PredictModelTrainConfiguration struct {
	Segment1           int64  `json:"segment_1"`
	Segment2           int64  `json:"segment_2"`
	Label1             string `json:"label_1"`
	Label2             string `json:"label_2"`
	SelectedAttributes []string
}

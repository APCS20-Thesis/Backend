package model

import "time"

type SourceTableMap struct {
	ID        int64 `gorm:"primaryKey"`
	TableId   int64
	SourceId  int64
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (SourceTableMap) TableName() string {
	return "source_table_map"
}

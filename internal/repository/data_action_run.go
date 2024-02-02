package repository

import (
	"gorm.io/gorm"

	"github.com/APCS20-Thesis/Backend/internal/model"
)

type DataActionRunRepository interface {
}

type dataActionRunRepo struct {
	*gorm.DB
	TableName string
}

func NewDataActionRunRepository(db *gorm.DB) DataActionRunRepository {
	return &dataActionRunRepo{db, model.DataActionRun{}.TableName()}
}

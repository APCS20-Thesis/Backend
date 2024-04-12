package repository

import (
	"context"
	"github.com/APCS20-Thesis/Backend/internal/model"
	"gorm.io/gorm"
)

type FileExportRecordRepository interface {
	ListFileExportRecords(ctx context.Context, tableId int64, accountUuid string) ([]model.FileExportRecord, error)
}

type fileExportRecordRepo struct {
	*gorm.DB
	TableName string
}

func NewFileExportRecordRepository(db *gorm.DB) FileExportRecordRepository {
	return &fileExportRecordRepo{db, model.FileExportRecord{}.TableName()}
}

func (r *fileExportRecordRepo) ListFileExportRecords(ctx context.Context, tableId int64, accountUuid string) ([]model.FileExportRecord, error) {
	var records []model.FileExportRecord
	err := r.WithContext(ctx).Table(r.TableName).
		Where("data_table_id = ? AND account_uuid = ?", tableId, accountUuid).
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

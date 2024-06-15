package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	AccountRepository
	DataSourceRepository
	DataActionRepository
	DataActionRunRepository
	DataTableRepository
	ConnectionRepository
	TransactionRepository
	FileExportRecordRepository
	SourceTableMapRepository
	SegmentRepository
	DataDestinationRepository
	DestTableMapRepository
	DestSegmentMapRepository
	DestMasterSegmentMapRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AccountRepository:              NewAccountRepository(db),
		DataSourceRepository:           NewDataSourceRepository(db),
		DataActionRepository:           NewDataActionRepository(db),
		DataActionRunRepository:        NewDataActionRunRepository(db),
		DataTableRepository:            NewDataTableRepository(db),
		ConnectionRepository:           NewConnectionRepository(db),
		TransactionRepository:          NewTransactionRepository(db),
		FileExportRecordRepository:     NewFileExportRecordRepository(db),
		SourceTableMapRepository:       NewSourceTableMapRepository(db),
		SegmentRepository:              NewSegmentRepository(db),
		DataDestinationRepository:      NewDataDestinationRepository(db),
		DestTableMapRepository:         NewDestTableMapRepository(db),
		DestSegmentMapRepository:       NewDestSegmentMapRepository(db),
		DestMasterSegmentMapRepository: NewDestMasterSegmentMapRepository(db),
	}
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

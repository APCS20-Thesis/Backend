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
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AccountRepository:          NewAccountRepository(db),
		DataSourceRepository:       NewDataSourceRepository(db),
		DataActionRepository:       NewDataActionRepository(db),
		DataActionRunRepository:    NewDataActionRunRepository(db),
		DataTableRepository:        NewDataTableRepository(db),
		ConnectionRepository:       NewConnectionRepository(db),
		TransactionRepository:      NewTransactionRepository(db),
		FileExportRecordRepository: NewFileExportRecordRepository(db),
		SourceTableMapRepository:   NewSourceTableMapRepository(db),
		SegmentRepository:          NewSegmentRepository(db),
		DataDestinationRepository:  NewDataDestinationRepository(db),
	}
}

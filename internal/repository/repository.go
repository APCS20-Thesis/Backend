package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	AccountRepository
	DataSourceRepository
	DataActionRepository
	DataActionRunRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AccountRepository:       NewAccountRepository(db),
		DataSourceRepository:    NewDataSourceRepository(db),
		DataActionRepository:    NewDataActionRepository(db),
		DataActionRunRepository: NewDataActionRunRepository(db),
	}
}

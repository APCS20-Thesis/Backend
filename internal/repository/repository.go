package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	AccountRepository
	DataSourceRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		AccountRepository:    NewAccountRepository(db),
		DataSourceRepository: NewDataSourceRepository(db),
	}
}

package repository

import "gorm.io/gorm"

type Repository struct {
	AccountRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		NewAccountRepository(db),
	}
}

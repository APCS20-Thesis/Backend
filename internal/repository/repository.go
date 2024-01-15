package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		NewUserRepository(db),
	}
}

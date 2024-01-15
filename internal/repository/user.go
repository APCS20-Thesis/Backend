package repository

import (
	"github.com/APCS20-Thesis/Backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
}

type userRepo struct {
	*gorm.DB
	TableName string
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db, model.User{}.TableName()}
}

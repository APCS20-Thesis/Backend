package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        int64     `gorm:"primary_key"`
	Uuid      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username  string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	FirstName string
	LastName  string
	Email     string `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "user"
}

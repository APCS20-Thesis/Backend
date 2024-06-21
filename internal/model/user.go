package model

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID        int64     `gorm:"primary_key"`
	Uuid      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username  string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	FirstName string
	LastName  string
	Email     string `gorm:"not null"`

	Phone    string
	Country  string
	Company  string
	Position string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Account) TableName() string {
	return "account"
}

type Setting struct {
	NotifyCreateSource        bool
	NotifyCreateDestination   bool
	NotifyCreateMasterSegment bool
	NotifyCreateSegment       bool

	AccountUuid uuid.UUID

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Setting) TableName() string {
	return "setting"
}

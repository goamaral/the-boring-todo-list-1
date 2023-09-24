package gorm_provider

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type AbstractEntity interface {
	GetID() uint
}

type Entity struct {
	ID        uint `gorm:"primary_key" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e Entity) GetID() uint {
	return e.ID
}

type EntityWithUUID struct {
	ID        uint `gorm:"primary_key"`
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e EntityWithUUID) GetID() uint {
	return e.ID
}

func (e *EntityWithUUID) BeforeCreate(tx *gorm.DB) error {
	e.UUID = ulid.Make().String()
	return nil
}

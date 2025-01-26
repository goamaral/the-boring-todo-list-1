package gorm_provider

import (
	"time"

	"github.com/google/uuid"
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
	Entity
	UUID uuid.UUID
}

func (e *EntityWithUUID) BeforeCreate(tx *gorm.DB) error {
	if e.UUID == uuid.Nil {
		var err error
		e.UUID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return nil
}

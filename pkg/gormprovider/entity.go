package gormprovider

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type AbstractEntity struct {
	Id        string    `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (e *AbstractEntity) BeforeCreate(tx *gorm.DB) error {
	e.Id = ulid.Make().String()
	return nil
}

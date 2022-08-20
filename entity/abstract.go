package entity

import (
	"time"
)

type abstractEntity struct {
	AbstractEntity

	Id        string
	CreatedAt time.Time
}

type AbstractEntity interface {
	GetTableName() string
	ToMap() map[string]interface{}
}

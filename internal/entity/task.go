package entity

import (
	"time"

	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

type Task struct {
	gorm_provider.EntityWithUUID

	Title  string     `json:"title"`
	DoneAt *time.Time `json:"doneAt"`
}

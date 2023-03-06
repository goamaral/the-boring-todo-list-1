package entity

import (
	"time"
)

type Task struct {
	AbstractEntity

	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completedAt"`
}

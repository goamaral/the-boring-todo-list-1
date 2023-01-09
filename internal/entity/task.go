package entity

import (
	"time"
)

type Task struct {
	AbstractEntity

	Title       string
	CompletedAt *time.Time
}

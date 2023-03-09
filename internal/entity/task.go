package entity

import (
	"time"

	gormprovider "example.com/the-boring-to-do-list-1/pkg/provider/gorm"
)

type Task struct {
	gormprovider.AbstractEntity

	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completedAt"`
}

package repository

import (
	"time"

	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

type TaskPatch struct {
	Title  gorm_provider.OptionalField[string]
	DoneAt gorm_provider.OptionalField[*time.Time] `gorm:"type:time"`
}

func (TaskPatch) TableName() string {
	return TasksTableName
}

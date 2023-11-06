package repository

import (
	"time"

	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

type TaskPatch struct {
	Title  gorm_provider.OptionalField[string]     `json:"title"`                   // TODO: Test if json tag can be removed
	DoneAt gorm_provider.OptionalField[*time.Time] `json:"doneAt" gorm:"type:time"` // TODO: Test if json tag can be removed
}

func (TaskPatch) TableName() string {
	return tasksTableName
}

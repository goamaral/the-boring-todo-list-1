package repository

import (
	"time"

	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"gorm.io/gorm"
)

type TaskFilter struct {
	ID   gorm_provider.QueryFieldFilter[uint]
	IDGt gorm_provider.QueryFieldFilter[uint]
	UUID gorm_provider.QueryFieldFilter[string]
	Done gorm_provider.QueryFieldFilter[bool]
}

func (opt TaskFilter) Apply(db *gorm.DB) *gorm.DB {
	if opt.ID.Defined {
		db = db.Where("id", opt.ID)
	}
	if opt.IDGt.Defined {
		db = db.Where("id > ?", opt.IDGt)
	}
	if opt.UUID.Defined {
		db = db.Where("uuid", opt.UUID)
	}
	if opt.Done.Defined {
		if opt.Done.Val {
			db = db.Where("done_at IS NOT NULL")
		} else {
			db = db.Where("done_at IS NULL")
		}
	}
	return db
}

type TaskPatch struct {
	Title  gorm_provider.QueryFieldFilter[string]     `json:"title"`
	DoneAt gorm_provider.QueryFieldFilter[*time.Time] `json:"doneAt" gorm:"type:time"`
}

func (TaskPatch) TableName() string {
	return tasksTableName
}

package repository

import (
	"gorm.io/gorm"
)

type TaskFilter struct {
	UUID       *string
	IDGt       *uint
	IsComplete *bool
}

func (opt TaskFilter) Apply(db *gorm.DB) *gorm.DB {
	if opt.UUID != nil {
		db = db.Where("uuid", *opt.UUID)
	}
	if opt.IDGt != nil {
		db = db.Where("id > ?", *opt.IDGt)
	}
	if opt.IsComplete != nil {
		if *opt.IsComplete {
			db = db.Where("completed_at IS NOT NULL")
		} else {
			db = db.Where("completed_at IS NULL")
		}
	}

	return db
}

type TaskPatch struct {
	Title *string `json:"title"`
}

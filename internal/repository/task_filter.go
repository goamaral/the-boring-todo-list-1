package repository

import "gorm.io/gorm"

type TaskFilter struct {
	Id         *string
	IsComplete *bool
}

func (opt TaskFilter) Apply(db *gorm.DB) *gorm.DB {
	if opt.Id != nil {
		db = db.Where("id", *opt.Id)
	}

	if opt.IsComplete != nil {
		if *opt.IsComplete {
			db = db.Where("completed_at IS NOT NULL")
		} else {
			db = db.Where("completed_at IS NULL")
		}
	}

	return db.Debug()
}

package repository

import "gorm.io/gorm"

type TaskFilter struct {
	Id         *string
	IsComplete *bool
}

func (opt TaskFilter) Apply(db *gorm.DB) *gorm.DB {
	if opt.Id != nil {
		db.Where("id = ?", opt.Id)
	}

	return db
}

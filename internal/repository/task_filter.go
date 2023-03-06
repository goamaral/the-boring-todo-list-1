package repository

import "gorm.io/gorm"

type TaskFilter struct {
	Id string
}

func (opt TaskFilter) Apply(db *gorm.DB) *gorm.DB {
	if opt.Id != "" {
		db.Where("id = ?", opt.Id)
	}

	return db
}

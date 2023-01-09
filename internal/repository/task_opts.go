package repository

import (
	gormprovider "example.com/the-boring-to-do-list-1/pkg/provider/gorm"
	"gorm.io/gorm"
)

type ListTasksOpts struct {
	PageId   string
	PageSize int
}

func (opts *ListTasksOpts) Apply(db *gorm.DB) *gorm.DB {
	if opts == nil {
		return db
	}

	if opts.PageId != "" {
		db = db.Where("id > ?", opts.PageId)
	}

	if opts.PageSize != 0 {
		db = db.Limit(opts.PageSize)
	} else {
		db = db.Limit(gormprovider.DefaultPageSize)
	}

	return db
}

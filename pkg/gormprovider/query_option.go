package gormprovider

import "gorm.io/gorm"

type QueryOption interface {
	Apply(*gorm.DB) *gorm.DB
}

func ApplyQueryOpts(qry *gorm.DB, opts ...QueryOption) *gorm.DB {
	for _, opt := range opts {
		qry = opt.Apply(qry)
	}
	return qry
}

/* PaginationOption */
type PaginationOption struct {
	PageId   string
	PageSize int
}

func (opt PaginationOption) Apply(db *gorm.DB) *gorm.DB {
	if opt.PageId != "" {
		db = db.Where("id > ?", opt.PageId)
	}

	if opt.PageSize != 0 {
		db = db.Limit(opt.PageSize)
	} else {
		db = db.Limit(DefaultPageSize)
	}

	return db
}

/* SelectOption */
func SelectOption(fields ...string) selectOption {
	return selectOption(fields)
}

type selectOption []string

func (opt selectOption) Apply(db *gorm.DB) *gorm.DB {
	return db.Select(opt)
}

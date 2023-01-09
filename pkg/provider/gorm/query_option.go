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

package repository

import "gorm.io/gorm"

type UserFilter struct {
	Id       *string
	Username *string
}

func (opt UserFilter) Apply(db *gorm.DB) *gorm.DB {
	if opt.Id != nil {
		db = db.Where("id", *opt.Id)
	}

	if opt.Username != nil {
		db = db.Where("username", *opt.Username)
	}

	return db
}

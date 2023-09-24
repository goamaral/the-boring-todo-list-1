package gorm_provider

import (
	"fmt"

	"gorm.io/gorm"
)

type QueryOption interface {
	Apply(*gorm.DB) *gorm.DB
}

func ApplyQueryOptions(qry *gorm.DB, opts ...QueryOption) *gorm.DB {
	for _, opt := range opts {
		qry = opt.Apply(qry)
	}
	return qry
}

/* ORDER OPTION */
type OrderOption string

func (o OrderOption) Apply(qry *gorm.DB) *gorm.DB {
	return qry.Order(string(o))
}

/* DEBUG OPTION */
type debugOption struct{}

func DebugOption() debugOption {
	return debugOption{}
}

func (o debugOption) Apply(qry *gorm.DB) *gorm.DB {
	return qry.Debug()
}

/* UNSCOPED OPTION */
type unscopedOption struct{}

func UnscopedOption() unscopedOption {
	return unscopedOption{}
}

func (o unscopedOption) Apply(qry *gorm.DB) *gorm.DB {
	return qry.Unscoped()
}

/* HARD DELETE OPTION */
func HardDeleteOption() unscopedOption {
	return unscopedOption{}
}

/* SELECT OPTION */
type selectOption struct {
	query any
	args  []any
}

func SelectOption(query string, args ...any) selectOption {
	return selectOption{query: query, args: args}
}

func (o selectOption) Apply(qry *gorm.DB) *gorm.DB {
	return qry.Select(o)
}

/* JOINS OPTION */
type JoinsOption struct {
	Table      string
	Column     string
	JoinTable  string
	JoinColumn string
}

func (o JoinsOption) Apply(qry *gorm.DB) *gorm.DB {
	return qry.Joins(fmt.Sprintf("JOIN %s ON %s.%s = %s.%s", o.JoinTable, o.Table, o.Column, o.JoinTable, o.JoinColumn))
}

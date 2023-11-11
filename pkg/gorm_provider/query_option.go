package gorm_provider

import (
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueryOptionApplier interface {
	Apply(*gorm.DB) *gorm.DB
}

func ApplyQueryOptions(qry *gorm.DB, opts ...any) *gorm.DB {
	for _, opt := range opts {
		switch o := opt.(type) {
		case QueryOptionApplier:
			qry = o.Apply(qry)
		case []QueryOptionApplier:
			qry = ApplyQueryOptions(qry, lo.ToAnySlice(o)...)
		case clause.Expression:
			qry = qry.Clauses(o)
		case []clause.Expression:
			qry = ApplyQueryOptions(qry, lo.ToAnySlice(o)...)
		default:
			qry.AddError(fmt.Errorf("ApplyQueryOptions: unsupported option (%#v)", o))
		}
	}
	return qry
}

/* Option with args */
type optionWithQueryArgs struct {
	Clause clause.Expression
	Query  any
	Args   []any
}

func (o optionWithQueryArgs) Apply(qry *gorm.DB) *gorm.DB {
	switch o.Clause.(type) {
	case clause.Select:
		return qry.Select(o.Query, o.Args...)
	default:
		qry.AddError(fmt.Errorf("optionWithQueryArgs does not support (%#v)", o.Clause))
		return qry
	}
}

/* SELECT OPTION */
func SelectOption(query any, args ...any) optionWithQueryArgs {
	return optionWithQueryArgs{Clause: clause.Select{}, Query: query, Args: args}
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

/* JOINS OPTION */
type joinsOption struct {
	JoinType   clause.JoinType
	Table      string
	Column     string
	JoinTable  string
	JoinColumn string
}

func JoinsOption(joinType clause.JoinType, table string, column string, joinTable string, joinColumn string) joinsOption {
	return joinsOption{
		JoinType:   joinType,
		Table:      table,
		Column:     column,
		JoinTable:  joinTable,
		JoinColumn: joinColumn,
	}
}

func (o joinsOption) Apply(qry *gorm.DB) *gorm.DB {
	return qry.Joins(fmt.Sprintf("%s JOIN %s ON %s.%s = %s.%s", o.JoinType, o.JoinTable, o.Table, o.Column, o.JoinTable, o.JoinColumn))
}

/* PRELOAD OPTION */
type PreloadOption string

func (o PreloadOption) Apply(db *gorm.DB) *gorm.DB {
	return db.Preload(string(o))
}

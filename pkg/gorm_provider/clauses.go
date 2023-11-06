package gorm_provider

import (
	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

/* CLAUSE COLUMNS */
func ClauseColumns(columns ...string) []clause.Column {
	return lo.Map(columns, func(col string, _ int) clause.Column {
		return clause.Column{Name: col}
	})
}

/* SELECT CLAUSE */
func SelectClause(columns ...string) clause.Expression {
	return clause.Select{Columns: ClauseColumns(columns...)}
}

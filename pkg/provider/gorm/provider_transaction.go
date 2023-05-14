package gormprovider

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type TxContext struct {
	context.Context
	txDB *gorm.DB
}

func (p Provider) NewTransaction(ctx context.Context, fc func(txCtx TxContext) error, opts ...*sql.TxOptions) error {
	return p.db.Transaction(func(txDB *gorm.DB) error {
		return fc(TxContext{Context: ctx, txDB: txDB})
	})
}

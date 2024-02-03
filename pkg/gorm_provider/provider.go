package gorm_provider

import (
	"context"

	"gorm.io/gorm"
)

type AbstractProvider interface {
	GetDb() *gorm.DB
	NewTransaction(ctx context.Context, fc func(context.Context) error) error
}

type Provider struct {
	DB *gorm.DB
}

func NewProvider(dialector gorm.Dialector) (Provider, error) {
	db, err := gorm.Open(dialector)
	if err != nil {
		return Provider{}, err
	}
	return Provider{DB: db}, nil
}

func GetDbFromContextOr(ctx context.Context, db *gorm.DB) *gorm.DB {
	ctxWithTx, ok := ctx.(TxContext)
	if ok && ctxWithTx.Tx != nil {
		return ctxWithTx.Tx
	}
	return db
}

func (p Provider) GetDb() *gorm.DB {
	return p.DB
}

type TxContext struct {
	context.Context
	Tx *gorm.DB
}

func (p Provider) NewTransaction(ctx context.Context, fc func(context.Context) error) error {
	tx := GetDbFromContextOr(ctx, p.DB).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	return tx.Transaction(func(tx *gorm.DB) error {
		return fc(TxContext{Context: ctx, Tx: tx})
	})
}

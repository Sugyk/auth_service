package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// struct-key for Tx in context
type txKey struct{}

type TxManager interface {
	WithTx(context.Context, func(context.Context) error) error
}

type txManager struct {
	db *pgxpool.Pool
}

func (t *txManager) WithTx(ctx context.Context, fn func(context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)
	if fnErr := fn(txCtx); fnErr != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original error: %w)", rbErr, fnErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func NewTxManager(db *pgxpool.Pool) TxManager {
	return &txManager{
		db: db,
	}
}

type Executor interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

func TxFromCtx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

func GetExecutor(ctx context.Context, db *pgxpool.Pool) Executor {
	if tx, ok := TxFromCtx(ctx); ok {
		return tx
	}
	return db
}

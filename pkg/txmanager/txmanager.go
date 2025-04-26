package txmanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Tx is an interface to work with pgx.Conn, pgxpool.Conn or pgxpool.Pool
// StmtContext and Stmt are not implemented!
type Tx interface {
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type txKey string

var ctxWithTx = txKey("tx")

type PGXManager struct {
	db *pgxpool.Pool
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
	DoWithSettings(ctx context.Context, options pgx.TxOptions, fn func(ctx context.Context) error) error
}

func NewManager(db *pgxpool.Pool) *PGXManager {
	return &PGXManager{db: db}
}

func (m *PGXManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.DoWithSettings(ctx, pgx.TxOptions{}, fn)
}

func (m *PGXManager) DoWithSettings(ctx context.Context, options pgx.TxOptions, fn func(ctx context.Context) error) (rErr error) {
	tx, err := m.db.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if rErr != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				rErr = errors.Join(rErr, fmt.Errorf("rollback tx: %w", rollbackErr))
			}
		}
	}()

	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(error); ok {
				rErr = e
			} else {
				rErr = fmt.Errorf("panic: %v", rec)
			}
		}
	}()

	err = fn(putInCtx(ctx, tx))
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func putInCtx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxWithTx, tx)
}

type TxGetter interface {
	GetTx(ctx context.Context) Tx
}

func NewGetter(db *pgxpool.Pool) *Getter {
	return &Getter{db: db}
}

type Getter struct {
	db *pgxpool.Pool
}

func (g *Getter) GetTx(ctx context.Context) Tx {
	tx := ctx.Value(ctxWithTx)

	if t, ok := tx.(pgx.Tx); ok {
		return t
	}

	return nil
}

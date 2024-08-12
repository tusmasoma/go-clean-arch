package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/config"
	"github.com/tusmasoma/go-clean-arch/repository"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, CtxTxKey(), tx)

	var done bool
	defer func() {
		ctx = context.WithValue(ctx, CtxTxKey(), nil)
		if !done {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error("Failed to rollback transaction: %v", rollbackErr)
				if err == nil {
					err = rollbackErr
				}
			}
		}
	}()

	if err = fn(ctx); err != nil {
		return err
	}

	done = true
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

type TxKey string

func CtxTxKey() TxKey {
	return "tx"
}

func TxFromCtx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(CtxTxKey()).(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}

func NewMySQLDB(ctx context.Context) (*sql.DB, error) {
	conf, err := config.NewDBConfig(ctx)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Critical("Failed to connect to database", log.Fstring("dsn", dsn), log.Ferror(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Critical("Failed to ping database", log.Ferror(err))
		return nil, err
	}

	log.Info("Successfully connected to database", log.Fstring("dsn", dsn))
	return db, nil
}

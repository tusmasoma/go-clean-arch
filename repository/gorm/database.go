package gorm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/config"
	"github.com/tusmasoma/go-clean-arch/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := tr.db.WithContext(ctx).Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if tx.Error != nil {
		return tx.Error
	}

	ctx = context.WithValue(ctx, CtxTxKey(), tx)

	var done bool
	defer func() {
		ctx = context.WithValue(ctx, CtxTxKey(), nil)
		if !done {
			if rollbackErr := tx.Rollback(); rollbackErr.Error != nil {
				log.Error("Failed to rollback transaction", log.Ferror(rollbackErr.Error))
			}
		}
	}()

	if err := fn(ctx); err != nil {
		return err
	}

	done = true
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

type TxKey string

func CtxTxKey() TxKey {
	return "tx"
}

func TxFromCtx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(CtxTxKey()).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

const (
	dbPrefix = "MYSQL_"
)

func NewMySQLDB(ctx context.Context) (*gorm.DB, error) {
	conf, err := config.NewDBConfig(ctx, dbPrefix)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // ping is automatically called
	if err != nil {
		return nil, err
	}

	return db, nil
}

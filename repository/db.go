package repository

import (
	"context"
	"database/sql"
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/middleware/logger"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type TxRepoItf interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type TxStruct struct {
	db *sql.DB
}

type txKey struct {
	key string
}

func NewTx(db *sql.DB) *TxStruct {
	return &TxStruct{db: db}
}

func (t *TxStruct) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return dto.CustomError{
			ErrorStr:    fmt.Errorf("transaction failed: %v", err).Error(),
			InternalErr: err.Error(),
			Status:      constant.InternalServerError,
		}
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			logger.Log.Info("close transaction: %v", errTx)
		}
	}()

	if err := fn(context.WithValue(ctx, txKey{constant.TRANSACTION_KEY}, tx)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.Log.Errorf("Transaction rolled back: %v, %v", err, rbErr)
			return dto.CustomError{
				ErrorStr:    fmt.Errorf("tx err : %v, rb err: %v", err, rbErr).Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}
		logger.Log.Errorf("Transaction rolled back: %v", err)
		return err
	}

	if errCommit := tx.Commit(); errCommit != nil {
		return dto.CustomError{
			ErrorStr:    fmt.Errorf("transaction failed: %v", errCommit).Error(),
			InternalErr: errCommit.Error(),
			Status:      constant.InternalServerError,
		}
	}
	logger.Log.Info("queries committed")
	return nil
}

func contextToTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{constant.TRANSACTION_KEY}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func ChooseDbOrTx(ctx context.Context, db DBTX) DBTX {
	if tx := contextToTx(ctx); tx != nil {
		return tx
	}
	return db
}

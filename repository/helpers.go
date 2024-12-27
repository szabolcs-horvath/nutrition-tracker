package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/szabolcs-horvath/nutrition-tracker/generated"
)

const transactionActive = "TRANSACTION_ACTIVE"
const queriesPointer = "QUERIES_POINTER"

func GetQueries() (*sqlc.Queries, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return sqlc.New(db), nil
}

func DoInTransaction(ctx context.Context, db *sql.DB, action func(childCtx context.Context, queries *sqlc.Queries) error) error {
	if isTransaction(ctx) {
		// If request is already running in a transaction, don't begin a new one, reuse the *sqlc.Queries
		queries, err := getQueriesPointer(ctx)
		if err != nil {
			return err
		}
		return action(ctx, queries)
	} else {
		ctx1 := context.WithValue(ctx, transactionActive, true)
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		queries := sqlc.New(db).WithTx(tx)
		ctx2 := context.WithValue(ctx1, queriesPointer, queries)

		err = action(ctx2, queries)
		if err != nil {
			return err
		}

		return tx.Commit()
	}
}

func isTransaction(ctx context.Context) bool {
	if value, ok := ctx.Value(transactionActive).(bool); ok {
		return value
	}
	return false
}

func getQueriesPointer(ctx context.Context) (*sqlc.Queries, error) {
	if value, ok := ctx.Value(queriesPointer).(*sqlc.Queries); ok {
		return value, nil
	}
	return nil, fmt.Errorf("couldn't find queries pointer in context %v", ctx)
}

//func DoInTransaction(ctx context.Context, db *sql.DB, action func() error) error {
//	tx, err := db.BeginTx(ctx, nil)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		err = tx.Rollback()
//		if err != nil {
//			slog.ErrorContext(ctx, "[DoInTransaction] Failure during transaction rollback", err)
//			return
//		}
//	}()
//	err = action()
//	if err != nil {
//		slog.ErrorContext(ctx, "[DoInTransaction] Unexpected error happened. Rolling back the transaction...", err)
//		rollbackErr := tx.Rollback()
//		if rollbackErr != nil {
//			slog.ErrorContext(ctx, "[DoInTransaction] Failure during transaction rollback", rollbackErr)
//		}
//		return err
//	}
//	err = tx.Commit()
//	if err != nil {
//		return err
//	}
//	return nil
//}

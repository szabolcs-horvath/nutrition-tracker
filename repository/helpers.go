package repository

import (
	"context"
	"database/sql"
	"github.com/szabolcs-horvath/nutrition-tracker/generated"
)

func GetQueries() (*sqlc.Queries, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return sqlc.New(db), nil
}

func DoInTransaction(ctx context.Context, db *sql.DB, action func(*sqlc.Queries) error) error {
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	queries := sqlc.New(db).WithTx(tx)
	err = action(queries)
	if err != nil {
		return err
	}
	return tx.Commit()
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

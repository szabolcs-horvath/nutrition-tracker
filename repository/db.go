package repository

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "sqlite/nutrition-tracker.db")
}

func DoInTransaction(ctx context.Context, db *sql.DB, action func() error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			slog.ErrorContext(ctx, "Failure during transaction rollback", err)
			return
		}
	}()
	err = action()
	if err != nil {
		slog.ErrorContext(ctx, "Unexpected error happened. Rolling back the transaction...", err)
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			slog.ErrorContext(ctx, "Failure during transaction rollback", rollbackErr)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

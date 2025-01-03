package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
)

const (
	DbFileKey = "DB_FILE"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite3", util.GetEnvSafe(DbFileKey))
}

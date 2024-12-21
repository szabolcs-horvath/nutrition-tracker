package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
)

func GetDB() (*sql.DB, error) {
	return sql.Open("sqlite3", util.SafeGetEnv("DB_FILE"))
}

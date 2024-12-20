package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/xneverov/todo-list/internal/config"
)

var db *sql.DB
var dbFile = config.Get("TODO_DBFILE")

const createTableQuery = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date INTEGER NOT NULL DEFAULT 0,
	title TEXT NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat TEXT NOT NULL DEFAULT "" CHECK(length(repeat) <= 128)
);
CREATE INDEX scheduler_date ON scheduler (date);`

func Init() error {
	shouldCreateTable := !fileExists()

	if err := open(); err != nil {
		return err
	}

	if shouldCreateTable {
		return createTable()
	}

	return nil
}

func open() error {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to establish a connection to the database: %w", err)
	}

	return nil
}

func fileExists() bool {
	_, err := os.Stat(dbFile)
	return err == nil
}

func createTable() error {
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func Get() *sql.DB {
	return db
}

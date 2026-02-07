package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/mattn/go-sqlite3"
)

func saveToFile(memDB *sql.DB, filename string) error {
	// Open destination file
	fileDB, err := sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}
	defer fileDB.Close()

	// Get connections
	memConn, err := memDB.Conn(context.Background())
	if err != nil {
		return err
	}
	defer memConn.Close()

	fileConn, err := fileDB.Conn(context.Background())
	if err != nil {
		return err
	}
	defer fileConn.Close()

	// Perform backup
	return memConn.Raw(func(memDriver any) error {
		return fileConn.Raw(func(fileDriver any) error {
			srcConn, ok := memDriver.(*sqlite3.SQLiteConn)
			if !ok {
				return fmt.Errorf("invalid source connection type")
			}

			dstConn, ok := fileDriver.(*sqlite3.SQLiteConn)
			if !ok {
				return fmt.Errorf("invalid destination connection type")
			}

			backup, err := dstConn.Backup("main", srcConn, "main")
			if err != nil {
				return err
			}
			defer backup.Finish()

			_, err = backup.Step(-1) // -1 = copy entire database
			return err
		})
	})
}
func NewDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	_, _ = db.Exec("PRAGMA foreign_keys = ON;")
	_, _ = db.Exec("PRAGMA busy_timeout = 5000;") // wait 5s if locked
	return db, nil
}
func DBSchemaInit(db *sql.DB) error {
	schema, err := os.ReadFile("./internal/store/sqlite/schema.sql")

	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))

	if err != nil {
		return err
	}
	return nil
}

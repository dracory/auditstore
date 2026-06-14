package auditstore

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func initTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	return db
}

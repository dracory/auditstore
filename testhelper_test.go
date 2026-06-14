package auditstore

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func initTestDB() *sql.DB {
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		panic(err)
	}
	return db
}

package auditstore

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Store defines a udit store
type Store struct {
	automigrateEnabled bool
	auditTableName     string
	db                 *sql.DB
	dbDriverName       string
	debugEnabled       bool
}

type NewStoreOptions struct {
	AuditTableName     string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new entity store
func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
		auditTableName:     opts.AuditTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.auditTableName == "" {
		return nil, errors.New("audit store: AuditTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("audit store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = store.DriverName(store.db)
	}

	if store.automigrateEnabled {
		if err := store.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate: %w", err)
		}
	}

	return store, nil
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() error {
	sql := st.SqlCreateAuditTable()

	err := st.SqlExec(sql)
	if err != nil {
		log.Println(sql)
		log.Println(err.Error())
		return err
	}

	return nil
}

// DriverName finds the driver name from database
func (st *Store) DriverName(db *sql.DB) string {
	dv := reflect.ValueOf(db.Driver())
	driverFullName := dv.Type().String()
	if strings.Contains(driverFullName, "mysql") {
		return "mysql"
	}
	if strings.Contains(driverFullName, "postgres") || strings.Contains(driverFullName, "pq") {
		return "postgres"
	}
	if strings.Contains(driverFullName, "sqlite") {
		return "sqlite"
	}
	if strings.Contains(driverFullName, "mssql") {
		return "mssql"
	}
	return driverFullName
}

func (st *Store) SqlExec(sqlStr string) error {
	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)
	if err != nil {
		if st.debugEnabled {
			log.Println(err)
		}
		return err
	}

	return nil
}

// DebugEnable - enables / disables the debug option
func (st *Store) DebugEnable(debug bool) {
	st.debugEnabled = debug
}

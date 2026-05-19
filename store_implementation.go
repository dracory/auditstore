package auditstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// storeImplementation implements StoreInterface for audit operations
type storeImplementation struct {
	db                 *sql.DB
	debugEnabled       bool
	automigrateEnabled bool
	dbDriverName       string
	auditTableName     string
	logger             *slog.Logger
}

var _ StoreInterface = (*storeImplementation)(nil)

// NewStoreOptions contains options for creating a new Store
type NewStoreOptions struct {
	DB                 *sql.DB
	AuditTableName     string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new audit store
func NewStore(options NewStoreOptions) (StoreInterface, error) {
	if options.DB == nil {
		return nil, errors.New("database is required")
	}

	if options.AuditTableName == "" {
		return nil, errors.New("audit table name is required")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		logger:             logger,
		db:                 options.DB,
		auditTableName:     options.AuditTableName,
		automigrateEnabled: options.AutomigrateEnabled,
		debugEnabled:       options.DebugEnabled,
	}

	// Set the database driver name if not provided
	if store.dbDriverName == "" {
		store.dbDriverName = database.DatabaseType(store.db)
	}

	if store.automigrateEnabled {
		err := store.MigrateUp(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return store, nil
}

// AutoMigrate creates the audit table if it doesn't exist (deprecated - use MigrateUp)
func (st *storeImplementation) AutoMigrate() error {
	return st.MigrateUp(context.Background())
}

// MigrateUp creates the audit table
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	sqlStr := st.sqlAuditTableCreate()

	if sqlStr == "" {
		return errors.New("audit table create SQL is empty")
	}

	if st.debugEnabled {
		st.logger.Info("Running migration", "sql", sqlStr)
	}

	var err error
	if txToUse != nil {
		_, err = txToUse.ExecContext(ctx, sqlStr)
	} else {
		_, err = st.db.ExecContext(ctx, sqlStr)
	}

	if err != nil {
		if st.debugEnabled {
			st.logger.Error("Migration failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the audit table
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	var txToUse *sql.Tx
	if len(tx) > 0 {
		txToUse = tx[0]
	}

	sqlStr := st.sqlAuditTableDrop()

	if sqlStr == "" {
		return errors.New("audit table drop SQL is empty")
	}

	if st.debugEnabled {
		st.logger.Info("Running migration", "sql", sqlStr)
	}

	var err error
	if txToUse != nil {
		_, err = txToUse.ExecContext(ctx, sqlStr)
	} else {
		_, err = st.db.ExecContext(ctx, sqlStr)
	}

	if err != nil {
		if st.debugEnabled {
			st.logger.Error("Migration failed", "error", err)
		}
		return err
	}

	return nil
}

// DriverName finds the driver name from database
func (st *storeImplementation) DriverName(db *sql.DB) string {
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

func (st *storeImplementation) SqlExec(sqlStr string) error {
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

// EnableDebugMode enables or disables debug mode
func (st *storeImplementation) EnableDebugMode(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// GetAuditTableName returns the audit table name
func (st *storeImplementation) GetAuditTableName() string {
	return st.auditTableName
}

// SetAuditTableName sets the audit table name
func (st *storeImplementation) SetAuditTableName(tableName string) {
	st.auditTableName = tableName
}

// AuditGet retrieves an audit record by its ID
func (st *storeImplementation) AuditGet(id string) (RecordInterface, error) {
	if st.db == nil {
		return nil, errors.New("database is not initialized")
	}

	if id == "" {
		return nil, errors.New("audit ID is required")
	}

	query := goqu.Dialect(st.dbDriverName).From(st.auditTableName).
		Where(goqu.C("id").Eq(id)).
		Limit(1)

	sqlStr, sqlParams, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	if st.debugEnabled {
		st.logger.Debug("AuditGet query", "query", sqlStr, "params", sqlParams)
	}

	// Execute the query and get results as a slice of maps
	modelMaps, err := database.SelectToMapString(database.Context(context.Background(), st.db), sqlStr, sqlParams...)
	if err != nil {
		return nil, err
	}

	if len(modelMaps) == 0 {
		return nil, nil
	}

	// Convert the first result to a Record object
	record := NewRecordFromExistingData(modelMaps[0])
	return record, nil
}

// AuditList retrieves a list of audit records based on a query
func (st *storeImplementation) AuditList(query RecordQueryInterface) ([]RecordInterface, error) {
	// Build the select dataset
	ds, _, err := query.ToSelectDataset(st.dbDriverName, st.auditTableName)
	if err != nil {
		return nil, err
	}

	// Generate SQL and parameters
	sqlStr, sqlParams, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	if st.debugEnabled {
		st.logger.Debug("AuditList query", "query", sqlStr, "params", sqlParams)
	}

	// Execute the query and get results as a slice of maps
	modelMaps, err := database.SelectToMapString(database.Context(context.Background(), st.db), sqlStr, sqlParams...)
	if err != nil {
		return nil, err
	}

	// Convert the maps to RecordInterface objects
	records := make([]RecordInterface, 0, len(modelMaps))
	for _, modelMap := range modelMaps {
		record := NewRecordFromExistingData(modelMap)
		records = append(records, record)
	}

	return records, nil
}

// AuditCount retrieves the count of audit records based on a query
func (st *storeImplementation) AuditCount(query RecordQueryInterface) (int64, error) {
	// Build the select dataset for count
	ds, _, err := query.ToSelectDataset(st.dbDriverName, st.auditTableName)
	if err != nil {
		return 0, err
	}

	// Convert to count query
	countDs := ds.Select(goqu.COUNT("*"))

	// Generate SQL and parameters
	sqlStr, args, err := countDs.ToSQL()
	if err != nil {
		return 0, err
	}

	if st.debugEnabled {
		log.Printf("AuditCount SQL: %s, Args: %v", sqlStr, args)
	}

	var count int64
	err = st.db.QueryRow(sqlStr, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// AuditDelete deletes an audit record by its ID
func (st *storeImplementation) AuditDelete(id string) error {
	if id == "" {
		return errors.New("audit ID is required")
	}

	query := goqu.Dialect(st.dbDriverName).
		Delete(st.auditTableName).
		Where(goqu.C("id").Eq(id))

	sqlStr, args, err := query.ToSQL()
	if err != nil {
		return err
	}

	if st.debugEnabled {
		log.Printf("AuditDelete SQL: %s, Args: %v", sqlStr, args)
	}

	_, err = st.db.Exec(sqlStr, args...)
	return err
}

// DebugEnable is kept for backward compatibility
func (st *storeImplementation) DebugEnable(debug bool) {
	st.EnableDebugMode(debug)
}

// AuditCreate creates a new audit record
func (s *storeImplementation) AuditCreate(record RecordInterface) error {
	if s.db == nil {
		return errors.New("database is not initialized")
	}

	if record.ID() == "" {
		time.Sleep(1 * time.Millisecond)
		record.SetID(uid.MicroUid())
	}

	record.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	// Get the data from the record object
	data := record.Data()

	// Build the SQL query
	sqlStr, sqlParams, err := goqu.Dialect(s.dbDriverName).
		Insert(s.auditTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if err != nil {
		return err
	}

	if s.debugEnabled {
		s.logger.Debug("Audit create query", "query", sqlStr, "params", sqlParams)
	}

	// Execute the query
	_, err = database.Execute(database.Context(context.Background(), s.db), sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	record.MarkAsNotDirty()

	return nil
}

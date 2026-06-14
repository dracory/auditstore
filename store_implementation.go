package auditstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/dracory/neat"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// storeImplementation implements StoreInterface for audit operations
type storeImplementation struct {
	db                 *neat.Database
	debugEnabled       bool
	automigrateEnabled bool
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

	neatDB, err := neat.NewFromSQLDB(options.DB)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		logger:             logger,
		db:                 neatDB,
		auditTableName:     options.AuditTableName,
		automigrateEnabled: options.AutomigrateEnabled,
		debugEnabled:       options.DebugEnabled,
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

// MigrateUp creates the audit table.
// Note: the neat schema builder does not expose a transaction handle, so the tx
// parameter is accepted for interface compatibility but cannot be forwarded.
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.auditTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: table already exists", "table", st.auditTableName)
		}
		return nil
	}

	err := st.db.Schema().Create(st.auditTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 21)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_OBJECT_TYPE, 100)
		table.String(COLUMN_OBJECT_ID, 40)
		table.Text(COLUMN_VALUE_OLD)
		table.Text(COLUMN_VALUE_NEW)
		table.String(COLUMN_AUTHOR_ID, 40)
		table.DateTime(COLUMN_CREATED_AT)
	})

	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateUp failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the audit table.
// Note: the neat schema builder does not expose a transaction handle, so the tx
// parameter is accepted for interface compatibility but cannot be forwarded.
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if !st.db.Schema().HasTable(st.auditTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateDown: table does not exist", "table", st.auditTableName)
		}
		return nil
	}

	err := st.db.Schema().Drop(st.auditTableName)
	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateDown failed", "error", err)
		}
		return err
	}
	return nil
}

// EnableDebugMode enables or disables debug mode
func (st *storeImplementation) EnableDebugMode(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.db.EnableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
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
	if id == "" {
		return nil, errors.New("audit ID is required")
	}

	record := &recordImplementation{}
	err := st.db.Query().Table(st.auditTableName).Where("id", id).First(record)
	if err != nil {
		if err.Error() == "no rows found" {
			return nil, nil
		}
		if st.debugEnabled {
			st.logger.Debug("AuditGet error", "error", err)
		}
		return nil, err
	}

	return record, nil
}

// AuditList retrieves a list of audit records based on a query
func (st *storeImplementation) AuditList(query RecordQueryInterface) ([]RecordInterface, error) {
	q, err := query.ToQuery(st.db)
	if err != nil {
		return nil, err
	}

	var results []recordImplementation
	err = q.Table(st.auditTableName).Get(&results)
	if err != nil {
		if st.debugEnabled {
			st.logger.Debug("AuditList error", "error", err)
		}
		return nil, err
	}

	records := make([]RecordInterface, len(results))
	for i := range results {
		records[i] = &results[i]
	}

	return records, nil
}

// AuditCount retrieves the count of audit records based on a query
func (st *storeImplementation) AuditCount(query RecordQueryInterface) (int64, error) {
	q, err := query.ToQuery(st.db)
	if err != nil {
		return 0, err
	}

	var count int64
	err = q.Table(st.auditTableName).Count(&count)
	if err != nil {
		if st.debugEnabled {
			st.logger.Debug("AuditCount error", "error", err)
		}
		return 0, err
	}
	return count, nil
}

// AuditDelete deletes an audit record by its ID
func (st *storeImplementation) AuditDelete(id string) error {
	if id == "" {
		return errors.New("audit ID is required")
	}

	_, err := st.db.Query().Table(st.auditTableName).Where("id", id).Delete()
	return err
}

// DebugEnable is kept for backward compatibility
func (st *storeImplementation) DebugEnable(debug bool) {
	st.EnableDebugMode(debug)
}

// AuditCreate creates a new audit record.
// Accepts any RecordInterface and builds the insert map from the interface's
// getters, so custom implementations are supported without type assertions.
func (st *storeImplementation) AuditCreate(record RecordInterface) error {
	if record.ID() == "" {
		record.SetID(neatuid.GenerateShortID())
	}

	if record.CreatedAt() == "" {
		record.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	// Pass time.Time so the ORM handles dialect-specific timestamp serialisation.
	row := map[string]any{
		COLUMN_ID:          record.ID(),
		COLUMN_OBJECT_TYPE: record.ObjectType(),
		COLUMN_OBJECT_ID:   record.ObjectID(),
		COLUMN_VALUE_OLD:   record.ValueOld(),
		COLUMN_VALUE_NEW:   record.ValueNew(),
		COLUMN_AUTHOR_ID:   record.AuthorID(),
		COLUMN_CREATED_AT:  record.CreatedAtCarbon().StdTime(),
	}

	err := st.db.Query().Table(st.auditTableName).Create(row)
	if err != nil {
		if st.debugEnabled {
			st.logger.Debug("AuditCreate error", "error", err)
		}
	}
	return err
}

package auditstore

import (
	"context"
	"database/sql"
)

// StoreInterface defines the interface for audit store operations
type StoreInterface interface {
	// GetAuditTableName returns the audit table name
	GetAuditTableName() string
	// SetAuditTableName sets the audit table name
	SetAuditTableName(tableName string)

	// MigrateDown drops the audit table
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	// MigrateUp creates the audit table
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	// EnableDebugMode enables or disables debug mode
	EnableDebugMode(debug bool)

	// DebugEnable is kept for backward compatibility
	DebugEnable(debug bool)

	// AuditCreate creates a new audit record
	AuditCreate(audit RecordInterface) error

	// AuditGet retrieves an audit record by its ID
	AuditGet(id string) (RecordInterface, error)

	// AuditList retrieves a list of audit records based on a query
	AuditList(query RecordQueryInterface) ([]RecordInterface, error)

	// AuditCount retrieves the count of audit records based on a query
	AuditCount(query RecordQueryInterface) (int64, error)

	// AuditDelete deletes an audit record by its ID
	AuditDelete(id string) error

	// AutoMigrate runs database migrations
	AutoMigrate() error
}

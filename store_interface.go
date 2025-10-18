package auditstore


// StoreInterface defines the interface for audit store operations
type StoreInterface interface {
	// EnableDebugMode enables or disables debug mode
	EnableDebugMode(debug bool)

	// DebugEnable is kept for backward compatibility
	DebugEnable(debug bool)

	// AuditCreate creates a new audit record
	AuditCreate(audit AuditInterface) error

	// AuditGet retrieves an audit record by its ID
	AuditGet(id string) (AuditInterface, error)

	// AuditList retrieves a list of audit records based on a query
	AuditList(query AuditQueryInterface) ([]AuditInterface, error)

	// AuditCount retrieves the count of audit records based on a query
	AuditCount(query AuditQueryInterface) (int64, error)

	// AuditDelete deletes an audit record by its ID
	AuditDelete(id string) error

	// AutoMigrate runs database migrations
	AutoMigrate() error
}

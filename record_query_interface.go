package auditstore

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

// AuditQueryInterface defines the interface for audit query operations
type AuditQueryInterface interface {
	// SetLimit sets the maximum number of records to return
	SetLimit(limit int) AuditQueryInterface

	// SetOffset sets the number of records to skip
	SetOffset(offset int) AuditQueryInterface

	// SetOrderBy sets the field to order by and the sort direction
	SetOrderBy(field string, ascending bool) AuditQueryInterface

	// SetObjectType filters by object type
	SetObjectType(objectType string) AuditQueryInterface

	// SetObjectID filters by object ID
	SetObjectID(objectID string) AuditQueryInterface

	// SetAuthorID filters by author ID
	SetAuthorID(authorID string) AuditQueryInterface

	// SetCreatedAfter filters records created after the specified time
	SetCreatedAfter(t time.Time) AuditQueryInterface

	// SetCreatedBefore filters records created before the specified time
	SetCreatedBefore(t time.Time) AuditQueryInterface

	// ToSelectDataset builds a goqu.SelectDataset with the current query parameters
	ToSelectDataset(driver string, table string) (*goqu.SelectDataset, []interface{}, error)
}

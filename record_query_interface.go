package auditstore

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

// RecordQueryInterface defines the interface for record query operations
type RecordQueryInterface interface {
	// SetLimit sets the maximum number of records to return
	SetLimit(limit int) RecordQueryInterface

	// SetOffset sets the number of records to skip
	SetOffset(offset int) RecordQueryInterface

	// SetOrderBy sets the field to order by and the sort direction
	SetOrderBy(field string, ascending bool) RecordQueryInterface

	// SetObjectType filters by object type
	SetObjectType(objectType string) RecordQueryInterface

	// SetObjectID filters by object ID
	SetObjectID(objectID string) RecordQueryInterface

	// SetAuthorID filters by author ID
	SetAuthorID(authorID string) RecordQueryInterface

	// SetCreatedAfter filters records created after the specified time
	SetCreatedAfter(t time.Time) RecordQueryInterface

	// SetCreatedBefore filters records created before the specified time
	SetCreatedBefore(t time.Time) RecordQueryInterface

	// ToSelectDataset builds a goqu.SelectDataset with the current query parameters
	ToSelectDataset(driver string, table string) (*goqu.SelectDataset, []interface{}, error)
}

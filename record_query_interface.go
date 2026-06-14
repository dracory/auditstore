package auditstore

import (
	"time"

	"github.com/dracory/neat"
	"github.com/dracory/neat/contracts/database/orm"
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

	// ToQuery builds a orm.Query with the current query parameters
	ToQuery(db *neat.Database) (orm.Query, error)

	// Validate ensures the query has valid parameters
	Validate() error
}

package auditstore

import (
	"errors"
	"time"

	"github.com/dracory/neat"
	"github.com/dracory/neat/contracts/database/orm"
)

type recordQueryImplementation struct {
	limit            int
	limitSet         bool
	offset           int
	offsetSet        bool
	orderBy          string
	orderAsc         bool
	orderBySet       bool
	objectType       string
	objectTypeSet    bool
	objectID         string
	objectIDSet      bool
	authorID         string
	authorIDSet      bool
	createdAfter     time.Time
	createdAfterSet  bool
	createdBefore    time.Time
	createdBeforeSet bool
}

// NewRecordQuery creates a new RecordQuery instance
func NewRecordQuery() RecordQueryInterface {
	return &recordQueryImplementation{
		orderAsc: true, // Default to ascending order
	}
}

func (q *recordQueryImplementation) SetLimit(limit int) RecordQueryInterface {
	q.limit = limit
	q.limitSet = true
	return q
}

func (q *recordQueryImplementation) SetOffset(offset int) RecordQueryInterface {
	q.offset = offset
	q.offsetSet = true
	return q
}

func (q *recordQueryImplementation) SetOrderBy(field string, ascending bool) RecordQueryInterface {
	q.orderBy = field
	q.orderAsc = ascending
	q.orderBySet = true
	return q
}

func (q *recordQueryImplementation) SetObjectType(objectType string) RecordQueryInterface {
	q.objectType = objectType
	q.objectTypeSet = true
	return q
}

func (q *recordQueryImplementation) SetObjectID(objectID string) RecordQueryInterface {
	q.objectID = objectID
	q.objectIDSet = true
	return q
}

func (q *recordQueryImplementation) SetAuthorID(authorID string) RecordQueryInterface {
	q.authorID = authorID
	q.authorIDSet = true
	return q
}

func (q *recordQueryImplementation) SetCreatedAfter(t time.Time) RecordQueryInterface {
	q.createdAfter = t
	q.createdAfterSet = true
	return q
}

func (q *recordQueryImplementation) SetCreatedBefore(t time.Time) RecordQueryInterface {
	q.createdBefore = t
	q.createdBeforeSet = true
	return q
}

func (q *recordQueryImplementation) Validate() error {
	if q.limitSet && q.limit < 0 {
		return errors.New("limit cannot be negative")
	}

	if q.offsetSet && q.offset < 0 {
		return errors.New("offset cannot be negative")
	}

	if q.orderBySet && q.orderBy == "" {
		return errors.New("order_by is required when order_by is set")
	}

	if q.objectTypeSet && q.objectType == "" {
		return errors.New("object_type is required")
	}

	if q.objectIDSet {
		if q.objectID == "" {
			return errors.New("object_id is required")
		}

		if !q.objectTypeSet || q.objectType == "" {
			return errors.New("object_type is required when object_id is set")
		}
	}

	if q.authorIDSet && q.authorID == "" {
		return errors.New("author_id is required")
	}

	if q.createdAfterSet {
		if q.createdAfter.IsZero() {
			return errors.New("created_after is required")
		}
	}

	if q.createdBeforeSet {
		if q.createdBefore.IsZero() {
			return errors.New("created_before is required")
		}
	}

	if q.createdAfterSet && q.createdBeforeSet && q.createdAfter.After(q.createdBefore) {
		return errors.New("created_after cannot be after created_before")
	}

	return nil
}

// ToQuery builds a orm.Query with the current query parameters
func (q *recordQueryImplementation) ToQuery(db *neat.Database) (orm.Query, error) {
	if err := q.Validate(); err != nil {
		return nil, err
	}

	// Set default order by if not set
	orderBy := q.orderBy
	if orderBy == "" {
		orderBy = COLUMN_CREATED_AT
	}

	query := db.Query()

	// Apply filters using raw SQL fragments so operators are included in the
	// query string (neat's Where treats the first arg as a column when there is
	// exactly one additional arg and no operator in the string).
	if q.objectType != "" {
		query = query.Where(COLUMN_OBJECT_TYPE+" = ?", q.objectType)
	}

	if q.objectID != "" {
		query = query.Where(COLUMN_OBJECT_ID+" = ?", q.objectID)
	}

	if q.authorID != "" {
		query = query.Where(COLUMN_AUTHOR_ID+" = ?", q.authorID)
	}

	if !q.createdAfter.IsZero() {
		query = query.Where(COLUMN_CREATED_AT+" >= ?", q.createdAfter)
	}
	if !q.createdBefore.IsZero() {
		query = query.Where(COLUMN_CREATED_AT+" <= ?", q.createdBefore)
	}

	// Apply ordering
	direction := "asc"
	if !q.orderAsc {
		direction = "desc"
	}
	query = query.OrderBy(orderBy, direction)

	// Apply pagination
	if q.limitSet && q.limit > 0 {
		query = query.Limit(q.limit)
	}

	// Apply offset whenever it was explicitly set, including zero
	if q.offsetSet {
		query = query.Offset(q.offset)
	}

	return query, nil
}

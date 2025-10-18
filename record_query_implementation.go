package auditstore

import (
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
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

// ToSelectDataset builds a goqu.SelectDataset with the current query parameters
func (q *recordQueryImplementation) ToSelectDataset(driver string, table string) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if err := q.Validate(); err != nil {
		return nil, nil, err
	}

	// Set default order by if not set
	if q.orderBy == "" {
		q.orderBy = COLUMN_CREATED_AT
	}

	// Initialize the query
	selectDataset = goqu.Dialect(driver).From(table)

	// Apply filters
	if q.objectType != "" {
		selectDataset = selectDataset.Where(goqu.C(COLUMN_OBJECT_TYPE).Eq(q.objectType))
	}

	if q.objectID != "" {
		selectDataset = selectDataset.Where(goqu.C(COLUMN_OBJECT_ID).Eq(q.objectID))
	}

	if q.authorID != "" {
		selectDataset = selectDataset.Where(goqu.C(COLUMN_AUTHOR_ID).Eq(q.authorID))
	}

	// Apply date range filters
	if !q.createdAfter.IsZero() && !q.createdBefore.IsZero() {
		selectDataset = selectDataset.Where(
			goqu.C(COLUMN_CREATED_AT).Gte(q.createdAfter),
			goqu.C(COLUMN_CREATED_AT).Lte(q.createdBefore),
		)
	} else if !q.createdAfter.IsZero() {
		selectDataset = selectDataset.Where(goqu.C(COLUMN_CREATED_AT).Gte(q.createdAfter))
	} else if !q.createdBefore.IsZero() {
		selectDataset = selectDataset.Where(goqu.C(COLUMN_CREATED_AT).Lte(q.createdBefore))
	}

	// Apply ordering
	if q.orderAsc {
		selectDataset = selectDataset.Order(goqu.I(q.orderBy).Asc())
	} else {
		selectDataset = selectDataset.Order(goqu.I(q.orderBy).Desc())
	}

	// Apply pagination
	if q.limit > 0 {
		selectDataset = selectDataset.Limit(uint(q.limit))
	}

	if q.offset > 0 {
		selectDataset = selectDataset.Offset(uint(q.offset))
	}

	// Set columns to select
	columns = []any{"*"}

	return selectDataset, columns, nil
}

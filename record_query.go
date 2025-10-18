package auditstore

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

type RecordQuery struct {
	limit         int
	offset        int
	orderBy       string
	orderAsc      bool
	objectType    string
	objectID      string
	authorID      string
	createdAfter  time.Time
	createdBefore time.Time
}

// NewRecordQuery creates a new RecordQuery instance
func NewRecordQuery() RecordQueryInterface {
	return &RecordQuery{
		orderAsc: true, // Default to ascending order
	}
}

func (q *RecordQuery) SetLimit(limit int) RecordQueryInterface {
	q.limit = limit
	return q
}

func (q *RecordQuery) SetOffset(offset int) RecordQueryInterface {
	q.offset = offset
	return q
}

func (q *RecordQuery) SetOrderBy(field string, ascending bool) RecordQueryInterface {
	q.orderBy = field
	q.orderAsc = ascending
	return q
}

func (q *RecordQuery) SetObjectType(objectType string) RecordQueryInterface {
	q.objectType = objectType
	return q
}

func (q *RecordQuery) SetObjectID(objectID string) RecordQueryInterface {
	q.objectID = objectID
	return q
}

func (q *RecordQuery) SetAuthorID(authorID string) RecordQueryInterface {
	q.authorID = authorID
	return q
}

func (q *RecordQuery) SetCreatedAfter(t time.Time) RecordQueryInterface {
	q.createdAfter = t
	return q
}

func (q *RecordQuery) SetCreatedBefore(t time.Time) RecordQueryInterface {
	q.createdBefore = t
	return q
}

// ToSelectDataset builds a goqu.SelectDataset with the current query parameters
func (q *RecordQuery) ToSelectDataset(driver string, table string) (selectDataset *goqu.SelectDataset, columns []any, err error) {
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

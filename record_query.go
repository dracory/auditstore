package auditstore

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

type AuditQuery struct {
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

// NewAuditQuery creates a new AuditQuery instance
func NewAuditQuery() AuditQueryInterface {
	return &AuditQuery{
		orderAsc: true, // Default to ascending order
	}
}

func (q *AuditQuery) SetLimit(limit int) AuditQueryInterface {
	q.limit = limit
	return q
}

func (q *AuditQuery) SetOffset(offset int) AuditQueryInterface {
	q.offset = offset
	return q
}

func (q *AuditQuery) SetOrderBy(field string, ascending bool) AuditQueryInterface {
	q.orderBy = field
	q.orderAsc = ascending
	return q
}

func (q *AuditQuery) SetObjectType(objectType string) AuditQueryInterface {
	q.objectType = objectType
	return q
}

func (q *AuditQuery) SetObjectID(objectID string) AuditQueryInterface {
	q.objectID = objectID
	return q
}

func (q *AuditQuery) SetAuthorID(authorID string) AuditQueryInterface {
	q.authorID = authorID
	return q
}

func (q *AuditQuery) SetCreatedAfter(t time.Time) AuditQueryInterface {
	q.createdAfter = t
	return q
}

func (q *AuditQuery) SetCreatedBefore(t time.Time) AuditQueryInterface {
	q.createdBefore = t
	return q
}

// ToSelectDataset builds a goqu.SelectDataset with the current query parameters
func (q *AuditQuery) ToSelectDataset(driver string, table string) (selectDataset *goqu.SelectDataset, columns []any, err error) {
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

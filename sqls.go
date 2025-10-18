package auditstore

import "github.com/dracory/sb"

// sqlAuditTableCreate returns a SQL string for creating the audit table
func (st *storeImplementation) sqlAuditTableCreate() string {
	sql := sb.NewBuilder(st.dbDriverName).
		Table(st.auditTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		Column(sb.Column{
			Name:   COLUMN_OBJECT_TYPE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
		}).
		Column(sb.Column{
			Name:   COLUMN_OBJECT_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: COLUMN_VALUE_OLD,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_VALUE_NEW,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   COLUMN_AUTHOR_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		})

	// Add indexes
	// sql.Index(COLUMN_OBJECT_TYPE)
	// sql.Index(COLUMN_OBJECT_ID)
	// sql.Index(COLUMN_CREATED_AT)

	return sql.CreateIfNotExists()
}

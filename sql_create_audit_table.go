package auditstore

// SqlCreateFunnelTable returns a SQL string for creating the funnel table
func (st *Store) SqlCreateAuditTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.auditTableName + ` (
		id				varchar(40)	NOT NULL PRIMARY KEY,
		object_type		varchar(40) NOT NULL,
		object_id		varchar(40) NOT NULL,
		value_old		text NOT NULL,
		value_new		text NOT NULL,
		author_id		varchar(40) NOT NULL,
		created_at		datetime NOT NULL
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.auditTableName + `" (
		"id"			varchar(40) NOT NULL PRIMARY KEY,
		"object_type"	varchar(40) NOT NULL,
		"object_id"		varchar(40) NOT NULL,
		"value_old"		text NOT NULL,
		"value_new"		text NOT NULL,
		"author_id"		varchar(40) NOT NULL,
		"created_at"	timestamptz(6) NOT NULL
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.auditTableName + `" (
		"id"			varchar(40) NOT NULL PRIMARY KEY,
		"object_type"	varchar(40) NOT NULL,
		"object_id"		varchar(40) NOT NULL,
		"value_old"		text NOT NULL,
		"value_new"		text NOT NULL,
		"author_id"		varchar(40) NOT NULL,
		"created_at"	datetime NOT NULL
	)
	`

	sql := "unsupported driver " + st.dbDriverName

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}

package auditstore

import (
	"encoding/json"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/uid"
)

type AuditCreate struct {
	ID         string
	ObjectType string
	ObjectID   string
	ValueOld   interface{}
	ValueNew   interface{}
	AuthorID   string
}

// AuditCreate creates a new audit record
func (st *Store) AuditCreate(auditCreate AuditCreate) error {
	valueOldJSON, _ := json.Marshal(auditCreate.ValueOld)
	valueNewJSON, _ := json.Marshal(auditCreate.ValueNew)

	audit := Audit{}
	if auditCreate.ID == "" {
		time.Sleep(1 * time.Millisecond)
		audit.ID = uid.MicroUid()
	}
	audit.ObjectID = auditCreate.ObjectID
	audit.ObjectType = auditCreate.ObjectType
	audit.AuthorID = auditCreate.AuthorID
	audit.CreatedAt = time.Now()

	audit.ValueOld = string(valueOldJSON)
	audit.ValueNew = string(valueNewJSON)

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.auditTableName).Rows(audit).ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}

package auditstore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// auditImplementation is the private implementation of AuditInterface
type auditImplementation struct {
	dataobject.DataObject
}

// NewAudit creates a new audit record with default values
func NewAudit() AuditInterface {
	audit := &auditImplementation{
		DataObject: *dataobject.New(),
	}

	audit.SetID(uid.MicroUid())
	audit.SetObjectType("")
	audit.SetObjectID("")
	audit.SetValueOld("")
	audit.SetValueNew("")
	audit.SetAuthorID("")
	audit.SetCreatedAt(carbon.Now().ToDateTimeString())

	return audit
}

func NewAuditFromExistingData(data map[string]string) AuditInterface {
	audit := &auditImplementation{
		DataObject: *dataobject.NewFromData(data),
	}
	return audit
}

// ID returns the unique identifier of the audit record
func (a *auditImplementation) ID() string {
	return a.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the audit record
func (a *auditImplementation) SetID(id string) {
	a.Set(COLUMN_ID, id)
}

// ObjectType returns the type of the audited object
func (a *auditImplementation) ObjectType() string {
	return a.Get(COLUMN_OBJECT_TYPE)
}

// SetObjectType sets the type of the audited object
func (a *auditImplementation) SetObjectType(objectType string) {
	a.Set(COLUMN_OBJECT_TYPE, objectType)
}

// ObjectID returns the ID of the audited object
func (a *auditImplementation) ObjectID() string {
	return a.Get(COLUMN_OBJECT_ID)
}

// SetObjectID sets the ID of the audited object
func (a *auditImplementation) SetObjectID(objectID string) {
	a.Set(COLUMN_OBJECT_ID, objectID)
}

// ValueOld returns the old value of the audited object (JSON string)
func (a *auditImplementation) ValueOld() string {
	return a.Get(COLUMN_VALUE_OLD)
}

// SetValueOld sets the old value of the audited object (JSON string)
func (a *auditImplementation) SetValueOld(valueOld string) {
	a.Set(COLUMN_VALUE_OLD, valueOld)
}

// ValueNew returns the new value of the audited object (JSON string)
func (a *auditImplementation) ValueNew() string {
	return a.Get(COLUMN_VALUE_NEW)
}

// SetValueNew sets the new value of the audited object (JSON string)
func (a *auditImplementation) SetValueNew(valueNew string) {
	a.Set(COLUMN_VALUE_NEW, valueNew)
}

// AuthorID returns the ID of the user who made the change
func (a *auditImplementation) AuthorID() string {
	return a.Get(COLUMN_AUTHOR_ID)
}

// SetAuthorID sets the ID of the user who made the change
func (a *auditImplementation) SetAuthorID(authorID string) {
	a.Set(COLUMN_AUTHOR_ID, authorID)
}

// CreatedAt returns the timestamp when the audit record was created
func (a *auditImplementation) CreatedAt() string {
	return a.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the timestamp when the audit record was created
func (a *auditImplementation) SetCreatedAt(createdAt string) {
	a.Set(COLUMN_CREATED_AT, createdAt)
}

// CreatedAtCarbon returns the created at timestamp as a Carbon instance
func (a *auditImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(a.CreatedAt())
}

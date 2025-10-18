package auditstore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// recordImplementation is the private implementation of RecordInterface
type recordImplementation struct {
	dataobject.DataObject
}

// NewRecord creates a new audit record with default values
func NewRecord() RecordInterface {
	record := &recordImplementation{
		DataObject: *dataobject.New(),
	}

	record.SetID(uid.MicroUid())
	record.SetObjectType("")
	record.SetObjectID("")
	record.SetValueOld("")
	record.SetValueNew("")
	record.SetAuthorID("")
	record.SetCreatedAt(carbon.Now().ToDateTimeString())

	return record
}

func NewRecordFromExistingData(data map[string]string) RecordInterface {
	record := &recordImplementation{
		DataObject: *dataobject.NewFromData(data),
	}
	return record
}

// ID returns the unique identifier of the audit record
func (r *recordImplementation) ID() string {
	return r.Get(COLUMN_ID)
}

// SetID sets the unique identifier of the audit record
func (r *recordImplementation) SetID(id string) {
	r.Set(COLUMN_ID, id)
}

// ObjectType returns the type of the audited object
func (r *recordImplementation) ObjectType() string {
	return r.Get(COLUMN_OBJECT_TYPE)
}

// SetObjectType sets the type of the audited object
func (r *recordImplementation) SetObjectType(objectType string) {
	r.Set(COLUMN_OBJECT_TYPE, objectType)
}

// ObjectID returns the ID of the audited object
func (r *recordImplementation) ObjectID() string {
	return r.Get(COLUMN_OBJECT_ID)
}

// SetObjectID sets the ID of the audited object
func (r *recordImplementation) SetObjectID(objectID string) {
	r.Set(COLUMN_OBJECT_ID, objectID)
}

// ValueOld returns the old value of the audited object (JSON string)
func (r *recordImplementation) ValueOld() string {
	return r.Get(COLUMN_VALUE_OLD)
}

// SetValueOld sets the old value of the audited object (JSON string)
func (r *recordImplementation) SetValueOld(valueOld string) {
	r.Set(COLUMN_VALUE_OLD, valueOld)
}

// ValueNew returns the new value of the audited object (JSON string)
func (r *recordImplementation) ValueNew() string {
	return r.Get(COLUMN_VALUE_NEW)
}

// SetValueNew sets the new value of the audited object (JSON string)
func (r *recordImplementation) SetValueNew(valueNew string) {
	r.Set(COLUMN_VALUE_NEW, valueNew)
}

// AuthorID returns the ID of the user who made the change
func (r *recordImplementation) AuthorID() string {
	return r.Get(COLUMN_AUTHOR_ID)
}

// SetAuthorID sets the ID of the user who made the change
func (r *recordImplementation) SetAuthorID(authorID string) {
	r.Set(COLUMN_AUTHOR_ID, authorID)
}

// CreatedAt returns the timestamp when the audit record was created
func (r *recordImplementation) CreatedAt() string {
	return r.Get(COLUMN_CREATED_AT)
}

// SetCreatedAt sets the timestamp when the audit record was created
func (r *recordImplementation) SetCreatedAt(createdAt string) {
	r.Set(COLUMN_CREATED_AT, createdAt)
}

// CreatedAtCarbon returns the created at timestamp as a Carbon instance
func (r *recordImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(r.CreatedAt())
}

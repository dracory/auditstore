package auditstore

import (
	"time"

	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// recordImplementation is the private implementation of RecordInterface
type recordImplementation struct {
	IDValue string `json:"id" db:"id"`

	ObjectTypeField string `db:"object_type"`
	ObjectIDField   string `db:"object_id"`
	ValueOldField   string `db:"value_old"`
	ValueNewField   string `db:"value_new"`
	AuthorIDField   string `db:"author_id"`

	CreatedAtValue time.Time `json:"created_at" db:"created_at"`
	UpdatedAtValue time.Time `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for the record
func (r *recordImplementation) TableName() string {
	return "audit_log" // Default, can be overridden by store
}

// NewRecord creates a new audit record with default values
func NewRecord() RecordInterface {
	record := &recordImplementation{}
	record.SetID(uid.MicroUid())
	record.SetCreatedAt(carbon.Now().ToDateTimeString())
	return record
}

func NewRecordFromExistingData(data map[string]string) RecordInterface {
	record := &recordImplementation{
		IDValue:         data[COLUMN_ID],
		ObjectTypeField: data[COLUMN_OBJECT_TYPE],
		ObjectIDField:   data[COLUMN_OBJECT_ID],
		ValueOldField:   data[COLUMN_VALUE_OLD],
		ValueNewField:   data[COLUMN_VALUE_NEW],
		AuthorIDField:   data[COLUMN_AUTHOR_ID],
	}
	record.SetCreatedAt(data[COLUMN_CREATED_AT])
	return record
}

// ID returns the unique identifier of the audit record
func (r *recordImplementation) ID() string {
	return r.IDValue
}

// SetID sets the unique identifier of the audit record
func (r *recordImplementation) SetID(id string) {
	r.IDValue = id
}

// ObjectType returns the type of the audited object
func (r *recordImplementation) ObjectType() string {
	return r.ObjectTypeField
}

// SetObjectType sets the type of the audited object
func (r *recordImplementation) SetObjectType(objectType string) RecordInterface {
	r.ObjectTypeField = objectType
	return r
}

// ObjectID returns the ID of the audited object
func (r *recordImplementation) ObjectID() string {
	return r.ObjectIDField
}

// SetObjectID sets the ID of the audited object
func (r *recordImplementation) SetObjectID(objectID string) RecordInterface {
	r.ObjectIDField = objectID
	return r
}

// ValueOld returns the old value of the audited object (JSON string)
func (r *recordImplementation) ValueOld() string {
	return r.ValueOldField
}

// SetValueOld sets the old value of the audited object (JSON string)
func (r *recordImplementation) SetValueOld(valueOld string) RecordInterface {
	r.ValueOldField = valueOld
	return r
}

// ValueNew returns the new value of the audited object (JSON string)
func (r *recordImplementation) ValueNew() string {
	return r.ValueNewField
}

// SetValueNew sets the new value of the audited object (JSON string)
func (r *recordImplementation) SetValueNew(valueNew string) RecordInterface {
	r.ValueNewField = valueNew
	return r
}

// AuthorID returns the ID of the user who made the change
func (r *recordImplementation) AuthorID() string {
	return r.AuthorIDField
}

// SetAuthorID sets the ID of the user who made the change
func (r *recordImplementation) SetAuthorID(authorID string) RecordInterface {
	r.AuthorIDField = authorID
	return r
}

// CreatedAt returns the timestamp when the audit record was created
func (r *recordImplementation) CreatedAt() string {
	if r.CreatedAtValue.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(r.CreatedAtValue).ToDateTimeString()
}

// SetCreatedAt sets the timestamp when the audit record was created
func (r *recordImplementation) SetCreatedAt(createdAt string) RecordInterface {
	if createdAt == "" {
		return r
	}
	r.CreatedAtValue = carbon.Parse(createdAt).StdTime()
	return r
}

// CreatedAtCarbon returns the created at timestamp as a Carbon instance
func (r *recordImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(r.CreatedAtValue)
}

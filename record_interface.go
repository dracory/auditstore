package auditstore

import (
	"github.com/dromara/carbon/v2"
)

// RecordInterface represents an audit record in the system
type RecordInterface interface {
	// ID returns the unique identifier of the audit record
	ID() string
	SetID(id string)

	// ObjectType returns the type of the audited object
	ObjectType() string
	SetObjectType(objectType string) RecordInterface

	// ObjectID returns the ID of the audited object
	ObjectID() string
	SetObjectID(objectID string) RecordInterface

	// ValueOld returns the old value of the audited object (JSON string)
	ValueOld() string
	SetValueOld(valueOld string) RecordInterface

	// ValueNew returns the new value of the audited object (JSON string)
	ValueNew() string
	SetValueNew(valueNew string) RecordInterface

	// AuthorID returns the ID of the user who made the change
	AuthorID() string
	SetAuthorID(authorID string) RecordInterface

	// CreatedAt returns the timestamp when the audit record was created
	CreatedAt() string
	SetCreatedAt(createdAt string) RecordInterface

	// CreatedAtCarbon returns the created at timestamp as a Carbon instance
	CreatedAtCarbon() *carbon.Carbon
}

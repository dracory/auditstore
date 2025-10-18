package auditstore

import "time"

type Audit struct {
	ID         string    `db:"id"`          // varchar 40, primary key
	ObjectType string    `db:"object_type"` // varchar 40
	ObjectID   string    `db:"object_id"`   // varchar 40
	ValueOld   string    `db:"value_old"`   // text
	ValueNew   string    `db:"value_new"`   // text
	AuthorID   string    `db:"author_id"`   // varchar 40
	CreatedAt  time.Time `db:"created_at"`  // varchar 40, primary key

}

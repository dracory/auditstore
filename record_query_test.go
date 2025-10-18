package auditstore

import (
	"testing"
	"time"
)

func TestRecordQueryValidate(t *testing.T) {
	t.Run("default query is valid", func(t *testing.T) {
		query := NewRecordQuery()

		if err := query.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("negative limit", func(t *testing.T) {
		err := NewRecordQuery().
			SetLimit(-1).
			Validate()

		if err == nil {
			t.Fatal("expected error for negative limit")
		}
	})

	t.Run("object id requires object type", func(t *testing.T) {
		err := NewRecordQuery().
			SetObjectID("object-1").
			Validate()

		if err == nil {
			t.Fatal("expected error when object type missing")
		}
	})

	t.Run("created after cannot be after created before", func(t *testing.T) {
		after := time.Now()
		before := after.Add(-time.Hour)

		err := NewRecordQuery().
			SetCreatedAfter(after).
			SetCreatedBefore(before).
			Validate()

		if err == nil {
			t.Fatal("expected error when created_after > created_before")
		}
	})
}

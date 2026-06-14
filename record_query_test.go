package auditstore

import (
	"testing"
	"time"

	_ "modernc.org/sqlite"
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

	t.Run("negative offset", func(t *testing.T) {
		err := NewRecordQuery().
			SetOffset(-1).
			Validate()

		if err == nil {
			t.Fatal("expected error for negative offset")
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

func TestRecordQueryToQuery(t *testing.T) {
	store := initStore(t)

	// Seed a few records with different authors and object types
	for i, author := range []string{"alice", "bob", "alice"} {
		objectID := "obj-1"
		if i == 1 {
			objectID = "obj-2"
		}
		rec := NewRecord().
			SetObjectType("post").
			SetObjectID(objectID).
			SetAuthorID(author)
		if err := store.AuditCreate(rec); err != nil {
			t.Fatalf("seed error: %v", err)
		}
	}

	t.Run("filter by author", func(t *testing.T) {
		audits, err := store.AuditList(NewRecordQuery().SetObjectType("post").SetAuthorID("alice"))
		if err != nil {
			t.Fatal(err)
		}
		if len(audits) != 2 {
			t.Fatalf("expected 2 records for alice, got %d", len(audits))
		}
	})

	t.Run("filter by object type and id", func(t *testing.T) {
		audits, err := store.AuditList(
			NewRecordQuery().SetObjectType("post").SetObjectID("obj-2"),
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(audits) != 1 {
			t.Fatalf("expected 1 record for obj-2, got %d", len(audits))
		}
	})

	t.Run("limit and offset zero returns first page", func(t *testing.T) {
		audits, err := store.AuditList(
			NewRecordQuery().SetLimit(2).SetOffset(0),
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(audits) != 2 {
			t.Fatalf("expected 2 records with limit=2 offset=0, got %d", len(audits))
		}
	})

	t.Run("offset skips records", func(t *testing.T) {
		all, err := store.AuditList(NewRecordQuery())
		if err != nil {
			t.Fatal(err)
		}
		page2, err := store.AuditList(NewRecordQuery().SetLimit(10).SetOffset(2))
		if err != nil {
			t.Fatal(err)
		}
		if len(page2) != len(all)-2 {
			t.Fatalf("expected %d records at offset 2, got %d", len(all)-2, len(page2))
		}
	})

	t.Run("date range includes all recent records", func(t *testing.T) {
		past := time.Now().Add(-24 * time.Hour)
		future := time.Now().Add(24 * time.Hour)
		audits, err := store.AuditList(
			NewRecordQuery().SetCreatedAfter(past).SetCreatedBefore(future),
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(audits) != 3 {
			t.Fatalf("expected 3 records in past-future range, got %d", len(audits))
		}
	})

	t.Run("date range excludes future records", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		farFuture := time.Now().Add(48 * time.Hour)
		audits, err := store.AuditList(
			NewRecordQuery().SetCreatedAfter(future).SetCreatedBefore(farFuture),
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(audits) != 0 {
			t.Fatalf("expected 0 records in far-future window, got %d", len(audits))
		}
	})

	t.Run("descending order", func(t *testing.T) {
		audits, err := store.AuditList(
			NewRecordQuery().SetOrderBy(COLUMN_CREATED_AT, false),
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(audits) == 0 {
			t.Fatal("expected records")
		}
		// Verify each record's created_at is >= the next one
		for i := 1; i < len(audits); i++ {
			prev := audits[i-1].CreatedAtCarbon()
			curr := audits[i].CreatedAtCarbon()
			if prev.Lt(curr) {
				t.Errorf("expected descending order at index %d", i)
			}
		}
	})
}
